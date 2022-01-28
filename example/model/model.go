/*
 * File: model.go
 * Created Date: 2022-01-29 01:29:06
 * Author: ysj
 * Description:
 */

package model

import "github.com/globalsign/mgo/bson"

type Perm struct {
	Action string
	Label  string
}

type EmbedFields struct {
	EmbedF1 string
	embedF2 string
}

type AccessRolePerms struct {
	Id       bson.ObjectId
	IdHex    string
	Role     string
	Roll     *int
	Actions  []string
	Perms    []*Perm
	PermMap  map[string]*Perm
	From     string
	EmbedFields
}
