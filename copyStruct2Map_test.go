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
					Child: &AccessRolePerms1{
						Id1Hex: "620b7c65eb37b696fe9eef25",
						Role:   "embedstruct",
					},
				},
				to: bson.M{},
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
				to, ok := tt.args.to.(bson.M)
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

				if len(to["perms"].([]*bson.M)) != len(from.Perms) {
					t.Fail()
				}
				for i := 0; i < len(to["perms"].([]*bson.M)); i++ {
					if (*to["perms"].([]*bson.M)[i])["action"] != from.Perms[i].Action {
						t.Fail()
					}
					if (*to["perms"].([]*bson.M)[i])["label"] != from.Perms[i].Label {
						t.Fail()
					}
				}

				for k, v := range from.PermMap {
					toPerm, ok := to["permMap"].(bson.M)[k]
					if !ok {
						t.Fail()
						break
					}
					toPermV, ok := toPerm.(*bson.M)
					if !ok {
						t.Fail()
					}
					if v.Action != (*toPermV)["action"] {
						t.Fail()
					}
					if v.Label != (*toPermV)["label"] {
						t.Fail()
					}
				}

				toChild, ok := to["child"].(*bson.M)
				if !ok {
					t.Fail()
				}
				if (*toChild)["id1Hex"] != "620b7c65eb37b696fe9eef25" {
					t.Fail()
				}
				if (*toChild)["role"] != "embedstruct" {
					t.Fail()
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
	targetPerms := []*map[string]interface{}{
		{"action": "GET", "label": "rest-get-method"},
		{"action": "PUT", "label": "rest-put-method"},
	}
	targetPermMap := map[string]*map[string]interface{}{
		"get": {"action": "GET", "label": "rest-get-method"},
		"put": {"action": "PUT", "label": "rest-put-method"},
	}
	targetChild := &map[string]interface{}{
		"to":     "child",
		"id1Hex": "620b7c65eb37b696fe9eef25",
		"role":   "embedstruct",
	}
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
					PermMap:   map[string]*Perm1{"put": {Action: "PUT", Label: "rest-put-method"}},
					Child: &AccessRolePerms1{
						Id1Hex: "620b7c65eb37b696fe9eef25",
						Role:   "embedstruct",
					},
				},
				to: map[interface{}]interface{}{
					"perms":   []*Perm1{{Action: "GET", Label: "rest-get-method"}},
					"permMap": map[interface{}]*Perm1{"get": {Action: "GET", Label: "rest-get-method"}},
					"child":   &AccessRolePerms1{From: "child"},
				},
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
		{
			name: "copystruct2map-convert",
			args: args{
				from: AccessRolePerms1{
					CreatedAt: time.Now(),
					UpdatedAt: "2022/02/16",
					Id1:       bson.NewObjectId(),
					Id2:       primitive.NewObjectID(),
					Id1Hex:    "61f04828eb37b662c8f3b085",
					Id2Hex:    "61f04828eb37b662c8f3b085",
				},
				to: bson.M{},
				opt: &Option{
					Converters: map[string]func(interface{}) interface{}{
						"CreatedAt": func(v interface{}) interface{} {
							return v.(time.Time).Format("2006-01-02 15:04:05")
						},
						"UpdatedAt": func(v interface{}) interface{} {
							t, _ := time.Parse("2006/01/02", v.(string))
							return t
						},
						"Id1": func(v interface{}) interface{} {
							return v.(bson.ObjectId).Hex()
						},
						"Id2": func(v interface{}) interface{} {
							return v.(primitive.ObjectID).Hex()
						},
						"Id3": func(v interface{}) interface{} {
							return v.(*primitive.ObjectID).Hex()
						},
						"Id1Hex": func(v interface{}) interface{} {
							return bson.ObjectIdHex(v.(string))
						},
						"Id2Hex": func(v interface{}) interface{} {
							oid, _ := primitive.ObjectIDFromHex(v.(string))
							return oid
						},
					},
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

				toPerms, ok := to["perms"].([]*map[interface{}]interface{})
				if !ok {
					t.Fail()
				}
				if len(toPerms) != len(targetPerms) {
					t.Fail()
				}
				for i := 0; i < len(toPerms); i++ {

					if (*toPerms[i])["action"] != (*targetPerms[i])["action"] {
						t.Fail()
					}
					if (*toPerms[i])["label"] != (*targetPerms[i])["label"] {
						t.Fail()
					}
				}

				toPermMap, ok := to["permMap"].(map[interface{}]interface{})
				if !ok {
					t.Fail()
				}
				if len(toPermMap) != len(targetPermMap) {
					t.Fail()
				}
				for k, v := range targetPermMap {
					toPermV, ok := toPermMap[k]
					if !ok {
						t.Fail()
					}

					toPermV_, ok := toPermV.(*map[interface{}]interface{})
					if !ok {
						t.Fail()
					}
					for kk, vv := range *v {
						toPermVV, ok := (*toPermV_)[kk]
						if !ok {
							t.Fail()
						}
						if toPermVV != vv {
							t.Fail()
						}
					}
				}

				toChild, ok := to["child"].(*map[interface{}]interface{})
				if !ok {
					t.Fail()
				}
				for k, v := range *targetChild {
					toChildV, ok := (*toChild)[k]
					if !ok {
						t.Fail()
					}
					if (k == "to" || k == "role") && toChildV != v {
						t.Fail()
					} else if k == "id1Hex" {
						if toChildV != bson.ObjectIdHex("620b7c65eb37b696fe9eef25") {
							t.Fail()
						}
					}
				}
				// ignore
				if _, ok := to["embedF1"]; ok {
					t.Fail()
				}
			case "copystruct2map-convert":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(bson.M)
				if !ok {
					t.Fail()
				}

				CopyWithOption(&to, from, tt.args.opt)

				if to["createdAt"] != from.CreatedAt.Format("2006-01-02 15:04:05") {
					t.Fail()
				}

				if to["updatedAt"].(time.Time).Format("2006/01/02") != from.UpdatedAt {
					t.Fail()
				}

				if to["id1"] != from.Id1.Hex() {
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
			}
		})
	}
}
