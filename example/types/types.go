// Package types contains destination types used by the runnable examples.
package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Perm is the destination permission shape.
type Perm struct {
	anonymous1 string
	anonymous2 int
	Action     string
	Label      string
}

// AccessRolePerms demonstrates a destination with different field types.
type AccessRolePerms struct {
	Id1       *string
	Id2       string
	Id3       string
	Id1Hex    bson.ObjectId
	Id2Hex    primitive.ObjectID
	Role      *string
	Roll      *int
	Actions   []string
	Perms     []*Perm
	PermMap   map[string]*Perm
	To        []byte
	CreatedAt string
	UpdatedAt time.Time
	EmbedF1   string
}
