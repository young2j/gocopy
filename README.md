# GoCopy

GoCopy is a reflection-based Go library for copying compatible slices, maps,
and structs. It handles nested pointers and collections, field tags, custom
converters, and MongoDB's `bson.M` (`map[string]interface{}`). Field metadata
and common conversion paths are cached for low allocation overhead.

## Installation

Go 1.16 or newer is supported.

```sh
go get github.com/young2j/gocopy@latest
```

## Quick start

The destination must be writable, so pass a pointer. `CopyE` is useful at input
boundaries because it returns conversion errors instead of panicking:

```go
package main

import (
    "fmt"
    "log"

    "github.com/young2j/gocopy"
)

type User struct {
    UserName string
    Age      int
    Address  struct {
        City string
    }
}

func main() {
    input := map[string]interface{}{
        "userName": "alice",
        "age":      int64(30),
        "address":  map[string]interface{}{"city": "Shanghai"},
    }

    var user User
    if err := gocopy.CopyE(&user, input); err != nil {
        log.Fatal(err)
    }

    output := map[string]interface{}{}
    gocopy.CopyWithOption(&output, user, &gocopy.Option{ToCase: "Snake"})
    fmt.Printf("%#v\n", output)
}
```

The example prints keys such as `user_name` and recursively converts the
nested `address` value to a map.

## Supported conversions

| Source | Destination | Notes |
| --- | --- | --- |
| `[]T` | `[]U` | Assignable, convertible, or nested element conversion |
| `map[K]V` | `map[K2]V2` | Nested maps, slices, structs, and pointers supported |
| struct | struct | Matches exported fields by name |
| struct | map / `bson.M` | Tags, case conversion, nested map output |
| map | struct | String keys or interface keys containing strings |

For every conversion, `to` should be a pointer (for example, `&dst`). Nil
source values are ignored at the top level; nil collection fields retain the
library's documented nil/empty representation.

## API

| Function | Purpose | Failure behavior |
| --- | --- | --- |
| `Copy(to, from)` | Copy with default options | Panics on unsupported or invalid conversion |
| `CopyWithOption(to, from, opt)` | Copy with an `Option` | Panics on unsupported or invalid conversion |
| `CopyE(to, from)` | Copy with default options | Returns conversion errors |
| `CopyWithOptionE(to, from, opt)` | Copy with an `Option` | Returns conversion errors |

`Copy` and `CopyWithOption` are convenient when conversion failures are
programmer errors. Use the `E` variants when source data is external or may
have an incompatible type. Passing `nil` for `opt` is equivalent to
`&Option{}`.

All four functions accept source values or pointers. The destination must be a
writable pointer; for example, use `&dst`, not `dst`. A conversion changes the
destination in place and does not return a new value.

Typical calls are identical for every supported shape:

```go
var structDst User
gocopy.Copy(&structDst, structSrc) // struct → struct

var sliceDst []User
gocopy.Copy(&sliceDst, sliceSrc) // []T → []U

var mapDst map[string]interface{}
gocopy.Copy(&mapDst, structSrc) // struct → map / bson.M
```

## Options

| Field | Applies to | Default | Behavior |
| --- | --- | --- | --- |
| `Append` | slices, maps, nested fields | `false` | Merge into existing destination values instead of replacing them |
| `NameFromTo` | struct conversions | empty | Rename source fields: `map[string]string{"Source": "Target"}` |
| `ToCase` | struct↔map | `LowerCamel` | Key casing: `Camel`, `Snake`, `ScreamingSnake`, `Kebab`, or `ScreamingKebab` are also supported |
| `IgnoreZero` | struct→map | `false` | Omit zero-valued fields |
| `IgnoreFields` | struct conversions | empty | Skip listed source field names |
| `IgnoreLevel` | embedded fields | `0` | Apply `IgnoreFields` through this embedding depth |
| `Converters` | struct conversions | empty | Per-field function: `map[string]func(interface{}) interface{}` |
| `TimeToString` | `time.Time` fields | empty | Format with optional `loc` and `layout` settings |
| `StringToTime` | string fields | empty | Parse with optional `loc` and `layout` settings |

Example:

```go
opt := &gocopy.Option{
    Append:       true,
    NameFromTo:   map[string]string{"CreatedAt": "Created"},
    IgnoreZero:   true,
    IgnoreFields: []string{"InternalID"},
    TimeToString: map[string]map[string]string{
        "CreatedAt": {"loc": "UTC", "layout": "2006-01-02"},
    },
}
gocopy.CopyWithOption(&destination, source, opt)
```

For time options, omitting `loc` uses `Asia/Shanghai`; omitting `layout` uses
`2006-01-02 15:04:05`. In append mode, nested slices and maps are merged while
scalar fields are replaced by the source value.

### Field tags

Tags are checked in this order: `gocopy`, `json`, then `bson`. A tag name
controls the map key or lookup name; `omitempty` skips empty values and `-`
skips the field entirely:

```go
type Record struct {
    UserID string `gocopy:"user_id,omitempty"`
    Secret string `json:"-"`
}
```

## Examples, tests, and benchmarks

The runnable tour covers slices, maps, struct conversions, BSON, append mode,
time conversion, and custom converters:

```sh
go run ./example
go test ./...
go test -race ./...
go vet ./...
```

Run benchmarks with:

```sh
go test -run '^$' -bench . -benchmem -benchtime=1s -cpu=4
```

Results from Apple M4 Pro (darwin/arm64, Go toolchain on 2026-08-30):

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `BenchmarkCopy-4` | 745.9 | 592 | 7 |
| `BenchmarkCopier-4` (comparison) | 5432 | 13728 | 49 |
| `BenchmarkCopyPlainStruct-4` | 39.47 | 0 | 0 |
| `BenchmarkCopyPlainStruct2Map-4` | 174.0 | 336 | 2 |
| `BenchmarkCopyMap2Struct-4` | 228.6 | 56 | 2 |
| `BenchmarkCopyMap2StructInterfaceKeys-4` | 516.9 | 0 | 0 |
| `BenchmarkCopyStruct2Map-4` | 847.9 | 1384 | 18 |
| `BenchmarkCopyStruct2MapNested-4` | 1645 | 2920 | 41 |
| `BenchmarkCopyStruct2MapNestedValues-4` | 656.5 | 1264 | 20 |
| `BenchmarkCopySliceNested-4` | 347.6 | 152 | 3 |
| `BenchmarkCopyMapNested-4` | 296.3 | 336 | 4 |

Benchmark numbers vary by CPU, Go version, and workload; rerun the command for
your target environment.

## License

MIT. See [`LICENSE`](LICENSE).
