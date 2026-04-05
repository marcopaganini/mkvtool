// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// errorFromSlice converts the slice of strings into a single multi-line string
// and returns it, or returns nil if the error list is empty.
func errorFromSlice(errmsgs []string) error {
	if len(errmsgs) != 0 {
		return errors.New(strings.Join(errmsgs, "\n"))
	}
	return nil
}

// checkTwoArgs shows the help message for the current context and return an
// error if we don't have exactly two arguments.
func checkTwoArgs(c *cli.Context) error {
	if c.Args().Len() != 2 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return errors.New("need input and output files")
	}
	return nil
}

// checkMultiArgs shows the help message for the current context and return an
// error if we don't have at least one argument.
func checkMultiArgs(c *cli.Context) error {
	if c.Args().Len() < 1 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return errors.New("no files to process")
	}
	return nil
}

// checkTrackType checks the validity of a track type. It returns an error if
// the specified type is invalid, or the full name of the type
// (audio/video/subtitles).
func checkTrackType(t string) (string, error) {
	tracktypes := map[string]string{
		"a":         "audio",
		"v":         "video",
		"s":         "subtitles",
		"aud":       "audio",
		"vid":       "video",
		"sub":       "subtitles",
		"audio":     "audio",
		"video":     "video",
		"subtitles": "subtitles",
	}
	if name, ok := tracktypes[t]; ok {
		return name, nil
	}
	return "", fmt.Errorf("invalid track type (use a, v, or s): %v", t)
}

func runnerFromContext(ctx context.Context) *runner {
	ret, ok := ctx.Value(runnerKey).(*runner)
	if !ok {
		panic("internal error: Unable to retrieve runner from context.")
	}
	return ret
}

func actionMerge(c *cli.Context) error {
	outfile := c.String("output")
	if err := remux(c.Args().Slice(), outfile, *runnerFromContext(c.Context), c.Bool("subs")); err != nil {
		return err
	}

	if !c.Bool("quiet") {
		if c.Bool("json") {
			return showJSON([]matroska{mustParseFile(outfile)})
		}
		show(mustParseFile(outfile), false)
	}
	return nil
}

// actionRemux remuxes the input file into the output file, filtering the
// tracks using the --tracks and --lang options. Tracks can be specified
// directly using the track numbers (separated by comma), or by type/language
// by using the "--lang" option.
//
// For the --lang option, format the specifier as a comma separated list of
// "[type]:language" items, where type is "v" for video, "a" for audio, or "s" for
// subtitles and "language" is the three letter language code.
//
// You can also use special idenfifiers for the languages:
// - all - copy all tracks having this type.
// - none - copy no tracks having this type.
// - first - copy only the first track having this type.
//
// If the lang specifier starts with a "!" the track will be excluded.
//
// By default, all tracks are copied (v:all,a:all,s:all), so you don't need to start
// the track list with "all".
//
// Examples:
//	 # Copy tracks 1, 2, and 3 (independent of type).
//   mkvtool remux input.mkv output.mkv --tracks 1,2,3
//
//  # Copy all video tracks, all english audio tracks, and all subtitles except for russian.
//  mkvtool remux input.mkv output.mkv --lang 'v:eng,a:rus,!s:rus'

func actionRemux(c *cli.Context) error {
	if err := checkTwoArgs(c); err != nil {
		return err
	}

	infile := c.Args().Get(0)
	outfile := c.Args().Get(1)
	run := *runnerFromContext(c.Context)

	mkv := mustParseFile(infile)

	// Keep all tracks by default.
	keep := make(map[int]bool)
	for _, t := range mkv.Tracks {
		keep[t.ID] = true
	}

	// Tracks by type.
	tracksByType := map[string][]int{
		typeVideo:     {},
		typeAudio:     {},
		typeSubtitles: {},
	}
	for _, t := range mkv.Tracks {
		if _, ok := tracksByType[t.Type]; ok {
			tracksByType[t.Type] = append(tracksByType[t.Type], t.ID)
		}
	}

	// Process the lang list.
	if c.IsSet("lang") {
		typeMap := map[string]string{
			"v": typeVideo,
			"a": typeAudio,
			"s": typeSubtitles,
		}
		// seenInclusion tracks if we've seen an inclusion rule for each type.
		// The first inclusion rule for a type clears all tracks of that type.
		seenInclusion := make(map[string]bool)

		for _, part := range strings.Split(c.String("lang"), ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			exclude := strings.HasPrefix(part, "!")
			if exclude {
				part = part[1:]
			}

			// Split into type:lang.
			p := strings.SplitN(part, ":", 2)
			if len(p) != 2 {
				continue
			}
			t, lang := p[0], p[1]
			fullType, ok := typeMap[t]
			if !ok {
				continue
			}

			// If this is the first inclusion rule for this type, clear all tracks of this type.
			if !exclude && !seenInclusion[fullType] {
				for _, id := range tracksByType[fullType] {
					delete(keep, id)
				}
				seenInclusion[fullType] = true
			}

			// Apply the rule.
			switch lang {
			case "all":
				for _, id := range tracksByType[fullType] {
					keep[id] = true
					if exclude {
						delete(keep, id)
					}
				}
			case "none":
				// only makes sense for inclusion, as exclusion of none is a no-op.
				if !exclude {
					for _, id := range tracksByType[fullType] {
						delete(keep, id)
					}
				}
			case "first":
				if len(tracksByType[fullType]) > 0 {
					id := tracksByType[fullType][0]
					keep[id] = true
					if exclude {
						delete(keep, id)
					}
				}
			default:
				// Match language.
				if lang == "default" {
					lang = ""
				}
				for _, track := range mkv.Tracks {
					if track.Type == fullType && track.Properties.Language == lang {
						keep[track.ID] = true
						if exclude {
							delete(keep, track.ID)
						}
					}
				}
			}
		}
	}

	// Process the tracks list (comma separated IDs).
	if c.IsSet("tracks") {
		for _, part := range strings.Split(c.String("tracks"), ",") {
			var id int
			if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &id); err == nil {
				keep[id] = true
			}
		}
	}

	// Apply ignore flag (supersedes selection).
	if c.IsSet("ignore") {
		for _, track := range mkv.Tracks {
			if stringInSlice(track.Properties.TrackName, c.StringSlice("ignore")) {
				delete(keep, track.ID)
			}
		}
	}

	// Final list of tracks to keep, separated by type.
	keepList := map[string][]string{
		typeVideo:     {},
		typeAudio:     {},
		typeSubtitles: {},
	}
	for _, t := range mkv.Tracks {
		if keep[t.ID] {
			keepList[t.Type] = append(keepList[t.Type], fmt.Sprintf("%d", t.ID))
		}
	}

	// Construct mkvmerge command.
	cmdline := []string{"mkvmerge", "-o", outfile}

	// Add track selection flags.
	flags := map[string]struct {
		include string
		exclude string
	}{
		typeVideo:     {include: "-d", exclude: "-D"},
		typeAudio:     {include: "-a", exclude: "-A"},
		typeSubtitles: {include: "-s", exclude: "-S"},
	}

	// Order: video, audio, subtitles (not strictly required, but for consistency).
	for _, t := range []string{typeVideo, typeAudio, typeSubtitles} {
		if len(keepList[t]) == 0 {
			cmdline = append(cmdline, flags[t].exclude)
		} else {
			cmdline = append(cmdline, flags[t].include, strings.Join(keepList[t], ","))
		}
	}

	cmdline = append(cmdline, infile)

	if err := run.run(cmdline[0], cmdline[1:]...); err != nil {
		return err
	}

	if !c.Bool("quiet") {
		if c.Bool("json") {
			return showJSON([]matroska{mustParseFile(outfile)})
		}
		show(mustParseFile(outfile), false)
	}
	return nil
}

func actionPrint(c *cli.Context) error {
	if err := checkMultiArgs(c); err != nil {
		return err
	}

	var errmsgs []string

	for _, fname := range c.Args().Slice() {
		output, err := format(c.String("format"), fname)
		if err != nil {
			errmsgs = append(errmsgs, fmt.Sprintf("%s: %v", fname, err))
			continue
		}
		fmt.Println(output)
	}
	return errorFromSlice(errmsgs)
}

func actionRename(c *cli.Context) error {
	if err := checkMultiArgs(c); err != nil {
		return err
	}

	var errmsgs []string

	for _, fname := range readable(c.Args().Slice()) {
		err := rename(c.String("format"), fname, c.Bool("dry-run"))
		if err != nil {
			errmsgs = append(errmsgs, fmt.Sprintf("%s: %v", fname, err))
		}
	}
	return errorFromSlice(errmsgs)
}

func actionSetDefault(c *cli.Context) error {
	if err := checkMultiArgs(c); err != nil {
		return err
	}

	run := *runnerFromContext(c.Context)

	var errmsgs []string

	for _, fname := range readable(c.Args().Slice()) {
		mkv := mustParseFile(fname)
		err := setdefault(mkv, c.Int("track"), run)
		if err != nil {
			errmsgs = append(errmsgs, fmt.Sprintf("%s: %s", fname, err))
		}
	}
	return errorFromSlice(errmsgs)
}

func actionSetDefaultByLang(c *cli.Context) error {
	if err := checkMultiArgs(c); err != nil {
		return err
	}

	run := *runnerFromContext(c.Context)

	var errmsgs []string

	for _, fname := range readable(c.Args().Slice()) {
		mkv := mustParseFile(fname)
		track, err := trackByLanguageAndType(mkv, c.StringSlice("lang"), typeSubtitles, c.StringSlice("ignore"))
		if err != nil {
			errmsgs = append(errmsgs, fmt.Sprintf("%s: %v", fname, err))
			continue
		}
		err = setdefault(mkv, track, run)
		if err != nil {
			errmsgs = append(errmsgs, fmt.Sprintf("%s: %v", fname, err))
		}
	}
	return errorFromSlice(errmsgs)
}

func actionShow(c *cli.Context) error {
	if err := checkMultiArgs(c); err != nil {
		return err
	}

	var allMKV []matroska
	isJSON := c.Bool("json")

	for _, fname := range readable(c.Args().Slice()) {
		mkv := mustParseFile(fname)
		if isJSON {
			allMKV = append(allMKV, mkv)
		} else {
			show(mkv, c.Bool("uid"))
		}
	}

	if isJSON {
		return showJSON(allMKV)
	}
	return nil
}
