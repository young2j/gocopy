/*
 * File: model.go
 * Created Date: 2022-01-29 01:29:06
 * Author: ysj
 * Description:
 */

package model

import (
	"time"

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
	CreatedAt  time.Time
	UpdatedAt  string
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
	Child *AccessRolePerms
}
