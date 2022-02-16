/*
 * File: copyStruct2Map.go
 * Created Date: 2022-02-16 12:54:13
 * Author: ysj
 * Description:
 */

package gocopy

import (
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_copyStruct2Map(t *testing.T) {
	type args struct {
		to   interface{}
		from interface{}
	}
	roll := 100

	tests := []struct {
		name string
		args args
	}{
		{
			name: "copystruct2map",
			args: args{
				from: AccessRolePerms1{
					CreatedAt: time.Now(),
					UpdatedAt: "2022/02/16",
					Id1:       bson.NewObjectId(),
					Id2:       primitive.NewObjectID(),
					Id1Hex:    bson.NewObjectId().Hex(),
					Id2Hex:    primitive.NewObjectID().Hex(),
					Role:      "copystruct2map",
					Roll:      &roll,
					From:      "From",
					Actions:   []string{"PUT", "DELETE"},
					Perms:     []*Perm1{{Action: "PUT", Label: "rest-put-method"}},
					PermMap:   map[string]*Perm1{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
				},
				to: make(map[interface{}]interface{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copystruct2map":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[interface{}]interface{})
				if !ok {
					t.Fail()
				}
				Copy(&to, from)
				if to["createdAt"] != from.CreatedAt {
					t.Fail()
				}
				if to["updatedAt"] != from.UpdatedAt {
					t.Fail()
				}

				if to["id1"] != from.Id1 {
					t.Fail()
				}
				if to["id2"] != from.Id2 {
					t.Fail()
				}

				if to["id1Hex"] != from.Id1Hex {
					t.Fail()
				}

				if to["id2Hex"] != from.Id2Hex {
					t.Fail()
				}
				if to["from"] != from.From {
					t.Fail()
				}

				if to["role"] != from.Role {
					t.Fail()
				}
				if *to["roll"].(*int) != *from.Roll {
					t.Fail()
				}
				if to["embedF1"] != from.EmbedF1 {
					t.Fail()
				}

				if len(to["actions"].([]string)) != len(from.Actions) {
					t.Fail()
				}
				for i := 0; i < len(to["actions"].([]string)); i++ {
					if to["actions"].([]string)[i] != from.Actions[i] {
						t.Fail()
					}
				}

				if len(to["perms"].([]*Perm1)) != len(from.Perms) {
					t.Fail()
				}
				for i := 0; i < len(to["perms"].([]*Perm1)); i++ {
					if to["perms"].([]*Perm1)[i].Action != from.Perms[i].Action {
						t.Fail()
					}
					if to["perms"].([]*Perm1)[i].Label != from.Perms[i].Label {
						t.Fail()
					}
				}

				for k, v := range from.PermMap {
					toPerm, ok := to["permMap"].(map[string]*Perm1)[k]
					if !ok {
						t.Fail()
						break
					}
					if v.Action != toPerm.Action {
						t.Fail()
					}
					if v.Label != toPerm.Label {
						t.Fail()
					}
				}
			}
		})
	}
}
func Test_copyStruct2MapWithOption(t *testing.T) {
	type args struct {
		to   interface{}
		from interface{}
		opt  *Option
	}
	roll := 100

	tests := []struct {
		name string
		args args
	}{
		{
			name: "copystruct2map",
			args: args{
				from: AccessRolePerms1{
					UpdatedAt: "2022/02/16",
					CreatedAt: time.Now(),
					Id1:       bson.NewObjectId(),
					Id2:       primitive.NewObjectID(),
					Id1Hex:    bson.NewObjectId().Hex(),
					Id2Hex:    primitive.NewObjectID().Hex(),
					Role:      "copystruct2map",
					Roll:      &roll,
					From:      "From",
					Actions:   []string{"PUT", "DELETE"},
					Perms:     []*Perm1{{Action: "PUT", Label: "rest-put-method"}},
					PermMap:   map[string]*Perm1{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
				},
				to: make(map[interface{}]interface{}),
				opt: &Option{
					Append:           true,
					NameFromTo:       map[string]string{"From": "to", "Id1": "_id"},
					ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
					StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
					TimeToString:     map[string]map[string]string{"CreatedAt": {"layout": "2006-01-02", "loc": "America/New_York"}},
					StringToTime:     map[string]map[string]string{"UpdatedAt": {"layout": "2006/01/02"}},
					IgnoreZero:       true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copystruct2map":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[interface{}]interface{})
				if !ok {
					t.Fail()
				}
				CopyWithOption(&to, from, tt.args.opt)

				loc, err := time.LoadLocation("America/New_York")
				if err != nil {
					t.Fail()
				}
				if to["createdAt"] != from.CreatedAt.In(loc).Format("2006-01-02") {
					t.Fail()
				}

				loc, err = time.LoadLocation("Asia/Shanghai")
				if err != nil {
					t.Fail()
				}
				if to["updatedAt"].(time.Time).In(loc).Format("2006/01/02") != from.UpdatedAt {
					t.Fail()
				}

				if to["_id"] != from.Id1.Hex() {
					t.Fail()
				}
				if to["id2"] != from.Id2.Hex() {
					t.Fail()
				}

				if to["id1Hex"].(bson.ObjectId).Hex() != from.Id1Hex {
					t.Fail()
				}

				if to["id2Hex"].(primitive.ObjectID).Hex() != from.Id2Hex {
					t.Fail()
				}
				if to["to"] != from.From {
					t.Fail()
				}

				if to["role"] != from.Role {
					t.Fail()
				}
				if *to["roll"].(*int) != *from.Roll {
					t.Fail()
				}

				if len(to["actions"].([]string)) != len(from.Actions) {
					t.Fail()
				}
				for i := 0; i < len(to["actions"].([]string)); i++ {
					if to["actions"].([]string)[i] != from.Actions[i] {
						t.Fail()
					}
				}

				if len(to["perms"].([]*Perm1)) != len(from.Perms) {
					t.Fail()
				}
				for i := 0; i < len(to["perms"].([]*Perm1)); i++ {
					if to["perms"].([]*Perm1)[i].Action != from.Perms[i].Action {
						t.Fail()
					}
					if to["perms"].([]*Perm1)[i].Label != from.Perms[i].Label {
						t.Fail()
					}
				}

				for k, v := range from.PermMap {
					toPerm, ok := to["permMap"].(map[string]*Perm1)[k]
					if !ok {
						t.Fail()
						break
					}
					if v.Action != toPerm.Action {
						t.Fail()
					}
					if v.Label != toPerm.Label {
						t.Fail()
					}
				}

				if _, ok := to["embedF1"]; ok {
					t.Fail()
				}
			}
		})
	}
}
