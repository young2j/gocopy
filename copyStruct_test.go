/*
 * File: copyStruct.go
 * Created Date: 2022-01-26 06:15:08
 * Author: ysj
 * Description:  copy struct to struct
 */
package gocopy

import (
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_copyStruct(t *testing.T) {
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
			name: "copystruct",
			args: args{
				from: AccessRolePerms1{
					Role: "角色",
					Roll: &roll,
					EmbedFields: EmbedFields{
						EmbedF1: "embedF1",
					},
					Actions: []string{"xxx", "yyy"},
					Perms:   []*Perm1{{Action: "GET", Label: "rest-get-method"}},
					PermMap: map[string]*Perm1{"perm": {Action: "PUT", Label: "rest-put-method"}},
				},
				to: AccessRolePerms2{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copystruct":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(AccessRolePerms2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if *to.Role != from.Role {
					t.Fail()
				}
				if to.Roll != from.Roll {
					t.Fail()
				}
				if to.EmbedF1 != from.EmbedF1 {
					t.Fail()
				}

				if len(to.Actions) != len(from.Actions) {
					t.Fail()
				}
				for i := 0; i < len(to.Actions); i++ {
					if to.Actions[i] != from.Actions[i] {
						t.Fail()
					}
				}

				if len(to.Perms) != len(from.Perms) {
					t.Fail()
				}
				for i := 0; i < len(to.Perms); i++ {
					if to.Perms[i].Action != from.Perms[i].Action {
						t.Fail()
					}
					if to.Perms[i].Label != from.Perms[i].Label {
						t.Fail()
					}
				}

				for k, v := range from.PermMap {
					toPerm, ok := to.PermMap[k]
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

func Test_copyStructWithOption(t *testing.T) {
	type args struct {
		from interface{}
		to   interface{}
		opt  *Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copystruct",
			args: args{
				from: AccessRolePerms1{
					Id1:       bson.NewObjectId(),
					Id2:       primitive.NewObjectID(),
					Id1Hex:    "61f04828eb37b662c8f3b085",
					Id2Hex:    "61f04828eb37b662c8f3b085",
					Actions:   []string{"PUT", "DELETE"},
					From:      "fromtofield",
					CreatedAt: time.Now(),
					UpdatedAt: "2022/01/01",
				},
				to: AccessRolePerms2{
					Actions: []string{"GET", "POST"},
				},
				opt: &Option{
					ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
					StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
					Append:           true,
					NameFromTo:       map[string]string{"From": "To"},
					StringToTime:     map[string]map[string]string{"UpdatedAt": {"loc": "Asia/Shanghai", "layout": "2006/01/02"}},
					TimeToString:     map[string]map[string]string{"CreatedAt": nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copystruct":
				from, ok := tt.args.from.(AccessRolePerms1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(AccessRolePerms2)
				if !ok {
					t.Fail()
				}
				actions := append(to.Actions, from.Actions...)

				CopyWithOption(&to, from, tt.args.opt)

				if from.From != string(to.To) {
					t.Fail()
				}

				if from.Id1.Hex() != *to.Id1 {
					t.Fail()
				}
				if from.Id2.Hex() != to.Id2 {
					t.Fail()
				}
				if bson.ObjectIdHex(from.Id1Hex) != to.Id1Hex {
					t.Fail()
				}
				id2hex, err := primitive.ObjectIDFromHex(from.Id2Hex)
				if err != nil {
					t.Fail()
				}
				if id2hex != to.Id2Hex {
					t.Fail()
				}

				if len(to.Actions) != len(actions) {
					t.Fail()
				}
				for i, v := range to.Actions {
					if actions[i] != v {
						t.Fail()
					}
				}

				fromTimeStr := from.CreatedAt.Format("2006-01-02 15:04:05")
				if fromTimeStr != to.CreatedAt {
					t.Fail()
				}

				toTimeStr := to.UpdatedAt.Format("2006/01/02")
				if toTimeStr != from.UpdatedAt {
					t.Fail()
				}
			}
		})
	}
}
