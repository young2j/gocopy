/*
 * File: copy.go
 * Created Date: 2022-01-26 05:36:00
 * Author: ysj
 * Description: copy slice and map by type; copy struct by field name
 */
package gocopy

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Perm1 struct {
	Action string
	Label  string
}

type Perm2 struct {
	Action string
	Label  string
}

type EmbedFields struct {
	EmbedF1 string
	embedF2 string
}

type AccessRolePerms1 struct {
	Id1     bson.ObjectId
	Id2     primitive.ObjectID
	Id1Hex  string
	Id2Hex  string
	Role    string
	Roll    *int
	Actions []string
	Perms   []*Perm1
	PermMap map[string]*Perm1
	From    string
	EmbedFields
}

type AccessRolePerms2 struct {
	Id1     string
	Id2     string
	Id1Hex  bson.ObjectId
	Id2Hex  primitive.ObjectID
	Role    *string
	Roll    *int
	Actions []string
	Perms   []*Perm2
	PermMap map[string]*Perm2
	To      []byte
	EmbedF1 string
}

func TestCopy(t *testing.T) {
	type args struct {
		from interface{}
		to   interface{}
	}
	roll := 100
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copyslice-simple",
			args: args{
				from: []int{3, 4, 5},
				to:   make([]int, 0),
			},
		},
		{
			name: "copyslice-map",
			args: args{
				from: []map[string]int{{"key1": 1, "key2": 2}},
				to:   make([]map[string]int, 0),
			},
		},
		{
			name: "copyslice-struct",
			args: args{
				from: []Perm1{{Action: "GET", Label: "rest-get-method"}},
				to:   make([]Perm2, 0),
			},
		},
		{
			name: "copymap-simple",
			args: args{
				from: map[string]int{"key1": 1, "key2": 2},
				to:   make(map[string]int),
			},
		},
		{
			name: "copymap-slice",
			args: args{
				from: map[int][]string{1: {"a", "b", "c"}},
				to:   make(map[int][]string),
			},
		},
		{
			name: "copymap-struct",
			args: args{
				from: map[string]Perm1{"perm1": {Action: "POST", Label: "rest-post-method"}},
				to:   make(map[string]Perm2),
			},
		},
		{
			name: "copystruct",
			args: args{
				from: AccessRolePerms1{
					Role: "角色",
					Roll: &roll,
					EmbedFields: EmbedFields{
						EmbedF1: "embedF1",
					},
					Actions: []string{"xxx", "yyy"},
					Perms:   []*Perm1{{Action: "GET", Label: "rest-get-method"}},
					PermMap: map[string]*Perm1{"perm": {Action: "PUT", Label: "rest-put-method"}},
				},
				to: AccessRolePerms2{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copyslice-simple":
				from, ok := tt.args.from.([]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for i := 0; i < len(to); i++ {
					t.Logf("from[%d] == to[%d] = %v\n", i, i, from[i] == to[i])
					if from[i] != to[i] {
						t.Fail()
					}
				}
			case "copyslice-map":
				from, ok := tt.args.from.([]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for i := 0; i < len(to); i++ {
					elem := to[i]
					for k, v := range elem {
						t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, from[i][k] == v)
						if from[i][k] != v {
							t.Fail()
						}
					}
				}
			case "copyslice-struct":
				from, ok := tt.args.from.([]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for i := 0; i < len(to); i++ {
					elem := to[i]
					t.Logf("to.Action==from.Action: %v\n", elem.Action == from[i].Action)
					t.Logf("to.Label==from.Label: %v\n", elem.Label == from[i].Label)
					if elem.Action != from[i].Action {
						t.Fail()
					}
					if elem.Label != from[i].Label {
						t.Fail()
					}
				}
			case "copymap-simple":
				from, ok := tt.args.from.(map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for k, v := range to {
					t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, from[k] == v)
					if from[k] != v {
						t.Fail()
					}
				}
			case "copymap-slice":
				from, ok := tt.args.from.(map[int][]string)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[int][]string)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for k, v := range to {
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, v)
					for i := 0; i < len(v); i++ {
						if v[i] != from[k][i] {
							t.Fail()
						}
					}
				}
			case "copymap-struct":
				from, ok := tt.args.from.(map[string]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				for k, v := range to {
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, v)
					if to[k].Action != from[k].Action {
						t.Fail()
					}
					if to[k].Label != from[k].Label {
						t.Fail()
					}
				}
			case "copystruct":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(AccessRolePerms2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				t.Logf("to.Role:%v from.Role: %v\n", *to.Role, from.Role)
				t.Logf("to.Roll:%v from.Roll: %v\n", to.Roll, from.Roll)
				t.Logf("to.EmbedF1:%v from.EmbedF1: %v\n", to.EmbedF1, from.EmbedF1)
				t.Logf("to.Actions:%v from.Actions: %v\n", to.Actions, from.Actions)
				t.Logf("to.Perms:%v from.Perms: %v\n", to.Perms, from.Perms)
				t.Logf("to.PermMap:%v from.PermMap: %v\n", to.PermMap, from.PermMap)
				if *to.Role != from.Role {
					t.Fail()
				}
				if to.Roll != from.Roll {
					t.Fail()
				}
				if to.EmbedF1 != from.EmbedF1 {
					t.Fail()
				}
				for i := 0; i < len(to.Actions); i++ {
					if to.Actions[i] != from.Actions[i] {
						t.Fail()
					}
				}
				for i := 0; i < len(to.Perms); i++ {
					if to.Perms[i].Action != from.Perms[i].Action {
						t.Fail()
					}
					if to.Perms[i].Label != from.Perms[i].Label {
						t.Fail()
					}
				}
				for k, v := range to.PermMap {
					if v.Action != from.PermMap[k].Action {
						t.Fail()
					}
					if v.Label != from.PermMap[k].Label {
						t.Fail()
					}
				}
			}
		})
	}
}

func TestCopyWithOption(t *testing.T) {
	type args struct {
		from interface{}
		to   interface{}
		opt  *Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copyslice-withoption",
			args: args{
				from: []int{3, 4, 5},
				to:   []int{1, 2},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-withoption",
			args: args{
				from: map[string][]int{"key1": {3, 4}, "key2": {7, 8}},
				to:   map[string][]int{"key1": {1, 2}, "key3": {5, 6}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copystruct-withoption",
			args: args{
				from: AccessRolePerms1{
					Id1:     bson.NewObjectId(),
					Id2:     primitive.NewObjectID(),
					Id1Hex:  "61f04828eb37b662c8f3b085",
					Id2Hex:  "61f04828eb37b662c8f3b085",
					Actions: []string{"PUT", "DELETE"},
					From:    "fromtofield",
				},
				to: AccessRolePerms2{
					Actions: []string{"GET", "POST"},
				},
				opt: &Option{
					ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
					StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
					Append:           true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copyslice-withoption":
				from, ok := tt.args.from.([]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]int)
				if !ok {
					t.Fail()
				}
				CopyWithOption(&to, from, tt.args.opt)
				fromto := append(to, from...)
				for i := 0; i < len(to); i++ {
					if to[i] != fromto[i] {
						t.Fail()
					}
				}
			case "copymap-withoption":
				from, ok := tt.args.from.(map[string][]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string][]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string][]int)
				for k, v := range to {
					fromto[k] = append(fromto[k], v...)
				}
				for k, v := range from {
					fromto[k] = append(fromto[k], v...)
				}
				t.Logf("fromto:%#v\n", fromto)

				CopyWithOption(&to, from, tt.args.opt)

				for k, v := range fromto {
					for i := 0; i < len(v); i++ {
						if to[k][i] != fromto[k][i] {
							t.Fail()
						}
					}
				}
			case "copystruct-withoption":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(AccessRolePerms2)
				if !ok {
					t.Fail()
				}
				actions := append(to.Actions, from.Actions...)

				CopyWithOption(&to, from, tt.args.opt)

				if from.Id1.Hex() != to.Id1 {
					t.Fail()
				}
				if from.Id2.Hex() != to.Id2 {
					t.Fail()
				}
				if bson.ObjectIdHex(from.Id1Hex) != to.Id1Hex {
					t.Fail()
				}
				id2hex, err := primitive.ObjectIDFromHex(from.Id2Hex)
				if err != nil {
					t.Fail()
				}
				if id2hex != to.Id2Hex {
					t.Fail()
				}
				for i, v := range to.Actions {
					if actions[i] != v {
						t.Fail()
					}
				}
			}
		})
	}
}

func BenchmarkCopy(b *testing.B) {
	roll := 100
	from := AccessRolePerms1{
		Role: "角色",
		Roll: &roll,
		EmbedFields: EmbedFields{
			EmbedF1: "embedF1",
		},
		Actions: []string{"xxx", "yyy"},
		Perms:   []*Perm1{{Action: "GET", Label: "rest-get-method"}},
		PermMap: map[string]*Perm1{"perm": {Action: "PUT", Label: "rest-put-method"}},
	}
	to := AccessRolePerms2{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&to, from)
	}
}
func BenchmarkCopier(b *testing.B) {
	roll := 100
	from := AccessRolePerms1{
		Role: "角色",
		Roll: &roll,
		EmbedFields: EmbedFields{
			EmbedF1: "embedF1",
		},
		Actions: []string{"xxx", "yyy"},
		Perms:   []*Perm1{{Action: "GET", Label: "rest-get-method"}},
		PermMap: map[string]*Perm1{"perm": {Action: "PUT", Label: "rest-put-method"}},
	}
	to := AccessRolePerms2{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copier.Copy(&to, from)
		// copier.CopyWithOption(&to, from, copier.Option{DeepCopy: true})
	}
}
