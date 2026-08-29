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
