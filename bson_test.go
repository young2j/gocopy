/*
 * File: bson_test.go
 * Created Date: 2022-03-22 10:08:30
 * Author: ysj
 * Description:
 */

package gocopy

import (
	"fmt"
	"testing"

	"github.com/globalsign/mgo/bson"
)

func Test_bson(t *testing.T) {
	from := AccessRolePerms1{
		Perms: []*Perm1{{Action: "GET", Label: "get"}},
	}
	fromM := bson.M{}

	CopyWithOption(&fromM, from, &Option{
		IgnoreZero: true,
		Append:     true,
	})
	fmt.Printf("fromM: %#v\n", fromM)

	fromMset := bson.M{"$set": fromM}
	to := bson.M{"$set": bson.M{"tofield": 0}}
	CopyWithOption(&to, fromMset, &Option{
		IgnoreZero: true,
		Append:     true,
	})
	fmt.Printf("to: %#v\n", to)

	fromMap := map[string]map[string]int{"f": {"ff": 1}}
	toMap := make(map[string]interface{})
	Copy(&toMap, fromMap)
	fmt.Printf("toMap: %#v\n", toMap)
}
