// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type runner interface {
	run(string, ...string) error
}

// runner provides a simple and mockable interface to exec.Command()
type runCommand int

// run creates an *exec.Cmd object using exec.Command and runs
// it using exec.Run. The return is the return of exec.Run.
func (x runCommand) run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, _ = io.Copy(os.Stdout, stdout)
	_, _ = io.Copy(os.Stderr, stderr)

	return cmd.Wait()
}

// fakeRunCommand provides a runner for dry-run operations.
type fakeRunCommand int

// Fakerunner just logs the commands (dry-run)
func (x fakeRunCommand) run(name string, args ...string) error {
	var quoted []string

	for _, a := range args {
		quoted = append(quoted, strconv.Quote(a))
	}
	log.Printf("%q %s", name, strings.Join(quoted, " "))
	return nil
}
