// Package model contains source types used by the runnable examples.
package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Perm is a permission in the source model.
type Perm struct {
	Action string
	Label  string
}

// EmbedFields demonstrates promotion of anonymous embedded fields.
type EmbedFields struct {
	EmbedF1 string
	embedF2 string
}

// AccessRolePerms exercises nested pointers, slices, maps, and time values.
type AccessRolePerms struct {
	CreatedAt time.Time
	UpdatedAt string
	Id1       bson.ObjectId
	Id2       primitive.ObjectID
	Id3       *primitive.ObjectID
	Id1Hex    string
	Id2Hex    string
	Role      string
	Roll      *int
	Actions   []string
	Perms     []*Perm
	PermMap   map[string]*Perm
	From      string
	EmbedFields
	Child *AccessRolePerms
}
