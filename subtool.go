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
		// Commands use track starting at offset zero, hence the subtraction below.
		row := []interface{}{t.number - 1}
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
	return adddefault(p.fname, track+1, cmd)
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

	command := []string{
		"mkvextract",
		handler.fname,
		"tracks",
		fmt.Sprintf("%d:%s", track, temp),
	}
	if err := cmd.run(command[0], command[1:]...); err != nil {
		return trackFileInfo{}, err
	}
	return trackFileInfo{language: ti.language, fname: temp}, nil
}

// mux merges an input file (usually an mkv file) and multiple subtitles into a
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

		// show
		showCmd  = app.Command("show", "Show Information about a file.")
		showUID  = showCmd.Flag("uid", "Include track UIDs in the output.").Short('u').Bool()
		showFile = showCmd.Arg("input-file", "Matroska Input file.").Required().String()

		// setdefault
		setDefaultCmd   = app.Command("setdefault", "Set default subtitle tag on a track.")
		setDefaultTrack = setDefaultCmd.Arg("track", "Track number to set as default.").Required().Int64()
		setDefaultFile  = setDefaultCmd.Arg("mkvfile", "Matroska file.").Required().String()

		// setdefaultbylanguage
		setDefaultByLangCmd  = app.Command("setdefaultbylang", "Set default subtitle track by language.")
		setDefaultByLangList = setDefaultByLangCmd.Flag("lang", "Preferred languages (Use multiple times. Use 'default' for tracks with no language set.)").Required().Strings()
		setDefaultByLangFile = setDefaultByLangCmd.Arg("mkvfile", "Matroska file.").Required().String()

		// only
		setOnlyCmd    = app.Command("only", "Remove all subtitle tracks, except one.")
		setOnlyTrack  = setOnlyCmd.Arg("track", "Track number to keep.").Required().Int64()
		setOnlyInput  = setOnlyCmd.Arg("input", "Matroska Input file.").Required().String()
		setOnlyOutput = setOnlyCmd.Arg("output", "Matroska Output file.").Required().String()

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

	switch k {
	case showCmd.FullCommand():
		h := mustParseFile(*showFile)
		show(h, *showUID)
	case setDefaultCmd.FullCommand():
		h := mustParseFile(*setDefaultFile)
		if err := setdefault(h, *setDefaultTrack, run); err != nil {
			log.Fatalf("Error setting default track: %v", err)
		}
	case setDefaultByLangCmd.FullCommand():
		h := mustParseFile(*setDefaultByLangFile)
		track, err := trackByLanguage(h, *setDefaultByLangList)
		if err != nil {
			log.Fatalf("Error setting default track by language: %v", err)
		}
		if err := setdefault(h, track, run); err != nil {
			log.Fatalf("Error setting default track: %v", err)
		}
	case setOnlyCmd.FullCommand():
		h := mustParseFile(*setOnlyInput)
		tfi, err := extract(h, *setOnlyTrack, run)
		if err != nil {
			log.Fatalf("Error extracting track: %v", err)
		}
		err = submux(*setOnlyInput, *setOnlyOutput, true, run, tfi)
		if err != nil {
			log.Fatalf("Error adding subtitle: %v", err)
		}
		os.Remove(tfi.fname)
	}
}
