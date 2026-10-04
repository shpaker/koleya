# Rules for AI agents

koleya is grid-aligned movement for Go games that look old but feel new. Module
`github.com/shpaker/koleya`, one package, no dependencies.

## Architecture

```
koleya.go    package docs, Instant, Events
vec.go       Axis, Vec2
dir.go       Dir: eight directions and None
lattice.go   Lattice: the nodes vehicles rest on
profile.go   Profile, Steering, Classic, Surface
axis.go      the one-dimensional solver: drive, pick a node, dock on it
mover.go     Mover: the state a game keeps in its entity
step.go      Step: Four and Eight steering on top of the solver
```

Architecture and the right abstractions come first.

**Scope: kinematics on a grid.** koleya moves a point: speeds, braking, nodes.
Collisions, input devices, drawing, sound and the shape of a vehicle belong to the game.
The game rolls a position back and calls `Halt` or `HaltAxis`; koleya never asks about
the world.

- Pure math: no engine, no IO, no platform code (depguard checks). Works the same on
  every platform Go builds for.
- `Step` reads nothing but its arguments: no global state, no wall clock, no randomness.
  Time comes from the game as `dt`.
- `Mover` is a plain value the game keeps in its entity; its fields are private except
  `Pos`, which the game may roll back.
- `Profile` is data: vehicle types are presets the game defines. A new feel is a new
  field with a zero value that keeps the old behavior.
- All the braking math lives in the one-dimensional solver (`axis.go`); steering only
  decides what each axis wants.
- The promise: at rest, unless halted, both coordinates lie on lattice nodes; the speed
  never exceeds `MaxSpeed`. `FuzzRestsOnGrid` checks it; a change must keep it green.

## Principles

- One type — one responsibility.
- No deprecated APIs (e.g. `rand.Seed`).

## Naming

- Constructors `New*`; getters without `Get` (`Facing`, `Vel`), setters `Set*`.
- No abbreviations.
- Short clear names.
- Fields as private as possible, grouped by kind.

## Style

- Standard Go conventions; short variable declarations where they fit.
- Comments and docs in English, plain words.
- Document every exported identifier.
- Tests as simple and readable as possible. `example_test.go` mirrors the README's
  quick start, so the README never drifts from the code.

## Required after changes

- `just fmt`, `just lint`, `just lint-wasm`, `just test`, `just fuzz`; the linter stays
  clean for both targets (desktop and js/wasm).
- Find and update every use of a changed identifier; update the tests.
- Update README when the API changes: what a game developer sees — short and to the
  point.

## Git

- Any writing git operation (commit, push, merge, rebase, reset, tag and the like) only with
  the user's explicit confirmation.
- Commits are made as the user only.
- No mention of AI or neural networks (`Co-Authored-By`, `Generated with ...`, `claude/`,
  `ai/` and the like) in commits, their metadata, branch names and PRs.
- PR and release descriptions as short as possible; the only formatting is a bulleted list,
  no headings or emphasis.
- Conventional Commits: `feat`, `fix`, `perf`, `refactor` raise the minor version,
  `type!:` or a `BREAKING CHANGE:` footer the major (while at v0 — the minor); `docs`,
  `chore`, `ci`, `test`, `style`, `build` make no release.
- Release tags (`vMAJOR.MINOR.0`, as Go modules require) are set by CI
  (`.github/workflows/version.yml`), run by hand on main with a green CI: changes pile up,
  a release goes out when ready; it tags and publishes the release with its changelog.
  Never tag by hand.
