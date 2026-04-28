# goquarium

ASCII aquarium animation for your terminal, written in Go with `tcell`.

This project is a Go rewrite inspired by the classic asciiquarium and the Python port in `asciiquarium-python/`.

## Features

- Animated fish, sharks, whales, ships, monsters, and other sea creatures
- Environment elements like water lines, seaweed, and castle decoration
- Color rendering with per-character color masks
- Interactive keyboard controls
- Works in standard terminals on Linux/macOS/Windows (with terminal support)

## Documentation

- [Visual Entity Catalog](Entities.md) - Full reference for all rendered entities, grouped by class with implementation details.

## Requirements

- Go 1.22+
- A terminal with color support
- Recommended terminal size: at least `40x15`

## Build

Using Makefile:

```bash
make
```

Using Go directly:

```bash
go build
```

## Run

```bash
./goquarium
```

Or run without building:

```bash
go run .
```

## CLI Options

```bash
./goquarium --help
./goquarium --version
./goquarium --info
./goquarium --classic
```

## Controls

- `q` - Quit
- `p` - Pause/unpause
- `r` - Reset and respawn entities
- `i` - Toggle info overlay
- `Esc` - Close info overlay

## Development

Run tests:

```bash
make test
```

## Install / Uninstall

```bash
make install
make uninstall
```

By default, the Makefile installs to `$(HOME)/.local/usr/bin/goquarium`.

## Credits

- [Original Perl version `asciiquarium`](http://robobunny.com/projects/asciiquarium): Kirk Baucom
- [Python port `asciiquarium-python`](https://github.com/MKAbuMattar/asciiquarium-python): Mohammad Abu Mattar

## License

GNU GPL v3. See [LICENSE](LICENSE).
