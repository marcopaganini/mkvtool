package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/remko/go-mkvparse"
	"gopkg.in/alecthomas/kingpin.v2"
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

type MyParser struct {
	track   trackinfo
	tracks  []trackinfo
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

	for _, t := range p.tracks {
		tab.AppendRow([]interface{}{t.number, t.uid, t.tracktype, t.name, t.language, t.CodecID, t.flagDefault})
	}
	tab.Render()
}

// setdefault resets flagDefault on all subtitle tracks and sets it on the chosen track UID.
func setdefault(mkvfile string, p MyParser, trackUID int64, cmd runner) error {
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
	return adddefault(mkvfile, trackUID, cmd)
}

// adddefault adds the default flag to a given track UID.
func adddefault(mkvfile string, trackUID int64, cmd runner) error {
	return cmd.run("mkvpropedit", mkvfile, "--edit", fmt.Sprintf("track:=%d", trackUID), "--set", "flag-default=1")
}

func main() {
	var (
		app     = kingpin.New("subtool", "Subtitle operations on matroska containers")
		mkvfile = app.Flag("file", "Matroska input file").Short('f').Required().String()
		dryrun  = app.Flag("dry-run", "Dry-run mode (only show commands)").Short('n').Bool()

		//debug = app.Flag("debug", "Enable debug mode.").Bool()
		// show
		showCmd = app.Command("show", "Show Information about a file.")

		// setdefault
		setDefaultCmd = app.Command("setdefault", "Set default subtitle tag")
		trackUID      = setDefaultCmd.Arg("trackUID", "Track UID to set as default").Required().Int64()

		// Command runners.
		runCmd     runCommand
		fakeRunCmd fakeRunCommand
		run        runner
	)

	// Plain logs.
	log.SetFlags(0)

	k := kingpin.MustParse(app.Parse(os.Args[1:]))

	handler := MyParser{}
	err := mkvparse.ParsePath(*mkvfile, &handler)
	if err != nil {
		log.Fatal(err)
	}

	run = runCmd
	if *dryrun {
		run = fakeRunCmd
	}

	switch k {
	case showCmd.FullCommand():
		show(handler)
	case setDefaultCmd.FullCommand():
		if err := setdefault(*mkvfile, handler, *trackUID, run); err != nil {
			log.Fatal(err)
		}
	}
}
