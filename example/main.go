/*
 * File: test.go
 * Created Date: 2022-01-29 03:05:53
 * Author: ysj
 * Description: compare with copier
 */

package main

import (
	"fmt"

	// "github.com/jinzhu/copier"
	"github.com/young2j/gocopy/example/model"
	"github.com/young2j/gocopy/example/types"

	// "go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/globalsign/mgo/bson"
	"github.com/young2j/gocopy"
)

func main() {
	// slice to slice
	// 1. simple slice
	s1 := []int{3, 4, 5}
	s2 := make([]int, 0)
	gocopy.Copy(&s2, &s1)
	// copier.Copy(&s2, &s1)
	fmt.Println("==============================")
	fmt.Printf("s2: %v\n", s2)

	// 2. map slice
	ms1 := []map[string]int{{"key1": 1, "key2": 2}}
	ms2 := make([]map[string]int, 0)
	gocopy.Copy(&ms2, &ms1)
	// copier.Copy(&ms2, &ms1)
	fmt.Println("==============================")
	fmt.Printf("ms2: %v\n", ms2)

	// 3. struct slice
	ss1 := []model.Perm{{Action: "GET", Label: "rest-get-method"}}
	ss2 := make([]types.Perm, 0)
	gocopy.Copy(&ss2, &ss1)
	// copier.Copy(&ss2, &ss1)
	fmt.Println("==============================")
	fmt.Printf("ss2: %v\n", ss2)

	// map to map
	// 1. simple map
	m1 := map[string]int{"key1": 1, "key2": 2}
	m2 := make(map[string]int)
	gocopy.Copy(&m2, &m1)
	// copier.Copy(&m2, &m1)
	fmt.Println("==============================")
	fmt.Printf("m2: %v\n", m2)

	// 2. slice map
	sm1 := map[int][]string{1: {"a", "b", "c"}}
	sm2 := make(map[int][]string)
	gocopy.Copy(&sm2, &sm1)
	// copier.Copy(&sm2, &sm1)
	fmt.Println("==============================")
	fmt.Printf("sm2: %v\n", sm2)

	// 3. struct map
	stm1 := map[string]model.Perm{"perm1": {Action: "POST", Label: "rest-post-method"}}
	stm2 := make(map[string]types.Perm)
	gocopy.Copy(&stm2, &stm1)
	// copier.Copy(&stm2, &stm1)
	fmt.Println("==============================")
	fmt.Printf("stm2: %v\n", stm2)

	// 4. struct to struct
	roll := 100
	st1 := model.AccessRolePerms{
		Role: "角色",
		Roll: &roll,
		EmbedFields: model.EmbedFields{
			EmbedF1: "embedF1",
		},
		Actions: []string{"xxx", "yyy"},
		Perms:   []*model.Perm{{Action: "GET", Label: "fuck"}},
		PermMap: map[string]*model.Perm{"perm": {Action: "PUT", Label: "修改"}},
	}
	st2 := types.AccessRolePerms{}
	gocopy.Copy(&st2, &st1)
	// copier.Copy(&st2, &st1)
	fmt.Println("==============================")
	fmt.Printf("st2: %#v\n", st2)
	fmt.Printf("st2.Role: %v\n", *st2.Role)
	fmt.Printf("st2.Roll: %v\n", *st2.Roll)

	for _, v := range st2.Perms {
		fmt.Printf("Perms: %#v\n", v)
	}
	for k, v := range st2.PermMap {
		fmt.Printf("PermMap k:%v v:%#v\n", k, v)
	}

	// append mode
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
		Perms:   []*model.Perm{{Action: "PUT", Label: "put-label"}},
		PermMap: map[string]*model.Perm{"delete": {Action: "DELETE", Label: "delete-label"}},
	}
	ast2 := types.AccessRolePerms{
		Actions: []string{"GET", "POST"},
		Perms:   []*types.Perm{{Action: "GET", Label: "get-label"}},
		PermMap: map[string]*types.Perm{"get": {Action: "GET", Label: "get-label"}},
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

	// from field to another field
	// objectId to string and vice versa
	ost1 := model.AccessRolePerms{
		From: "fromto",
		Id:   bson.NewObjectId(),
		// Id:    primitive.NewObjectID(),
		IdHex: "61f04828eb37b662c8f3b085",
	}
	ost2 := types.AccessRolePerms{}
	opt := gocopy.Option{
		NameFromTo:       map[string]string{"From": "To"},
		ObjectIdToString: map[string]string{"Id": "mgo"},
		StringToObjectId: map[string]string{"IdHex": "mgo"},
		// ObjectIdToString: map[string]string{"Id": "official"},
		// StringToObjectId: map[string]string{"IdHex": "official"},
	}
	gocopy.CopyWithOption(&ost2, &ost1, &opt)
	fmt.Println("==============================")
	fmt.Printf("ost2.To: %v\n", ost2.To)
	fmt.Printf("ost2.Id: %v\n", ost2.Id)
	fmt.Printf("ost2.IdHex: %v\n", ost2.IdHex)
}
