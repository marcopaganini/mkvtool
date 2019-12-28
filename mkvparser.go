package main

import (
	"github.com/remko/go-mkvparse"
	"log"
	"os"
	"time"
)

type mkvParser struct {
	track   trackinfo
	tracks  []trackinfo
	fname   string
	inTrack bool
}

func (p *mkvParser) HandleMasterBegin(id mkvparse.ElementID, info mkvparse.ElementInfo) (bool, error) {
	//fmt.Printf("==> %v\n", mkvparse.NameForElementID(id))
	if id == mkvparse.TrackEntryElement {
		p.inTrack = true
	}

	return true, nil
}

func (p *mkvParser) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	// If we're inside a track and found another track start, process the current one.
	if id == mkvparse.TrackEntryElement {
		p.tracks = append(p.tracks, p.track)
		p.track = trackinfo{}
	}
	return nil
}

func (p *mkvParser) HandleString(id mkvparse.ElementID, value string, info mkvparse.ElementInfo) error {
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

func (p *mkvParser) HandleInteger(id mkvparse.ElementID, value int64, info mkvparse.ElementInfo) error {
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

func (p *mkvParser) HandleFloat(id mkvparse.ElementID, value float64, info mkvparse.ElementInfo) error {
	return nil
}

func (p *mkvParser) HandleDate(id mkvparse.ElementID, value time.Time, info mkvparse.ElementInfo) error {
	return nil
}

func (p *mkvParser) HandleBinary(id mkvparse.ElementID, value []byte, info mkvparse.ElementInfo) error {
	return nil
}

// mustParseFile parses the MKV file and returns a handler, or aborts with an
// error message in case of problems.
func mustParseFile(fname string) mkvParser {
	handler := mkvParser{fname: fname}
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Only parse the sections we want.
	if err = mkvparse.ParseSections(f, &handler, mkvparse.TracksElement); err != nil {
		log.Fatal(err)
	}
	return handler
}
