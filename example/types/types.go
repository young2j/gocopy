/*
 * File: types.go
 * Created Date: 2022-01-29 01:29:56
 * Author: ysj
 * Description:
 */
package types

import "github.com/globalsign/mgo/bson"

type Perm struct {
	Action string
	Label  string
}

type AccessRolePerms struct {
	Id      string
	IdHex   bson.ObjectId
	Role    *string
	Roll    *int
	Actions []string
	Perms   []*Perm
	PermMap map[string]*Perm
	To      []byte
	EmbedF1 string
}
