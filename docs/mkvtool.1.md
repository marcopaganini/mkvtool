# NAME

mkvtool

# SYNOPSIS

  **mkvtool [options] command [ARGs...]**

# DESCRIPTION

A swiss army knife for easy operations on Matroska containers.

This program provides a simpler front-end to common operations involving
Matroska containers. It hides details and complexities from utilities in the
mkvtoolnix suite, such as mkvmerge, mkvpropedit, and others.

# OPTIONS

  **--help**: Show context-sensitive help.

  **-n**, **--dry-run**: Dry-run mode (only show commands or output.)

# COMMANDS

## **help [COMMAND]**

Show a simple help page.

## **merge options [INPUT_FILEs]...**

Merge multiple input files (containing their respective media tracks) into
`<output-file>`.

This operation takes multiple input files (of any format recognized by
mkvmerge) and creates a new MKV file containing all tracks read from the input
files. A common use of this feature is to add a subtitle track to an existing
file that contains only A/V tracks (E.g., to create a MKV file with the A/V
tracks from an MP4 file and subtitles from a SRT file.)

**Options**:

* **--subs**: Copy all subtitles in `<input-file>` (default=true). Use
  `--nosubs` to replace all input subtitle tracks in a MKV file with one or
  more subtitle tracks coming from external files (E.g, SRT files).

  Note: The program does not check the contents of the input files. It is
  entirely possible to create files with multiple video, audio, and subtitle
  tracks (although we typically want only one video track in the output file.)

*  **-o OUTPUT_FILE, --output=OUTPUT_FILE**: Output file.

## **remux [options] INPUT_FILE OUTPUT_FILE**

Remux the input file into the output file, with the option to filter tracks.

**Options**:

* **--tracks**: Command separated list of raw track numbers to be copied. Use `mkvtool show` to
  see the full list of tracks with numbers. This option is mutually exclusive with `--lang`.

* **--lang**: Command separated list of track types/language specifiers to be
  copied.  The format the specifier is `[type]:language`, where type is "v" for
  video, "a" for audio, or "s" for subtitles and "language" is the three letter
  language code. The "language" field can also be **all**, meaning "copy all
  tracks of this type", **none**, for "copy no tracks of this type", and
  **first**, meaning "Only copy the first track of this type".

  If the lang specifier starts with a "!" the track will be excluded.

  By default, all tracks are copied (v:all,a:all,s:all), so you don't need to
  start the track list with "all".

* **--ignore**: Tracks matching the string specified by this option will not be
  copied to the output.  Useful to prevent the copy of Forced or SDH subtitles
  for example.

**Examples:**

* Copy all video tracks, English audio, and all subtitle tracks from in.mkv to
  out.mkv.  Note how by default, all tracks are copied so we don't need to use
  `v:all` and `a:all`:

```
mkvtool remux --lang='a:eng' in.mkv out.mkv
```

* Copy all video tracks, all audio tracks except Russian and French, and no subtitle tracks:

```
mkvtool remux --lang='!a:rus,!a:fre,s:none'
```

* Copy all tracks, but exclude any tracks with "Forced" in the description (as in "forced subtitles"):

```
mkvtool remux --ignore="Forced" in.mkv out.mkv
```

* Copy only all video tracks, the first audio track, and all subtitles in the English language:

```
mkvtool remux --lang='a:first,s:eng'
```

## **rename INPUT_FILEs...**

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

## **setdefault TRACK FILEs...**

Set the track specified with the `<track>` argument as the default track
in `<mkvfile>`.

Default tracks tell the player which track (in the presence of multiple tracks
of the same type) should be played. It's most commonly used for audio and
subtitle tracks. Not all players obey these settings (bust empirically, most
appear to do the right thing.)

## **setdefaultbylang [options] FILEs...**

Set the track with the first matching language as the default track.

* **--lang=LANGUAGE**: languages to match. It's possible to specify `--lang` multiple
  times, in which case the first match will prevail.

There's also the "default" meta-language, which matches the default
language in the MKV file (normally, shown as an empty language when
using the `show` option).

For example, using `--lang=en --lang=es --lang=default` will cause the program
to first attempt to find a subtitle track in the English language. If no
subtitle tracks in English exist, it will attempt Spanish next. As a last
resort, a subtitle track with the "default" language will be matched and set as
the default subtitle track (if it exists.)

* **--ignore=STRING**: Ignore any tracks containing the specified string in the
  description.  This is useful to ignore "Forced" subtitles, for example.

This command only works in subtitle tracks for now.

**Example:**

This will set the first subtitle track in English (change to your favorite
language) to be the default subtitle track. Failing that, it will attempt the
default language and "und" (some files have one subtitle track marked as
"undefined".)

```
mkvtool setdefaultbylang --lang=eng --lang=default --lang=und --ignore="force" *.mkv
```


## **show [options] INPUT_FILEs...**

Shows a listing of all tracks in the file.

Options:

* **-u, --uid**: Include track UIDs in the output (table output only).

* **-J, --json**: Output information in JSON format (as an array). When used,
  mkvtool prints full information about the files.

## **version**

Show version information.

# Author

- (C) 2021-2026 by Marco Paganini <paganini at paganini dot net>

