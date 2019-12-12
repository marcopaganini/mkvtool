package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/remko/go-mkvparse"
	"os"
	"time"
)

type trackinfo struct {
	number      int64
	uid         int64
	name        string
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
		//fmt.Printf("Got track %+v\n", p.track)
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
	switch id {
	case mkvparse.TrackNumberElement:
		p.track.number = value
	case mkvparse.TrackUIDElement:
		p.track.uid = value
	case mkvparse.FlagDefaultElement:
		if value != 0 {
			p.track.flagDefault = true
		}
	}
	//fmt.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleFloat(id mkvparse.ElementID, value float64, info mkvparse.ElementInfo) error {
	//fmt.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleDate(id mkvparse.ElementID, value time.Time, info mkvparse.ElementInfo) error {
	//fmt.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleBinary(id mkvparse.ElementID, value []byte, info mkvparse.ElementInfo) error {
	/*
		if id == mkvparse.SeekIDElement {
			fmt.Printf("%s- %v: %x\n", indent(info.Level), mkvparse.NameForElementID(id), value)
		} else {
			fmt.Printf("%s- %v: <binary>\n", indent(info.Level), mkvparse.NameForElementID(id))
		}
	*/
	return nil
}

func print(p MyParser) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"Number", "UID", "Name", "Language", "Codec", "Default"})

	for _, t := range p.tracks {
		tab.AppendRow([]interface{}{t.number, t.uid, t.name, t.language, t.CodecID, t.flagDefault})
	}
	tab.Render()
}

func main() {
	handler := MyParser{}
	err := mkvparse.ParsePath(os.Args[1], &handler)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
	print(handler)
}
