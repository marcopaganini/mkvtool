package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app    = kingpin.New("subtool", "Subtitle operations on matroska containers.")
		dryrun = app.Flag("dry-run", "Dry-run mode (only show commands).").Short('n').Bool()

		// merge
		mergeCmd    = app.Command("merge", "Merge input tracks and files (subtitle/video/audio) into an output file.")
		mergeOutput = mergeCmd.Flag("output", "Output file.").Required().Short('o').String()
		mergeSubs   = mergeCmd.Flag("subs", "Copy subs from video file.").Default("true").Bool()
		mergeInputs = mergeCmd.Arg("input-files", "Input files.").Required().Strings()

		// only
		onlyCmd       = app.Command("only", "Remove all subtitle tracks, except one.")
		setOnlyTrack  = onlyCmd.Arg("track", "Track number to keep.").Required().Int()
		setOnlyInput  = onlyCmd.Arg("input", "Matroska Input file.").Required().String()
		setOnlyOutput = onlyCmd.Arg("output", "Matroska Output file.").Required().String()

		// remux
		remuxCmd       = app.Command("remux", "Remux input file into an output file.")
		remuxCmdInput  = remuxCmd.Arg("input-file", "Matroska Input file.").Required().String()
		remuxCmdOutput = remuxCmd.Arg("output-file", "Matroska Output file.").Required().String()

		// rename
		renameCmd   = app.Command("rename", "Rename file based on scene information in filename.")
		renameFiles = renameCmd.Arg("input-files", "Matroska file(s).").Required().Strings()

		// setdefault
		setDefaultCmd   = app.Command("setdefault", "Set default subtitle tag on a track.")
		setDefaultTrack = setDefaultCmd.Arg("track", "Track number to set as default.").Required().Int()
		setDefaultFiles = setDefaultCmd.Arg("mkvfile", "Matroska file.").Required().Strings()

		// setdefaultbylanguage
		setDefaultByLangCmd    = app.Command("setdefaultbylang", "Set default subtitle track by language.")
		setDefaultByLangList   = setDefaultByLangCmd.Flag("lang", "Preferred languages (Use multiple times. Use 'default' for tracks with no language set.)").Required().Strings()
		setDefaultByLangIgnore = setDefaultByLangCmd.Flag("ignore", "Ignore tracks with this string in the name (can be used multiple times.)").Strings()
		setDefaultByLangFiles  = setDefaultByLangCmd.Arg("mkvfiles", "Matroska file(s).").Required().Strings()

		// show
		showCmd   = app.Command("show", "Show Information about file(s).")
		showUID   = showCmd.Flag("uid", "Include track UIDs in the output.").Short('u').Bool()
		showFiles = showCmd.Arg("input-files", "Matroska Input files.").Required().Strings()

		// version
		versionCmd = app.Command("version", "Show version information.")

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
	// Just print version number and exit.
	case versionCmd.FullCommand():
		fmt.Printf("Build Version: %s\n", BuildVersion)
		os.Exit(0)

	case mergeCmd.FullCommand():
		err = remux(*mergeInputs, *mergeOutput, run, *mergeSubs)

	case onlyCmd.FullCommand():
		mkv := mustParseFile(*setOnlyInput)
		var tfi trackFileInfo
		tfi, err = extract(mkv, *setOnlyTrack, run)
		if err != nil {
			break
		}
		err = submux(*setOnlyInput, *setOnlyOutput, true, run, tfi)
		// Attempt to remove even on error.
		_ = os.Remove(tfi.fname)

	case remuxCmd.FullCommand():
		err = remux([]string{*remuxCmdInput}, *remuxCmdOutput, run, true)

	case renameCmd.FullCommand():
		for _, f := range *renameFiles {
			rename(f, *dryrun)
		}

	case setDefaultCmd.FullCommand():
		for _, f := range *setDefaultFiles {
			mkv := mustParseFile(f)
			err = setdefault(mkv, *setDefaultTrack, run)
			if err != nil {
				break
			}
		}

	case setDefaultByLangCmd.FullCommand():
		for _, f := range *setDefaultByLangFiles {
			mkv := mustParseFile(f)
			var track int
			track, err = trackByLanguage(mkv, *setDefaultByLangList, *setDefaultByLangIgnore)
			if err != nil {
				break
			}
			err = setdefault(mkv, track, run)
		}

	case showCmd.FullCommand():
		for _, f := range *showFiles {
			mkv := mustParseFile(f)
			show(mkv, *showUID)
		}
	}

	// Print error message, if any
	if err != nil {
		log.Fatalf("Error during %s: %v", k, err)
	}
}
