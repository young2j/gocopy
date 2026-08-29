# GoCopy

GoCopy is a small, reflection-based Go library for copying compatible data
without hand-written field-by-field conversion code. It supports nested values,
pointer/value combinations, and common map representations such as
`map[string]interface{}` and MongoDB's `bson.M`.

## Features

- slice ↔ slice and map ↔ map conversion, including nested collections;
- struct ↔ struct conversion by exported field name;
- struct ↔ map conversion and map → struct decoding;
- field tags: `gocopy`, then `json` or `bson` as fallbacks;
- append, rename, case conversion, zero-value omission, and custom converters;
- time formatting/parsing and error-returning `CopyE` APIs;
- cached field metadata and direct fast paths for common conversions.

## Installation

Go 1.16 or newer is supported.

```sh
go get github.com/young2j/gocopy@latest
```

## Quick start

The destination must be writable, normally a pointer:

```go
type User struct {
    UserName string `gocopy:"user_name"`
    Age      int
}

var user User
if err := gocopy.CopyE(&user, map[string]interface{}{
    "user_name": "alice",
    "age":       int64(30),
}); err != nil {
    log.Fatal(err)
}

result := map[string]interface{}{}
gocopy.Copy(&result, user)
```

## API and options

`Copy(to, from)` uses defaults. `CopyWithOption` accepts an `Option`; the
`CopyE` and `CopyWithOptionE` variants return conversion errors instead of
panicking. `Append` merges existing slices/maps. `NameFromTo` renames fields
(`source: target`). `ToCase` controls struct-to-map keys and map-to-struct
lookup (`LowerCamel`, `Snake`, `Camel`, `Kebab`, and screaming variants). `IgnoreZero` and
`IgnoreFields` omit data, while `Converters` handles application-specific
types. `TimeToString` and `StringToTime` accept `loc` and `layout` settings.

## Examples and verification

See [`example/`](example) and [`example/README.md`](example/README.md):

```sh
go run ./example
go test ./...
go test -race ./...
go vet ./...
```

Run benchmarks with:

```sh
go test -bench=. -benchmem -benchtime=1s -cpu=4
```

## License

MIT. See [`LICENSE`](LICENSE).
