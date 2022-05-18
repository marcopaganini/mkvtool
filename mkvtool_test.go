// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"testing"
)

func TestFormat(t *testing.T) {
	casetests := []struct {
		fname     string
		mask      string
		want      string
		wantError bool
	}{
		// Basic case: Series name, season, spisode, year, resolution.
		{
			fname: "Series Title S01E02 HDTV x264 (2022) [1080p] FOOBAR.mkv",
			mask:  "%{title} %{season} %{episode} %{quality} %{codec} %{year} %{resolution}",
			want:  "Series Title 1 2 HDTV x264 2022 1080p",
		},
		// Formatting specifiers
		{
			fname: "Series Title S01E02 (2022) [1080p].mkv",
			mask:  "%{title} S%02.2{season}E%02.2{episode} [%{resolution}]",
			want:  "Series Title S01E02 [1080p]",
		},
		// Proper capitalization of title
		{
			fname: "a bad title that makes one of a kind 2022.mkv",
			mask:  "%{title} %{year}",
			want:  "A Bad Title That Makes One Of A Kind 2022",
		},
		// Invalid tag.
		{
			fname:     "Series Title S01E02 [1080p].mkv",
			mask:      "%{bad} S%02.2{season}E%02.2{episode} [%{resolution}]",
			wantError: true,
		},
		// Missing information in the parsed filename.
		{
			fname:     "Series Title S01E02 [1080p].mkv",
			mask:      "%{title} S%02.2{season}E%02.2{episode} (%{year}) [%{resolution}]",
			wantError: true,
		},
	}

	for _, tt := range casetests {
		got, err := format(tt.mask, tt.fname)
		if !tt.wantError {
			if err != nil {
				t.Fatalf("Got error %q want no error", err)
			}
			if got != tt.want {
				t.Fatalf("command diff: Got %v, want %v", got, tt.want)
			}
			continue
		}
		// Here, we want to see an error.
		if err == nil {
			t.Errorf("Got no error, want error")
		}
	}
}
