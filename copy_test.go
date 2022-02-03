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
	anonymous1 string
	anonymous2 int
	Action     string
	Label      string
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
