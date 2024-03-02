# NAME

mkvtool

# SYNOPSIS

  **mkvtool [\<flags\>] \<command\> [\<args\> ...]**

# DESCRIPTION

A swiss army knife for easy operations on Matroska containers.

This program provides a simpler front-end to common operations involving
Matroska containers. It hides details and complexities from utilities in the
mkvtoolnix suite, such as mkvmerge, mkvpropedit, and others.

# OPTIONS

  **--help**: Show context-sensitive help (also try `--help-long` and `--help-man`).

  **-n**, **--dry-run**: Dry-run mode (only show commands or output.)

# COMMANDS

## **help [\<command\>...]**

Show help.

## **merge --output=OUTPUT [\<flags\>] \<input-files\>...**

Merge multiple input files (containing their respective media tracks) into
`<output-file>`.

This operation takes multiple input files (of any format recognized by
mkvmerge) and creates a new MKV file containing all tracks read from the input
files. A common use of this feature is to add a subtitle track to an existing
file that contains only A/V tracks (E.g., to create a MKV file with the A/V
tracks from an MP4 file and subtitles from a SRT file.)

The boolean flag `--subs` causes the program to copy all subtitles in
`<input-file>` (default=true). Use `--nosubs` to replace all input subtitle
tracks in a MKV file with one or more subtitle tracks coming from external
files (E.g, SRT files).

Note: The program does not check the contents of the input files. It is
entirely possible to create files with multiple video, audio, and
subtitle tracks (although we typically want only one video track in the
output file.)

  **-o, --output=OUTPUT**: Output file.

  **--subs**:  Copy subs from video file (use `--nosubs` to ignore all subs in the source file.)

## **remux [\<flags\>] \<input-file\> \<output-file\>**

Remux the input file into the output file, filtering the tracks using the
`--tracks` and `--lang` options.

For the `--lang` option, format the specifier as a comma separated list of
`[type]:language` items, where type is "v" for video, "a" for audio, or "s" for
subtitles and "language" is the three letter language code.

You can also use special identifiers for the languages:
- **all**: copy all tracks having this type.
- **none**: copy no tracks having this type.
- **first**: copy only the first track having this type.

If the lang specifier starts with a "!" the track will be excluded.

By default, all tracks are copied (v:all,a:all,s:all), so you don't need to
start the track list with "all".

The `--ignore` flag supersedes any selection; tracks matching the ignore
strings will be removed even if explicitly selected.

  **-i, --ignore=IGNORE**: Ignore tracks with this string in the name (can be used multiple times.)
  **-l, --lang=LANG**: Language selection by type. E.g: v:all,a:eng,!s:rus
  **-t, --tracks=TRACKS**: Track number(s) to keep (comma separated list).

## **rename \<input-files\>...**

Rename `<input-files>` into a standardized format, using metadata in the
original filenames.

Attempt to fetch metadata (title, season, year, resolution, etc) from each of
the filenames specified in `<input-files>`. The filenames themselves are parsed
for data, so this is a "best effort" operation.  Each filename is renamed
following a simplified format, compatible with most streaming applications.

Note: For TV series, There's no clean way to fetch the episode name, so this
information will be lost during the rename. Most streaming servers fill that
information in their databases based on Title, Episode, and Season, so that
tends not to be a problem for most people.

## **setdefault \<track\> \<mkvfile\>...**

Set the track specified with the `<track>` argument as the default track
in `<mkvfile>`.

Default tracks tell the player which track (in the presence of multiple tracks
of the same type) should be played. It's most commonly used for audio and
subtitle tracks. Not all players obey these settings (bust empirically, most
appear to do the right thing.)

## **setdefaultbylang --lang=LANG [\<flags\>] \<mkvfiles\>...**

Set the track with the first matching language as the default track.

This option sets the default track based on the track language. Use the
`--lang` option to select a set of languages to match. It's possible to specify
`--lang` multiple times, in which case the first match will prevail.

There's also the "default" meta-language, which matches the default
language in the MKV file (normally, shown as an empty language when
using the `show` option).

For example, using `--lang=en --lang=es --lang=default` will cause the program
to first attempt to find a subtitle track in the English language. If no
subtitle tracks in English exist, it will attempt Spanish next. As a last
resort, a subtitle track with the "default" language will be matched and set as
the default subtitle track (if it exists.)

It's also possible to ignore strings in the title with the `--ignore` flag.
This is useful to ignore "Forced" subtitles, for example.

This command only works in subtitle tracks for now.

Useful command example:

```
$ mkvtool setdefaultbylang --lang=eng --lang=default --lang=und --ignore="force" *.mkv
```

This will set the first subtitle track in English (change to your favorite
language) to be the default subtitle track. Failing that, it will attempt the
default language and "und" (some files have one subtitle track marked as
"undefined".)

  **--lang=LANG**: Preferred languages (Use multiple times. Use 'default' for
    tracks with no language set.)

  **--ignore=IGNORE**: Ignore tracks with this string in the name (can be
    used multiple times.)

## **show \[\<flags\>\] \<input-files\>...**

Shows a listing of all tracks in the file.

  **-u, --uid**: Include track UIDs in the output.

## **version**

Show version information.

# Author

- (C) 2021-2026 by Marco Paganini <paganini at paganini dot net>

