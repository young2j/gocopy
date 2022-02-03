/*
 * File: types.go
 * Created Date: 2022-01-29 01:29:56
 * Author: ysj
 * Description:
 */
package types

import (
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Perm struct {
	anonymous1 string
	anonymous2 int
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
