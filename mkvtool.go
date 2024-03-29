// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/jedib0t/go-pretty/table"
	ParseTorrentName "github.com/middelink/go-parse-torrent-name"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// A friendly chat about Matroska metadata track numbers.
//
// Matroska tracks numbers are confusing. Tracks are stored in the file
// starting at 1 (ONE). Some mkvtoolnix commands such as mkvmerge and
// mkvextract expect tracks to start at offset zero (ZERO), while others like
// mkvpropedit, expect offset 1. Due to this, the following conventions were
// adopted here:
//
// - Tracks are always displayed starting at 0 (as the output of mkvmerge --identify)
// - Any actions using mkvpropedit automatically add one to the track number.
// - Any actions using mkvmerge or mkvextract will use the track number unchanged.
//
// Track Types. See https://www.matroska.org/technical/specs/index.html
const (
	typeSubtitle = "subtitles"
)

// trackFileInfo holds information about an exported track file.
type trackFileInfo struct {
	language string
	fname    string
}

// BuildVersion holds the git build number (set by make).
var BuildVersion string

// show lists all tracks in a file.
func show(mkv matroska, showUID bool) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	if showUID {
		tab.AppendHeader(table.Row{"Number", "UID", "Type", "Name", "Language", "Codec", "Default"})
	} else {
		tab.AppendHeader(table.Row{"Number", "Type", "Name", "Language", "Codec", "Default"})
	}

	for _, track := range mkv.Tracks {
		// Create a row with the desired columns.
		// mkvmerge reports tracks starting at zero, so we add one to match the file.
		row := []interface{}{track.ID}
		if showUID {
			row = append(row, uint64(track.Properties.UID))
		}
		row = append(row, track.Type, track.Properties.TrackName, track.Properties.Language, track.Codec)

		// Make default flag easier to see.
		if track.Properties.DefaultTrack {
			row = append(row, "<=====")
		} else {
			row = append(row, "")
		}
		tab.AppendRow(row)
	}
	tab.Render()
}

// setdefault resets flagDefault on all subtitle tracks and sets it on the chosen track UID.
func setdefault(mkv matroska, tracknum int, cmd runner) error {
	command := []string{
		"mkvpropedit",
		mkv.FileName,
	}

	for _, track := range mkv.Tracks {
		if track.Type == typeSubtitle {
			// mkvpropedit uses base 1 for track (not zero).
			command = append(command, "--edit", fmt.Sprintf("track:%d", track.ID+1), "--set", "flag-default=0")
		}
	}

	if err := cmd.run(command[0], command[1:]...); err != nil {
		return err
	}
	return adddefault(mkv, tracknum, cmd)
}

// trackByLanguage returns the track number (base 0) for the first track with
// one of the specified languages. The list of languages works as a priority,
// meaning that languages=["eng","fra"] will first attempt to find a track with
// the English language, and failing that, French. The special language
// "default" will cause tracks without a language code to be selected (Matroska
// has the concept of a "default language", which is usually English -- tracks
// with this language will not have a language code).
//
// The ignore slice contains a list of strings for case-insentive search
// against the track name. If the selected language contains one of the strings
// in this slice, it will be ignored. This is useful to select tracks by
// language while ignoring 'Forced' tracks.
func trackByLanguage(mkv matroska, languages []string, ignore []string) (int, error) {
	for _, lang := range languages {
		if lang == "default" {
			lang = ""
		}
		for _, track := range mkv.Tracks {
			// Match subtitle and language.
			if track.Type != typeSubtitle || track.Properties.Language != lang {
				continue
			}
			// Make sure track should not be ignored.
			if stringInSlice(track.Properties.TrackName, ignore) {
				continue
			}
			return track.ID, nil
		}
	}
	return 0, fmt.Errorf("no track with language(s): %s", strings.Join(languages, ","))
}

// stringInSlice returns true if a string exists inside a slice of strings.
// Comparison is case insensitive.
func stringInSlice(s string, slc []string) bool {
	for _, substr := range slc {
		if strings.Contains(strings.ToLower(s), strings.ToLower(substr)) {
			return true
		}
	}
	return false
}

// extract extracts a given track into a file.
func extract(mkv matroska, tracknum int, cmd runner) (trackFileInfo, error) {
	// Fetch language for the track. Fail if track does not exist.
	ok := false
	language := ""
	for _, track := range mkv.Tracks {
		if track.ID == tracknum {
			ok = true
			language = track.Properties.Language
			break
		}
	}
	if !ok {
		return trackFileInfo{}, fmt.Errorf("track #%d not found in file %s", tracknum, mkv.FileName)
	}

	// Extract into a temporary file
	tmpfile, err := ioutil.TempFile("", "mkvtool")
	if err != nil {
		return trackFileInfo{}, err
	}
	temp := tmpfile.Name()
	_ = tmpfile.Close()

	command := []string{
		"mkvextract",
		mkv.FileName,
		"tracks",
		fmt.Sprintf("%d:%s", tracknum, temp),
	}
	if err := cmd.run(command[0], command[1:]...); err != nil {
		return trackFileInfo{}, err
	}
	return trackFileInfo{language: language, fname: temp}, nil
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

// remux re-multiplexes the input file(s) into the output file. Setting subs to
// false will cause subs not to be copied.
func remux(infiles []string, outfile string, cmd runner, subs bool) error {
	cmdline := []string{"mkvmerge"}
	if !subs {
		cmdline = append(cmdline, "-S")
	}
	cmdline = append(cmdline, infiles...)
	cmdline = append(cmdline, "-o", outfile)

	return cmd.run(cmdline[0], cmdline[1:]...)
}

// adddefault adds the default flag to a given track UID.
func adddefault(mkv matroska, tracknum int, cmd runner) error {
	for _, track := range mkv.Tracks {
		if track.ID == tracknum {
			// mkvpropedit uses base 1 for tracks.
			return cmd.run("mkvpropedit", mkv.FileName, "--edit", fmt.Sprintf("track:%d", tracknum+1), "--set", "flag-default=1")
		}
	}
	return fmt.Errorf("file %s does not contain track %d", mkv.FileName, tracknum)
}

// rename renames a file according to the "Scene" information in the file.
func rename(mask, fname string, dryrun bool) error {
	newname, err := format(fname, mask)
	if err != nil {
		return err
	}
	dir, _ := filepath.Split(fname)
	newfile := filepath.Join(dir, newname)

	fmt.Printf("%s => %s\n", fname, newfile)
	if dryrun {
		return nil
	}
	return os.Rename(fname, newfile)
}

// format parses "Scene" information in the file and returns a string formatted
// according to a formatting mask. The mask may contain the following tokens:
//
// %[format]{audio}
// %[format]{codec}
// %[format]{container} (this matches the original extension)
// %[format]{episode}
// %[format]{episodeName}
// %[format]{excess}
// %[format]{extended}
// %[format]{garbage}
// %[format]{group}
// %[format]{hardcoded}
// %[format]{language}
// %[format]{proper}
// %[format]{quality}
// %[format]{region}
// %[format]{repack}
// %[format]{resolution}
// %[format]{season}
// %[format]{title}
// %[format]{website}
// %[format]{widescreen}
// %[format]{year}
//
// Where "format" is a printf style format sizing specification. Complete
// examples:
// - %-02.2{season} - Season formatted as two characters, left padded wth zeroes.
// - %-20{title} - Title truncated to 20 characters
//
// Anything not matching the %[format]{xxxx} construct will be interpreted literally.
//
// Formatting will fail if any element present in the mask cannot be resolved
// (a typical example is asking for episode numbers for movies).
func format(mask, fname string) (string, error) {
	// Split the filename so we can work on parts separately.
	_, file := filepath.Split(fname)

	parsed, err := ParseTorrentName.Parse(file)
	if err != nil {
		return "", err
	}
	fields := structs.Map(parsed)

	// tags are formatted as %[format]{value}
	re, err := regexp.Compile(`%((?:-?[\d]+)?(?:\.\d+)?){([a-z]+)}`)
	if err != nil {
		return "", err
	}

	var errlist []string

	formatted := re.ReplaceAllStringFunc(mask, func(match string) string {
		// Split matched tag into size formatting specifier and tag name.
		// Tag must be capitalized to match the keys in the map.
		e := re.FindStringSubmatch(match)
		sizespec := e[1]
		tag := cases.Title(language.English).String(e[2])

		if i, ok := fields[tag]; ok {
			switch t := i.(type) {
			case string:
				val := i.(string)
				if val == "" {
					break
				}
				// Special case for title: Capitalize
				if tag == "Title" {
					val = cases.Title(language.English).String(val)
				}
				return fmt.Sprintf("%"+sizespec+"s", val)
			case int:
				val := i.(int)
				if val <= 0 {
					break
				}
				return fmt.Sprintf("%"+sizespec+"d", i.(int))
			default:
				errlist = append(errlist, fmt.Sprintf("Internal error: Unknown type %T for %q", match, t))
				return "*ERROR*"
			}
		}
		errlist = append(errlist, fmt.Sprintf("Unable to parse data for %s", match))
		return "*ERROR*"
	})

	if len(errlist) != 0 {
		return "", fmt.Errorf("%s", strings.Join(errlist, ";"))
	}
	return formatted, nil
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

// mustParseFile parses the MKV file using the JSON output from mkmerge --identify.
// error message in case of problems.
func mustParseFile(fname string) matroska {
	var stdout bytes.Buffer

	cmd := exec.Command("mkvmerge", "--identify", "-F", "json", fname)
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Printf("mkvmerge failed: %v", err)
		log.Fatalf("--- Output ---\n%s\n", stdout.String())
	}

	// Decode JSON.
	var mkv matroska
	err = json.Unmarshal(stdout.Bytes(), &mkv)
	if err != nil {
		log.Fatalf("Error decoding JSON output from mkvmerge: %v", err)
	}
	return mkv
}
