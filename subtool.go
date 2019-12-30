package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// A friendly chat about Matroska metadata track numbers.
//
// Matroska tracks numbers are confusing. Tracks are stored in the file
// starting at 1 (ONE). Some mkvtoolnix commands such as mkvmerge and
// mkvextract expect tracks to start at offset zero (ZERO), while others like
// mkvpropedit, expect offset 1. Due to this, the following conventions were
// adopted here:
//
// - Tracks are always displayed as they appear in the file (base = 1), thus
//   all user selectable tracks expect offset = 1.
// - Any actions using mkvpropedit use the track number unchanged.
// - Any actions using mkvmerge or mkvextract will automatically subtract one
//   from the track number.

// Track Types. See https://www.matroska.org/technical/specs/index.html
const (
	typeVideo    = 1
	typeAudio    = 2
	typeComplex  = 3
	typeLogo     = 16
	typeSubtitle = 17
	typeButtons  = 18
	typeControl  = 32
)

type trackinfo struct {
	number      int64
	uid         int64
	name        string
	tracktype   int64
	language    string
	flagDefault bool
	CodecID     string
}

// trackFileInfo holds information about an exported track file.
type trackFileInfo struct {
	language string
	fname    string
}

// show lists all tracks in a file.
func show(p mkvParser, showUID bool) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	if showUID {
		tab.AppendHeader(table.Row{"Number", "UID", "Type", "Name", "Language", "Codec", "Default"})
	} else {
		tab.AppendHeader(table.Row{"Number", "Type", "Name", "Language", "Codec", "Default"})
	}

	for _, t := range p.tracks {
		// Create a row with the desired columns.
		row := []interface{}{t.number}
		if showUID {
			row = append(row, uint64(t.uid))
		}
		row = append(row, trackType(t.tracktype), t.name, t.language, t.CodecID)

		// Make default flag easier to see.
		if t.flagDefault {
			row = append(row, "<=====")
		} else {
			row = append(row, "")
		}
		tab.AppendRow(row)
	}
	tab.Render()
}

// setdefault resets flagDefault on all subtitle tracks and sets it on the chosen track UID.
func setdefault(p mkvParser, track int64, cmd runner) error {
	command := []string{
		"mkvpropedit",
		p.fname,
	}

	for _, t := range p.tracks {
		if t.tracktype == typeSubtitle {
			command = append(command, "--edit", fmt.Sprintf("track:%d", t.number), "--set", "flag-default=0")
		}
	}

	if err := cmd.run(command[0], command[1:]...); err != nil {
		return err
	}
	// Tracks selected by the user have offset = 0 so we make them offset = 1.
	return adddefault(p.fname, track, cmd)
}

// trackByLanguage returns the track number (base 1) for the first track with
// one of the specified languages. The list of languages works as a priority,
// meaning that languages=["eng","fra"] will first attempt to find a track with
// the English language, and failing that, French. The special language
// "default" will cause tracks without a language code to be selected (Matroska
// has the concept of a "default language", which is usually English -- tracks
// with this language will not have a language code).
func trackByLanguage(p mkvParser, languages []string) (int64, error) {
	for _, lang := range languages {
		if lang == "default" {
			lang = ""
		}
		for _, t := range p.tracks {
			if t.tracktype == typeSubtitle && t.language == lang {
				return t.number, nil
			}
		}
	}
	return 0, fmt.Errorf("no track with language(s): %s", strings.Join(languages, ","))
}

// extract extracts a given track into a file.
func extract(handler mkvParser, track int64, cmd runner) (trackFileInfo, error) {
	// Fetch language for the track. Fail if track does not exist.
	ti, err := trackInfo(handler, track)
	if err != nil {
		return trackFileInfo{}, err
	}

	// Extract into a temporary file
	tmpfile, err := ioutil.TempFile("", "subtool")
	if err != nil {
		return trackFileInfo{}, err
	}
	temp := tmpfile.Name()
	tmpfile.Close()

	// Note: mkvextract uses 0 for the first track number.
	command := []string{
		"mkvextract",
		handler.fname,
		"tracks",
		fmt.Sprintf("%d:%s", track-1, temp),
	}
	if err := cmd.run(command[0], command[1:]...); err != nil {
		return trackFileInfo{}, err
	}
	return trackFileInfo{language: ti.language, fname: temp}, nil
}

// submux merges an input file (usually an mkv file) and multiple subtitles into a
// destination, optionally removing all other subtitles from the source.
func submux(infile, outfile string, nosubs bool, cmd runner, subs ...trackFileInfo) error {
	cmdline := []string{"mkvmerge", "-o", outfile}

	if nosubs {
		cmdline = append(cmdline, "-S")
	}
	cmdline = append(cmdline, infile)

	for _, sub := range subs {
		cmdline = append(cmdline, "--language", fmt.Sprintf("0:%s", sub.language))
		cmdline = append(cmdline, sub.fname)
	}
	return cmd.run(cmdline[0], cmdline[1:]...)

}

// remux re-multiplexes the input file(s) into the output file without changes.
// This is useful to fix problems in poorly assembled Matroska files.
func remux(infiles []string, outfile string, cmd runner) error {
	cmdline := []string{"mkvmerge"}
	cmdline = append(cmdline, infiles...)
	cmdline = append(cmdline, "-o", outfile)

	return cmd.run(cmdline[0], cmdline[1:]...)
}

// adddefault adds the default flag to a given track UID.
func adddefault(mkvfile string, track int64, cmd runner) error {
	return cmd.run("mkvpropedit", mkvfile, "--edit", fmt.Sprintf("track:%d", track), "--set", "flag-default=1")
}

// trackInfo returns the trackinfo for a given track number, or error if it does not exist.
func trackInfo(handler mkvParser, track int64) (trackinfo, error) {
	for _, v := range handler.tracks {
		if v.number == track {
			return v, nil
		}
	}
	return trackinfo{}, fmt.Errorf("track number %d not found in file %s\n", track, handler.fname)
}

// trackType returns the string type of the track from the numeric track type value
// or Unknown(value) if the type is not known.
func trackType(t int64) string {
	var ttypes = map[int64]string{
		typeVideo:    "Video",
		typeAudio:    "Audio",
		typeComplex:  "Complex",
		typeLogo:     "Logo",
		typeSubtitle: "Subtitle",
		typeButtons:  "Buttons",
		typeControl:  "Control",
	}
	if v, ok := ttypes[t]; ok {
		return v
	}
	return fmt.Sprintf("Unknown(%d)", t)
}

// requirements returns nil if all required tools are installed and an error indicating
// the tools missing otherwise.
func requirements() error {
	var tools = []string{"mkvextract", "mkvmerge", "mkvpropedit"}

	missing := []string{}
	for _, t := range tools {
		_, err := exec.LookPath(t)
		if err != nil {
			missing = append(missing, t)
		}
	}
	if len(missing) != 0 {
		return fmt.Errorf("required 3rd party tool(s) missing: %s", strings.Join(missing, ","))
	}
	return nil
}

func main() {
	var (
		app    = kingpin.New("subtool", "Subtitle operations on matroska containers.")
		dryrun = app.Flag("dry-run", "Dry-run mode (only show commands).").Short('n').Bool()

		// add
		addCmd    = app.Command("add", "Add an input track (subtitle/video/audio) into an output file.")
		addOutput = addCmd.Flag("output", "Output file.").Required().Short('o').String()
		addInputs = addCmd.Arg("input-files", "Input files.").Required().Strings()

		// only
		onlyCmd       = app.Command("only", "Remove all subtitle tracks, except one.")
		setOnlyTrack  = onlyCmd.Arg("track", "Track number to keep.").Required().Int64()
		setOnlyInput  = onlyCmd.Arg("input", "Matroska Input file.").Required().String()
		setOnlyOutput = onlyCmd.Arg("output", "Matroska Output file.").Required().String()

		// remux
		remuxCmd       = app.Command("remux", "Remux input file into an output file.")
		remuxCmdInput  = remuxCmd.Arg("input-file", "Matroska Input file.").Required().String()
		remuxCmdOutput = remuxCmd.Arg("output-file", "Matroska Output file.").Required().String()

		// setdefault
		setDefaultCmd   = app.Command("setdefault", "Set default subtitle tag on a track.")
		setDefaultTrack = setDefaultCmd.Arg("track", "Track number to set as default.").Required().Int64()
		setDefaultFile  = setDefaultCmd.Arg("mkvfile", "Matroska file.").Required().String()

		// setdefaultbylanguage
		setDefaultByLangCmd  = app.Command("setdefaultbylang", "Set default subtitle track by language.")
		setDefaultByLangList = setDefaultByLangCmd.Flag("lang", "Preferred languages (Use multiple times. Use 'default' for tracks with no language set.)").Required().Strings()
		setDefaultByLangFile = setDefaultByLangCmd.Arg("mkvfile", "Matroska file.").Required().String()

		// show
		showCmd  = app.Command("show", "Show Information about a file.")
		showUID  = showCmd.Flag("uid", "Include track UIDs in the output.").Short('u').Bool()
		showFile = showCmd.Arg("input-file", "Matroska Input file.").Required().String()

		// Command runner.
		runCmd runCommand

		// Dry-run command runner (only print commands).
		fakeRunCmd fakeRunCommand

		run runner
	)

	if err := requirements(); err != nil {
		log.Fatalf("Requirements check: %v", err)
	}

	// Plain logs.
	log.SetFlags(0)

	k := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Run will resolve to a print-only version when dry-run is chosen.
	run = runCmd
	if *dryrun {
		run = fakeRunCmd
	}

	var err error

	switch k {
	case addCmd.FullCommand():
		err = remux(*addInputs, *addOutput, run)

	case onlyCmd.FullCommand():
		h := mustParseFile(*setOnlyInput)
		var tfi trackFileInfo
		tfi, err = extract(h, *setOnlyTrack, run)
		if err != nil {
			break
		}
		err = submux(*setOnlyInput, *setOnlyOutput, true, run, tfi)
		// Attempt to remove even on error.
		os.Remove(tfi.fname)

	case remuxCmd.FullCommand():
		err = remux([]string{*remuxCmdInput}, *remuxCmdOutput, run)

	case setDefaultCmd.FullCommand():
		h := mustParseFile(*setDefaultFile)
		err = setdefault(h, *setDefaultTrack, run)

	case setDefaultByLangCmd.FullCommand():
		h := mustParseFile(*setDefaultByLangFile)
		var track int64
		track, err = trackByLanguage(h, *setDefaultByLangList)
		if err != nil {
			break
		}
		err = setdefault(h, track, run)

	case showCmd.FullCommand():
		h := mustParseFile(*showFile)
		show(h, *showUID)
	}

	// Print error message, if any
	if err != nil {
		log.Fatalf("Error during %s: %v", k, err)
	}
}
