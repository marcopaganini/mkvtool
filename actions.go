// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
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

func runnerFromContext(ctx context.Context) *runner {
	ret, ok := ctx.Value(runnerKey).(*runner)
	if !ok {
		panic("internal error: Unable to retrieve runner from context.")
	}
	return ret
}

func actionMerge(c *cli.Context) error {
	return remux(c.Args().Slice(), c.String("output"), *runnerFromContext(c.Context), c.Bool("subs"))
}

func actionOnly(c *cli.Context) error {
	if err := checkTwoArgs(c); err != nil {
		return err
	}

	infile := c.Args().Get(0)
	outfile := c.Args().Get(1)

	run := *runnerFromContext(c.Context)

	mkv := mustParseFile(infile)
	tfi, err := extract(mkv, c.Int("track"), run)
	defer os.Remove(tfi.fname)
	if err != nil {
		return fmt.Errorf("%s: %v", infile, err)
	}
	return submux(infile, outfile, true, run)
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

func actionRemux(c *cli.Context) error {
	if err := checkTwoArgs(c); err != nil {
		return err
	}

	infile := c.Args().Get(0)
	outfile := c.Args().Get(1)
	run := *runnerFromContext(c.Context)

	return remux([]string{infile}, outfile, run, true)
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
		track, err := trackByLanguage(mkv, c.StringSlice("lang"), c.StringSlice("ignore"))
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
	for _, fname := range readable(c.Args().Slice()) {
		mkv := mustParseFile(fname)
		show(mkv, c.Bool("uid"))
	}
	return nil
}
