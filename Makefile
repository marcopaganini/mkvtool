.PHONY: arch clean install manpage

bin := mkvtool
bindir := /usr/local/bin
mandir := /usr/local/man/man1
archdir := arch
src := $(wildcard *.go)
git_tag := $(shell git describe --always --tags)

# Default target
${bin}: Makefile ${src}
	CGO_ENABLED=0 go build -v -ldflags "-X main.BuildVersion=${git_tag}" -o "${bin}"

manpage:
	pandoc --metadata title="${bin}(1)" -s --from gfm-smart --to man "docs/${bin}.1.md" -o "docs/${bin}.1"

clean:
	rm -f "${bin}"
	rm -f "docs/${bin}.1"
	rm -rf "${archdir}"

install: ${bin} ${manpage}
	install -m 755 "${bin}" "${bindir}"
	mkdir -p "${mandir}"
	install -m 644 "docs/${bin}.1" "${mandir}"

# Creates cross-compiled tarred versions (for releases).
arch: Makefile ${src} manpage
	for ga in "linux/amd64" "linux/386" "linux/arm" "linux/arm64" "linux/mips" "linux/mipsle"; do \
	  export GOOS="$${ga%/*}"; \
	  export GOARCH="$${ga#*/}"; \
	  dst="./${archdir}/$${GOOS}-$${GOARCH}"; \
	  mkdir -p "$${dst}"; \
	  echo "=== Building $${GOOS}/$${GOARCH} ==="; \
	  go build -v -ldflags "-X main.Build=${git_tag}" -o "$${dst}/${bin}"; \
	  [ -s LICENSE ] && install -m 644 LICENSE "$${dst}"; \
	  [ -s README.md ] && install -m 644 README.md "$${dst}"; \
	  [ -s docs/${bin}.1 ] && install -m 644 docs/${bin}.1 "$${dst}"; \
	  tar -C "${archdir}" -zcvf "${archdir}/${bin}-$${GOOS}-$${GOARCH}.tar.gz" "$${dst##*/}"; \
	done
