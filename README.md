# mkvtool

## Description

A handy frontend CLI to many common operations on Matroska (MKV) files using
[mkvtoolnix](https://mkvtoolnix.download/) (and some extras). I wrote this program after
spending considerable time with the (necessarily) complex command-line options
of the mkvtoolnix toolkit.

## Installation

### Pre-requisites

* `mkvtoolnix` (`sudo apt-get install mkvtoolnix` on debian based systems): **required**.
* `pandoc` (`sudo apt-get install pandoc` on debian based systems): **optional**.
  Only required to generate the man page).

### Download a binary version from github

Github automatically builds new versions for multiple platforms. You can download
a binary (and manpage) for your platform from the [mkvtool releases page](https://github.com/marcopaganini/mkvtool/releases).

Once downloaded, use the following commands to install it in your system:

```
$ ARCH="your architecture name" # (E.g. amd64)
$ sudo mkdir -p /usr/local/bin
$ sudo mkdir -p /usr/local/man/man1
$ sudo install -m 755 "$ARCH/mkvtool" /usr/local/bin
$ sudo install -m 644 "$ARCH/mkvtool.1" /usr/local/man/man1
```

Run `mkvtool` for a short help message or `man mkvtool` for the full man page.

### Clone and compile

This requires a recent version of Go Installed in your system. You also need
`pandoc` if you want to recreate the manpage.

```
$ git clone https://github.com/marcopaganini/mkvtool
$ cd mkvtool
$ make
$ make manpage # <-- optional
$ sudo make install
```

## Documentation

Documentation is [available online](docs/mkvtool.1.md) or by typing `man mkvtool`
in your shell prompt.

## Author

Marco Paganini <paganini@paganini.net>.

Comments, ideas, and PRs are welcome.
