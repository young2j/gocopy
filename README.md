# Features
* copy slice to slice by type
* copy map to map by type
* copy struct to struct by field name
* copy struct to map(support field name case and ignore zero value)
* support append mode
* support copy bson.ObjectId to string and vice versa
* support copy time.Time to string and vice versa

# Installation

```shell
go get -u github.com/young2j/gocopy@latest
```

# Usage

## API

* `Copy(to, from interface{})`
* `CopyWithOption(to,from interface{},opt *Option)`

> Note: The arg `to` must be a reference value(usually `&to`), or copy maybe fail. 

see more  at [`/example`](https://github.com/young2j/gocopy/tree/master/example)

define source types as follows:

```go
// model.go
package model
import (
  "github.com/globalsign/mgo/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Perm struct {
  Action string
  Label  string
}

type EmbedFields struct {
  EmbedF1 string
  embedF2 string
}

type AccessRolePerms struct {
  Id1     bson.ObjectId
  Id2     primitive.ObjectID
  Id1Hex  string
  Id2Hex  string
  Role    string
  Roll    *int
  Actions []string
  Perms   []*Perm
  PermMap map[string]*Perm
  From    string
  EmbedFields
}
```

define destination types as follows:

```go
// types.go
package types

import (
  "github.com/globalsign/mgo/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Perm struct {
  Action string
  Label  string
}

type AccessRolePerms struct {
  Id1     string
  Id2     string
  Id1Hex  bson.ObjectId
  Id2Hex  primitive.ObjectID
  Role    *string
  Roll    *int
  Actions []string
  Perms   []*Perm
  PermMap map[string]*Perm
  To      []byte
  EmbedF1 string
}
```

## Copy

### Copy slice

```go
// 1. simple slice
s1 := []int{3, 4, 5}
s2 := make([]int, 0)
gocopy.Copy(&s2, &s1)
fmt.Println("==============================")
fmt.Printf("s2: %v\n", s2)

// 2. map slice
ms1 := []map[string]int{{"key1": 1, "key2": 2}}
ms2 := make([]map[string]int, 0)
gocopy.Copy(&ms2, &ms1)
fmt.Println("==============================")
fmt.Printf("ms2: %v\n", ms2)

// 3. struct slice
ss1 := []model.Perm{{Action: "GET", Label: "rest-get-method"}}
ss2 := make([]types.Perm, 0)
gocopy.Copy(&ss2, &ss1)
fmt.Println("==============================")
fmt.Printf("ss2: %v\n", ss2)
```

```shell
==============================
s2: [3 4 5]
==============================
ms2: [map[key1:1 key2:2]]
==============================
ss2: [{GET rest-get-method}]
```

### Copy map

```go
// map to map
// 1. simple map
m1 := map[string]int{"key1": 1, "key2": 2}
m2 := make(map[string]int)
gocopy.Copy(&m2, &m1)
fmt.Println("==============================")
fmt.Printf("m2: %v\n", m2)

// 2. slice map
sm1 := map[int][]string{1: {"a", "b", "c"}}
sm2 := make(map[int][]string)
gocopy.Copy(&sm2, &sm1)
fmt.Println("==============================")
fmt.Printf("sm2: %v\n", sm2)

// 3. struct map
stm1 := map[string]model.Perm{"perm1": {Action: "POST", Label: "rest-post-method"}}
stm2 := make(map[string]types.Perm)
gocopy.Copy(&stm2, &stm1)
fmt.Println("==============================")
fmt.Printf("stm2: %v\n", stm2)
```

```shell
==============================
m2: map[key1:1 key2:2]
==============================
sm2: map[1:[a b c]]
==============================
stm2: map[perm1:{POST rest-post-method}]
```

### Copy struct

#### struct to struct

```go
// 4. struct to struct
roll := 100
st1 := model.AccessRolePerms{
  Role: "角色",
  Roll: &roll,
  EmbedFields: model.EmbedFields{
    EmbedF1: "embedF1",
  },
  Actions: []string{"GET", "POST"},
  Perms:   []*model.Perm{{Action: "GET", Label: "rest-get-method"}},
  PermMap: map[string]*model.Perm{"perm": {Action: "PUT", Label: "rest-put-method"}},
}
st2 := types.AccessRolePerms{}
gocopy.Copy(&st2, &st1)
fmt.Println("==============================")
fmt.Printf("st2.Role: %v\n", *st2.Role)
fmt.Printf("st2.Roll: %v\n", *st2.Roll)
fmt.Printf("st2.Actions: %v\n", st2.Actions)

for _, v := range st2.Perms {
  fmt.Printf("Perms: %#v\n", v)
}
for k, v := range st2.PermMap {
  fmt.Printf("PermMap k:%v v:%#v\n", k, v)
}
```

```shell
==============================
st2.Role: 角色
st2.Roll: 100
st2.Actions: [GET POST]
Perms: &types.Perm{Action:"GET", Label:"rest-get-method"}
PermMap k:perm v:&types.Perm{Action:"PUT", Label:"rest-put-method"}
```

#### struct to map/bson.M

```go
fromst := model.AccessRolePerms{
  CreatedAt: time.Now(),
  UpdatedAt: "2022/02/16",
  Id1:       bson.NewObjectId(),
  Id2:       primitive.NewObjectID(),
  Id1Hex:    bson.NewObjectId().Hex(),
  Id2Hex:    primitive.NewObjectID().Hex(),
  Role:      "copystruct2map",
  Roll:      &roll,
  From:      "From",
  Actions:   []string{"PUT", "DELETE"},
  Perms:     []*model.Perm{{Action: "PUT", Label: "rest-put-method"}},
  PermMap:   map[string]*model.Perm{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
  Child: &model.AccessRolePerms{
    Id1Hex: bson.NewObjectId().Hex(),
    Role:   "embedstruct",
  },
}
toBM := bson.M{}
gocopy.CopyWithOption(&toBM, fromst, &gocopy.Option{
  IgnoreZero:  true,
})
fmt.Println("==============================")
fmt.Printf("toBM[\"createdAt\"]: %v\n", toBM["createdAt"])
fmt.Printf("toBM[\"updatedAt\"]: %v\n", toBM["updatedAt"])
fmt.Printf("toBM[\"id1\"]: %v\n", toBM["id1"])
fmt.Printf("toBM[\"id2\"]: %v\n", toBM["id2"])
fmt.Printf("toBM[\"id1Hex\"]: %v\n", toBM["id1Hex"])
fmt.Printf("toBM[\"id2Hex\"]: %v\n", toBM["id2Hex"])
fmt.Printf("toBM[\"role\"]: %v\n", toBM["role"])
fmt.Printf("toBM[\"roll\"]: %v\n", *toBM["roll"].(*int))
fmt.Printf("toBM[\"actions\"]: %v\n", toBM["actions"])
fmt.Printf("toBM[\"perms\"]: %#v\n", toBM["perms"])
fmt.Printf("toBM[\"permMap\"]: %#v\n", toBM["permMap"])
fmt.Printf("toBM[\"child\"]: %#v\n", toBM["child"])
```

```shell
==============================
toBM["createdAt"]: 2022-02-17 18:31:38.317398 +0800 CST m=+0.002843322
toBM["updatedAt"]: 2022/02/16
toBM["id1"]: ObjectIdHex("620e240aeb37b6263f59fd51")
toBM["id2"]: ObjectID("620e240a9ca780f5f4bb7c0a")
toBM["id1Hex"]: 620e240aeb37b6263f59fd52
toBM["id2Hex"]: 620e240a9ca780f5f4bb7c0b
toBM["role"]: copystruct2map
toBM["roll"]: 100
toBM["actions"]: [PUT DELETE]
toBM["perms"]: []*bson.M{(*bson.M)(0xc00009a290)}
toBM["permMap"]: bson.M{"delete":(*bson.M)(0xc00009a2b0)}
toBM["child"]: &bson.M{"id1Hex":"620e240aeb37b6263f59fd53", "role":"embedstruct"}
```

## CopyWithOption

### Append mode

```go
opts := gocopy.Option{Append: true}
// 1. append slice
as1 := []int{3, 4, 5}
as2 := []int{1, 2}
gocopy.CopyWithOption(&as2, &as1, &opts)
fmt.Println("==============================")
fmt.Printf("as2: %v\n", as2)

// 2. append map
am1 := map[string]int{"key1": 1, "key2": 2}
am2 := map[string]int{"key0": 0, "key2": 3}
gocopy.CopyWithOption(&am2, &am1, &opts)
fmt.Println("==============================")
fmt.Printf("am2: %v\n", am2)

ams1 := map[string][]int{"key1": {1}, "key2": {2}}
ams2 := map[string][]int{"key0": {0}, "key2": {3}}
gocopy.CopyWithOption(&ams2, &ams1, &opts)
fmt.Println("==============================")
fmt.Printf("ams2: %v\n", ams2)

// 3. append struct map/slice field
ast1 := model.AccessRolePerms{
  Actions: []string{"PUT", "DELETE"},
  Perms:   []*model.Perm{{Action: "PUT", Label: "rest-put-method"}},
  PermMap: map[string]*model.Perm{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
}
ast2 := types.AccessRolePerms{
  Actions: []string{"GET", "POST"},
  Perms:   []*types.Perm{{Action: "GET", Label: "rest-get-method"}},
  PermMap: map[string]*types.Perm{"get": {Action: "GET", Label: "rest-get-method"}},
}
gocopy.CopyWithOption(&ast2, &ast1, &opts)
fmt.Println("==============================")
fmt.Printf("ast2.Actions: %v\n", ast2.Actions)
for i, perm := range ast2.Perms {
  fmt.Printf("ast2.Perms[%v]: %#v\n", i, perm)
}
for i, pm := range ast2.PermMap {
  fmt.Printf("ast2.PermMap[%v]: %#v\n", i, pm)
}
```

```shell
==============================
as2: [1 2 3 4 5]
==============================
am2: map[key0:0 key1:1 key2:2]
==============================
ams2: map[key0:[0] key1:[1] key2:[3 2]]
==============================
ast2.Actions: [GET POST PUT DELETE]
ast2.Perms[0]: &types.Perm{Action:"GET", Label:"rest-get-method"}
ast2.Perms[1]: &types.Perm{Action:"PUT", Label:"rest-put-method"}
ast2.PermMap[delete]: &types.Perm{Action:"DELETE", Label:"rest-delete-method"}
ast2.PermMap[get]: &types.Perm{Action:"GET", Label:"rest-get-method"}
```

### Specify field name

```go
// from field to another field
ost1 := model.AccessRolePerms{
  From: "fromto",
}
ost2 := types.AccessRolePerms{}
opt := gocopy.Option{
  NameFromTo:       map[string]string{"From": "To"},
}
gocopy.CopyWithOption(&ost2, &ost1, &opt)
fmt.Println("==============================")
fmt.Printf("ost2.To: %v\n", ost2.To)
```

```shell
==============================
ost2.To: fromto
```

### ObjectId and String

```go
// objectId to string and vice versa
from := model.AccessRolePerms{
  Id1:    bson.NewObjectId(),
  Id2:    primitive.NewObjectID(),
  Id1Hex: "61f04828eb37b662c8f3b085",
  Id2Hex: "61f04828eb37b662c8f3b085",
}
to := types.AccessRolePerms{
  Actions: []string{"GET", "POST"},
}
option := &gocopy.Option{
  ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
  StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
  Append:           true,
}
gocopy.CopyWithOption(&to, from, option)
fmt.Println("==============================")
fmt.Printf("from.Id1: %v to.Id1: %v \n", from.Id1, to.Id1)
fmt.Printf("from.Id2: %v to.Id2: %v \n", from.Id2, to.Id2)
fmt.Printf("from.Id1Hex: %v to.Id1Hex: %v \n", from.Id1Hex, to.Id1Hex)
fmt.Printf("from.Id2Hex: %v to.Id2Hex: %v \n", from.Id2Hex, to.Id2Hex)
```

```shell
==============================
from.Id1: ObjectIdHex("61f6cdf318ef1d4366bca973") to.Id1:61f6cdf318ef1d4366bca973
from.Id2: ObjectID("61f6cdf3cc541c1bc35a41fc") to.Id2:61f6cdf3cc541c1bc35a41fc
from.Id1Hex: 61f04828eb37b662c8f3b085 to.Id1Hex:ObjectIdHex("61f04828eb37b662c8f3b085")
from.Id2Hex: 61f04828eb37b662c8f3b085 to.Id2Hex:ObjectID("61f04828eb37b662c8f3b085")
```

### Time.Time and String

```go
from := model.AccessRolePerms{
  CreatedAt: time.Now(),
  UpdatedAt: "2022/02/11 15:04:05",
}
to := types.AccessRolePerms{}
option := gocopy.Option{
  // TimeToString: map[string]map[string]string{"CreatedAt": nil},
  // StringToTime: map[string]map[string]string{"UpdatedAt": nil},
  TimeToString: map[string]map[string]string{"CreatedAt": {"layout": "2006-01-02", "loc": "America/New_York"}},
  StringToTime: map[string]map[string]string{"UpdatedAt": {"layout": "2006/01/02 15:04:05"}},
}
gocopy.CopyWithOption(&to, from, &option)
fmt.Println("==============================")
fmt.Printf("time.Time to string-> to.CreatedAt: %v\n", to.CreatedAt)
fmt.Printf("string to time.Time-> to.UpdatedAt: %v\n", to.UpdatedAt)
```

```shell
==============================
time.Time to string-> to.CreatedAt: 2022-02-16
string to time.Time-> to.UpdatedAt: 2022-02-11 15:04:05 +0800 CST
```

# Related repository

* [`copier`](https://github.com/jinzhu/copier)

# Benchmark

```shell
go test -v . -bench=.  -benchmem -benchtime=1s -cpu=4
```

```shell
goos: darwin
goarch: amd64
pkg: github.com/young2j/gocopy
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkCopy
BenchmarkCopy-4     	  122139	      8884 ns/op	    5592 B/op	      81 allocs/op
BenchmarkCopier
BenchmarkCopier-4   	   62940	     18695 ns/op	   14640 B/op	     166 allocs/op
PASS
ok  	github.com/young2j/gocopy	4.999s
```



