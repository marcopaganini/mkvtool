package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/remko/go-mkvparse"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
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

type MyParser struct {
	track   trackinfo
	tracks  []trackinfo
	fname   string
	inTrack bool
}

func (p *MyParser) HandleMasterBegin(id mkvparse.ElementID, info mkvparse.ElementInfo) (bool, error) {
	//fmt.Printf("==> %v\n", mkvparse.NameForElementID(id))
	// Skip large elements.
	if id == mkvparse.CuesElement || id == mkvparse.ClusterElement {
		return false, nil
	}

	if id == mkvparse.TrackEntryElement {
		p.inTrack = true
	}

	return true, nil
}

func (p *MyParser) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	// If we're inside a track and found another track start, process the current one.
	if id == mkvparse.TrackEntryElement {
		p.tracks = append(p.tracks, p.track)
		p.track = trackinfo{}
	}
	return nil
}

func (p *MyParser) HandleString(id mkvparse.ElementID, value string, info mkvparse.ElementInfo) error {
	if !p.inTrack {
		return nil
	}
	switch id {
	case mkvparse.NameElement:
		p.track.name = value
	case mkvparse.LanguageElement:
		p.track.language = value
	case mkvparse.CodecIDElement:
		p.track.CodecID = value
	}
	//fmt.Printf("%v: %q\n", mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleInteger(id mkvparse.ElementID, value int64, info mkvparse.ElementInfo) error {
	if !p.inTrack {
		return nil
	}
	//fmt.Printf("%v: %v\n", mkvparse.NameForElementID(id), value)
	switch id {
	case mkvparse.TrackNumberElement:
		p.track.number = value
	case mkvparse.TrackUIDElement:
		p.track.uid = value
	case mkvparse.TrackTypeElement:
		p.track.tracktype = value
	case mkvparse.FlagDefaultElement:
		if value != 0 {
			p.track.flagDefault = true
		}
	}
	return nil
}

func (p *MyParser) HandleFloat(id mkvparse.ElementID, value float64, info mkvparse.ElementInfo) error {
	return nil
}

func (p *MyParser) HandleDate(id mkvparse.ElementID, value time.Time, info mkvparse.ElementInfo) error {
	return nil
}

func (p *MyParser) HandleBinary(id mkvparse.ElementID, value []byte, info mkvparse.ElementInfo) error {
	return nil
}

// show lists all tracks in a file.
func show(p MyParser) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"Number", "UID", "Type", "Name", "Language", "Codec", "Default"})

	// Commands use track starting at offset zero, hence the subtraction below.
	for _, t := range p.tracks {
		tab.AppendRow([]interface{}{t.number - 1, t.uid, t.tracktype, t.name, t.language, t.CodecID, t.flagDefault})
	}
	tab.Render()
}

// setdefault resets flagDefault on all subtitle tracks and sets it on the chosen track UID.
func setdefault(mkvfile string, p MyParser, track int64, cmd runner) error {
	command := []string{
		"mkvpropedit",
		mkvfile,
	}

	for _, t := range p.tracks {
		if t.tracktype == typeSubtitle {
			command = append(command, "--edit", fmt.Sprintf("track:=%d", t.uid), "--set", "flag-default=0")
		}
	}

	if err := cmd.run(command[0], command[1:]...); err != nil {
		return err
	}
	return adddefault(mkvfile, track, cmd)
}

// extract extracts a given track into a file.
func extract(handler MyParser, track int64, cmd runner) (trackFileInfo, error) {
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

	for track, sub := range subs {
		cmdline = append(cmdline, "--language", fmt.Sprintf("%d:%s", track, sub.language))
		cmdline = append(cmdline, sub.fname)
	}
	return cmd.run(cmdline[0], cmdline[1:]...)

}

// adddefault adds the default flag to a given track UID.
func adddefault(mkvfile string, track int64, cmd runner) error {
	return cmd.run("mkvpropedit", mkvfile, "--edit", fmt.Sprintf("track:%d", track), "--set", "flag-default=1")
}

// mustParseFile parses the MKV file and returns a handler, or aborts with an
// error message in case of problems.
func mustParseFile(fname string) MyParser {
	handler := MyParser{fname: fname}
	err := mkvparse.ParsePath(fname, &handler)
	if err != nil {
		log.Fatal(err)
	}
	return handler
}

// trackInfo returns the trackinfo for a given track number, or error if it does not exist.
func trackInfo(handler MyParser, track int64) (trackinfo, error) {
	for _, v := range handler.tracks {
		if v.number == track {
			return v, nil
		}
	}
	return trackinfo{}, fmt.Errorf("track number %d not found in file %s\n", track, handler.fname)
}

func main() {
	var (
		app    = kingpin.New("subtool", "Subtitle operations on matroska containers")
		dryrun = app.Flag("dry-run", "Dry-run mode (only show commands)").Short('n').Bool()

		// show
		showCmd  = app.Command("show", "Show Information about a file.")
		showFile = showCmd.Arg("input-file", "Matroska Input file").Required().String()

		// setdefault
		setDefaultCmd   = app.Command("setdefault", "Set default subtitle tag")
		setDefaultFile  = setDefaultCmd.Arg("input-file", "Matroska Input file").Required().String()
		setDefaultTrack = setDefaultCmd.Arg("track", "Track number to set as default").Required().Int64()

		// only
		setOnlyCmd    = app.Command("only", "Remove all subtitle tracks, except one")
		setOnlyTrack  = setOnlyCmd.Arg("track", "Track number to keep").Required().Int64()
		setOnlyInput  = setOnlyCmd.Arg("input", "Matroska Input file").Required().String()
		setOnlyOutput = setOnlyCmd.Arg("output", "Matroska Output file").Required().String()

		// Command runner.
		runCmd runCommand

		// Dry-run command runner (only print commands).
		fakeRunCmd fakeRunCommand

		run runner
	)

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
		show(h)
	case setDefaultCmd.FullCommand():
		h := mustParseFile(*setDefaultFile)
		if err := setdefault(*setDefaultFile, h, *setDefaultTrack, run); err != nil {
			log.Fatal(err)
		}
	case setOnlyCmd.FullCommand():
		h := mustParseFile(*setOnlyInput)
		tfi, err := extract(h, *setOnlyTrack, run)
		if err != nil {
			log.Fatal(err)
		}
		err = submux(*setOnlyInput, *setOnlyOutput, true, run, tfi)
		if err != nil {
			log.Fatal(err)
		}
		os.Remove(tfi.fname)
	}
}
