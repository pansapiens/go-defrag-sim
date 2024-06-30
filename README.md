# Disk Defragmentation Simulation

A toy visualization inspired by the [hypnotic DOS disk defragmentation program](https://youtu.be/lxZyxxHOw3Y?si=7IolvOg4tHWxgn-2&t=148), `scandisk.exe`.
I've tried to give it similar aesthetics, but haven't attempted historical or algorithmic accuracy.  
I don't think it quite captures the magic of the original, but here it is.

> Doesn't actually scan or modify your disk !

Written in collaboration with GPT-4o.

Quit with `q` or Ctrl-C.

## Usage

```
$ ./scandisk.exe --help

Usage of ./scandisk.exe:
  -d, --delay int    Base delay in milliseconds (default 100)
  -h, --height int   Height of the terminal
  -l, --loop         Run in an infinite loop
      --skip-scan    Skip the scanning step
  -w, --width int    Width of the terminal
pflag: help requested
```

Run in a 80x20 terminal, forever:

```sh
./scandisk.exe --loop --width 80 --height 20
```

## Building

### Requirements

- Go 1.16 or higher
- Packages:
    - `github.com/gdamore/tcell/v2`
    - `github.com/spf13/pflag`

### Build

```sh
go mod download
go build -o scandisk.exe main.go
```