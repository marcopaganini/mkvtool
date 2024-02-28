// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Custom key for *runner in WithValue context.
type key int

const runnerKey = key(iota)

// readable returns a slice of readable files in the input slice.
func readable(fnames []string) []string {
	var ret []string

	for _, f := range fnames {
		if _, err := os.Stat(f); err == nil {
			ret = append(ret, f)
		} else {
			log.Printf("Note: File %q is not readable. Skipping.", f)
		}
	}
	return ret
}

func main() {
	var (
		// Command runner.
		runCmd runCommand

		// Dry-run command runner (only print commands).
		fakeRunCmd fakeRunCommand

		// This is overriden to fakeRunCmd when using dry-run.
		run runner = runCmd

		dryrun bool
	)

	if err := requirements(); err != nil {
		log.Fatalf("Requirements check: %v", err)
	}

	// Plain logs.
	log.SetFlags(0)

	app := &cli.App{
		Name: "mkvtool",
		Authors: []*cli.Author{
			{
				Name:  "Marco Paganini",
				Email: "paganini@paganini.net",
			},
		},
		Usage:   "Easy operations on Matroska containers.",
		Version: BuildVersion,

		// Global Flags
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"n"},
				Value:       false,
				Usage:       "Dry-run mode (only show commands)",
				Destination: &dryrun,
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, "")
			return nil
		},
		Before: func(c *cli.Context) error {
			// Run will resolve to a print-only version when dry-run is chosen.
			if dryrun {
				fmt.Println("Dry-run mode: Will not modify any files.")
				run = fakeRunCmd
				c.Context = context.WithValue(c.Context, runnerKey, &run)
			}
			return nil
		},
	}

	// Commands.
	app.Commands = []*cli.Command{
		// merge
		{
			Name:      "merge",
			Usage:     "Merge input tracks and files (A/V/S) into an output file",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "output",
					Aliases:  []string{"o"},
					Usage:    "Output file",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "subs",
					Usage: "Copy subtitles from original video file",
					Value: true,
				},
			},
			Action: actionMerge,
		},

		// only
		{
			Name:      "only",
			Usage:     "Remove all subtitle tracks, except one",
			ArgsUsage: "input_file output_file",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "track",
					Aliases:  []string{"t"},
					Usage:    "Track number to keep",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "subs",
					Usage: "Copy subtitles from original video file",
					Value: true,
				},
			},
			Action: actionOnly,
		},

		// print
		{
			Name:      "print",
			Usage:     "Parse input filename and print scene information using a printf style mask.",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"f"},
					Value:   "%{title}.mkv",
					Usage:   "Formating mask",
				},
			},
			Action: actionPrint,
		},

		// remux
		{
			Name:      "remux",
			Usage:     "Remux input file into an output file",
			ArgsUsage: "input_file output_file",
			Action:    actionRemux,
		},

		// rename
		{
			Name:      "rename",
			Usage:     "Rename file based on scene information in filename.",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"f"},
					Value:   "%{title}.%{container}",
					Usage:   "Formating mask",
				},
			},
			Action: actionRename,
		},

		// setdefault
		{
			Name:      "setdefault",
			Usage:     "Set the default subtitle tag on a track.",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "track",
					Aliases:  []string{"t"},
					Usage:    "Track Number",
					Required: true,
				},
			},
			Action: actionSetDefault,
		},

		// setdefaultbylang
		{
			Name:      "setdefaultbylang",
			Usage:     "Set default subtitle track by language.",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "lang",
					Aliases:  []string{"l"},
					Usage:    "Preferred languages (Use multiple times. Use 'default' for tracks with no language set.)",
					Required: true,
				},
				&cli.StringSliceFlag{
					Name:    "ignore",
					Aliases: []string{"i"},
					Usage:   "Ignore tracks with this string in the name (can be used multiple times.)",
				},
			},
			Action: actionSetDefaultByLang,
		},

		// show
		{
			Name:      "show",
			Usage:     "Show information about files",
			ArgsUsage: "FILE(s)...",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "uid",
					Aliases: []string{"u"},
					Usage:   "Include track UIDs in the output",
				},
			},
			Action: actionShow,
		},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, runnerKey, &run)
	err := app.RunContext(ctx, os.Args)

	if err != nil {
		log.Fatalln("Execution failed:", err)
	}
}
