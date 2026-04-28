# AGENT.md

This file is a session handoff for future agent runs.
Use it to quickly understand the repo, prior decisions, and safe workflows.

## Project Identity

- Project: `goquarium`
- Language: Go
- Rendering library: `github.com/gdamore/tcell/v2`
- Module: `github.com/vak/goquarium`
- Entry point: `main.go`

## Current Repo Layout

Top-level Go files:

- `main.go` - process entrypoint
- `cli.go` - flag parsing and app dispatch
- `animation.go` - runtime loop, event handling, draw pipeline
- `entity.go` - entity model, movement, lifecycle, collisions
- `environment.go` - water, castle, seaweed
- `fish.go` - fish/bubbles/splat logic
- `special.go` - shark/ship/whale/monster/hook/ducks/dolphins/swan
- `setup.go` - startup scene composition
- `depth.go` - z-layer map
- `info.go`, `version.go` - user-facing text/version
- `entity_test.go` - focused unit tests

Non-code:

- `Makefile`
- `README.md`
- `LICENSE`

## Important Historical Decisions

1. **Top-level package layout**
   - Code intentionally moved from `cmd/...` + `internal/...` to root-level `package main`.
   - Build command is now `go build .`.

2. **Whale spout behavior parity**
   - Whale appears to dip slightly during spout frames.
   - This was confirmed to match the Python source behavior and is **intentional**.

3. **Surface sprite regression fixed**
   - Ship and other surface entities were repaired to use real multiline art (not literal `\n` text).
   - Regression test exists to catch this class of issue.

4. **Beginner-focused comments added**
   - All root Go files include explanatory comments aimed at non-expert readers.
   - Recent pass improved comments around callbacks, state machines, and frame/collision flow.

## Build/Test/Lint Workflow

Preferred quick validation:

```bash
gofmt -w *.go
go test ./...
go build .
```

Makefile workflow:

```bash
make
make test
make cover
```

## Known Makefile Note

- `install` and `uninstall` target paths are currently inconsistent:
  - install: `$(DESTDIR)/usr/bin/${PROG}`
  - uninstall: `$(DESTDIR)/bin/${PROG}`
- If you touch installation behavior, align these paths first.

## Behavioral Notes for Future Changes

- `Entity.CallbackArgs` is polymorphic:
  - often `[]float64{dx, dy, dz, frameStep}`
  - sometimes mode maps (for hook state machine)
- Collision pass is simple AABB O(n^2), intentionally readable.
- Render order depends on `Depth`; collision does not use depth.
- Many special entities chain via death callbacks into `RandomObject`.

## Tests Worth Keeping

- `TestSurfaceSpritesUseRealMultilineStrings` in `entity_test.go`
  - Protects against accidental literal `\n` sprite strings.
- Keep this test updated if adding new surface entities.

## If Starting a New Session

1. Run `go test ./...` to verify baseline.
2. If visuals look wrong, inspect sprite text and masks first.
3. For parity questions, compare against https://github.com/MKAbuMattar/asciiquarium-python/asciiquarium/entities/*.py.
4. Preserve beginner-friendly comments when refactoring.
