/*
 * File: types.go
 * Created Date: 2022-01-29 01:29:56
 * Author: ysj
 * Description:
 */
package types

import (
	"github.com/globalsign/mgo/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type Perm struct {
	Action string
	Label  string
}

type AccessRolePerms struct {
	Id string
	IdHex   bson.ObjectId
	// IdHex   primitive.ObjectID
	Role    *string
	Roll    *int
	Actions []string
	Perms   []*Perm
	PermMap map[string]*Perm
	To      []byte
	EmbedF1 string
}
