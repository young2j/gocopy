# GoCopy examples

Run the complete reflection/conversion tour from the repository root:

```sh
go run ./example
```

The program groups examples by conversion type:

- slices, including pointer elements and `Option{Append: true}`;
- maps with nested slices, maps, and structs;
- struct-to-struct matching and custom converters;
- struct-to-map output for `map[string]interface{}` and `bson.M`;
- map-to-struct input, tagged fields, and numeric conversion;
- time formatting/parsing through `TimeToString` and `StringToTime`.

## Converters

`Option.Converters` is a map keyed by the **source struct field name**. Each
function receives that field's original value as `interface{}` and returns the
value that should be assigned to the destination field (or map entry):

```go
type Source struct { Cents int64 }
type Destination struct { Amount float64 }

var dst Destination
err := gocopy.CopyWithOptionE(&dst, Source{Cents: 1299}, &gocopy.Option{
    NameFromTo: map[string]string{"Cents": "Amount"},
    Converters: map[string]func(interface{}) interface{}{
        "Cents": func(v interface{}) interface{} {
            cents := v.(int64)       // assert the source field type
            return float64(cents) / 100 // return the destination field type
        },
    },
})
```

When source and destination field names differ, add a `NameFromTo` entry as in
the example. The callback is invoked once per matching field. For pointer fields, `v` is
the pointer value; dereference it inside the callback when needed. The return
value must be assignable to the destination type. Use `CopyWithOptionE` when a
bad type assertion or conversion should be reported as an error.

A minimal map/struct round trip looks like this:

```go
type User struct {
    UserName string `gocopy:"user_name"`
    Age      int
}

var user User
err := gocopy.CopyE(&user, map[string]interface{}{
    "user_name": "alice",
    "age":       int64(30),
})
if err != nil {
    log.Fatal(err)
}

result := map[string]interface{}{}
gocopy.Copy(&result, user)
```

The destination must be writable (usually passed by pointer). See the root
[`README.md`](../README.md) for all options and supported tag formats.
