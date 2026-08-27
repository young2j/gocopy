# Repository Guidelines

## Project Structure & Module Organization

GoCopy is a single Go module (`github.com/young2j/gocopy`) with implementation files at the repository root. `copy.go` exposes `Copy`, `CopyWithOption`, and `Option`; the `copyMap.go`, `copySlice.go`, `copyStruct.go`, `copyStruct2Map.go`, and `helper.go` files contain the reflection-based copying paths. Root-level `*_test.go` files cover slice, map, struct, BSON, and option behavior. `example/` is a runnable usage sample, while `README.md` documents the public API and benchmarks.

## Build, Test, and Development Commands

- `go test ./...` runs the complete unit-test suite, including packages under `example/`.
- `go vet ./...` performs standard Go correctness checks.
- `gofmt -w $(find . -name '*.go' -type f)` formats all Go sources; run it on every changed file.
- `go run ./example` executes the sample program.
- `go test -v . -bench=. -benchmem -benchtime=1s -cpu=4` runs the repository benchmarks.

## Coding Style & Naming Conventions

Use standard `gofmt` output (tabs, grouped imports, one declaration per line) and idiomatic Go error/panic handling. Keep exported identifiers PascalCase and document new exported API. Match existing lowerCamelCase names for internal helpers such as `copyMap` and `indirectValue`. Keep conversion behavior in the appropriate copy* file rather than duplicating reflection logic.

## Testing Guidelines

Add or update a root-level `*_test.go` file beside the behavior under test. Existing tests use Go's `testing` package and names such as `Test_copySlice` and `Test_copyMapWithOption`; follow that area-based pattern. Exercise nil inputs, pointer/value combinations, and option flags when changing conversion code. Run `go test ./...` before submitting; include benchmark results when performance-sensitive code changes.

## Commit & Pull Request Guidelines

Use short, imperative commit subjects with the established `feat:` or `fix:` prefixes (for example, `fix: handle nil source value`). Pull requests should explain the behavior change, list verification commands and results, and link the related issue when one exists. Update `README.md` or `example/` when public API or usage changes; include focused screenshots only for documentation or example-output changes.
