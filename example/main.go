/*
 * File: test.go
 * Created Date: 2022-01-29 03:05:53
 * Author: ysj
 * Description: test cases
 */

package main

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/young2j/gocopy"
	"github.com/young2j/gocopy/example/model"
	"github.com/young2j/gocopy/example/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// slice to slice
	// 1. simple slice
	s1 := []int{3, 4, 5}
	s2 := make([]int, 0)
	gocopy.Copy(&s2, nil)
	gocopy.Copy(&s2, &s1)
	// copier.Copy(&s2, &s1)
	fmt.Println("==============================")
	fmt.Printf("simple slice-> s2: %v\n", s2)

	// 2. map slice
	ms1 := []map[string]int{{"key1": 1, "key2": 2}}
	ms2 := make([]map[string]int, 0)
	gocopy.Copy(&ms2, &ms1)
	// copier.Copy(&ms2, &ms1)
	fmt.Println("==============================")
	fmt.Printf("map slice-> ms2: %v\n", ms2)

	// 2.1 ptr map slice
	pms1 := []*map[string]int{{"key1": 1, "key2": 2}}
	pms2 := make([]*map[string]int, 0)
	gocopy.Copy(&pms2, &pms1)
	// copier.Copy(&pms2, &pms1)
	fmt.Printf("ptr map slice-> pms2: %v\n", pms2)

	// 2.2 from ptr map slice
	ms2 = make([]map[string]int, 0)
	gocopy.Copy(&ms2, &pms1)
	// copier.Copy(&ms2, &pms1)
	fmt.Printf("from ptr map slice-> ms2: %v\n", ms2)

	// 2.2 to ptr map slice
	pms2 = make([]*map[string]int, 0)
	gocopy.Copy(&pms2, &ms1)
	// copier.Copy(&pms2, &ms1)
	fmt.Printf("to ptr map slice-> pms2: %#v\n", pms2)

	// 3. struct slice
	ss1 := []model.Perm{{Action: "GET", Label: "rest-get-method"}}
	ss2 := make([]types.Perm, 0)
	gocopy.Copy(&ss2, &ss1)
	// copier.Copy(&ss2, &ss1)
	fmt.Println("==============================")
	fmt.Printf("struct slice-> ss2: %#v\n", ss2)

	// 3.1 ptr struct slice
	pss1 := []*model.Perm{{Action: "GET", Label: "rest-get-method"}}
	pss2 := make([]*types.Perm, 0)
	gocopy.Copy(&pss2, &pss1)
	// copier.Copy(&pss2, &pss1)
	for _, v := range pss2 {
		fmt.Printf("ptr struct slice-> pss2: %#v\n", v)
	}

	// 3.2 from ptr struct slice
	ss2 = make([]types.Perm, 0)
	gocopy.Copy(&ss2, &pss1)
	// copier.Copy(&ss2, &pss1)
	for _, v := range ss2 {
		fmt.Printf("from ptr struct slice-> ptrss2: %#v\n", v)
	}

	// 3.3 to ptr struct slice
	pss2 = make([]*types.Perm, 0)
	gocopy.Copy(&pss2, &ss1)
	// copier.Copy(&pss2, &ss1)
	for _, v := range ss2 {
		fmt.Printf("to ptr struct slice-> ptrss2: %#v\n", v)
	}

	// map to map
	// 1. simple map
	m1 := map[string]int{"key1": 1, "key2": 2}
	m2 := make(map[string]int)
	gocopy.Copy(&m2, &m1)
	// copier.Copy(&m2, &m1)
	fmt.Println("==============================")
	fmt.Printf("simple map-> m2: %v\n", m2)

	// 2. slice map
	sm1 := map[int][]string{1: {"a", "b", "c"}}
	sm2 := make(map[int][]string)
	gocopy.Copy(&sm2, &sm1)
	// copier.Copy(&sm2, &sm1)
	fmt.Printf("slice map-> sm2: %v\n", sm2)

	// 2.1 ptr slice map
	psm1 := map[int]*[]string{1: {"a", "b", "c"}}
	psm2 := make(map[int]*[]string)
	gocopy.Copy(&psm2, &psm1)
	// copier.Copy(&psm2, &psm1) // fail
	fmt.Printf("ptr slice map-> psm2: %#v\n", psm2)

	// 2.2 from ptr slice map
	sm2 = make(map[int][]string)
	gocopy.Copy(&sm2, &psm1)
	// copier.Copy(&sm2, &psm1)
	fmt.Printf("from ptr slice map-> sm2: %#v\n", sm2)

	// 2.3 to ptr slice map
	psm2 = make(map[int]*[]string)
	gocopy.Copy(&psm2, &sm1)
	// copier.Copy(&psm2, &sm1) // fail
	fmt.Printf("to ptr slice map-> psm2: %#v\n", psm2)

	// 3. map map
	mm1 := map[string]map[string]int{"level1k": {"level2k": 1}}
	mm2 := make(map[string]map[string]int)
	gocopy.Copy(&mm2, &mm1)
	// copier.Copy(&mm2, &mm1)
	fmt.Println("==============================")
	fmt.Printf("map map-> mm2: %v\n", mm2)

	// 3.1 ptr map map
	pmm1 := map[string]*map[string]int{"level1k": {"level2k": 1}}
	pmm2 := make(map[string]*map[string]int)
	gocopy.Copy(&pmm2, &pmm1)
	// copier.Copy(&pmm2, &pmm1)
	fmt.Printf("ptr map map-> pmm2: %v\n", pmm2)

	// 3.2 from ptr map map
	mm2 = make(map[string]map[string]int)
	gocopy.Copy(&mm2, &pmm1)
	// copier.Copy(&mm2, &pmm1)
	fmt.Printf("from ptr map map-> mm2: %v\n", mm2)

	// 3.3 to ptr map map
	pmm2 = make(map[string]*map[string]int)
	gocopy.Copy(&pmm2, &mm1)
	// copier.Copy(&pmm2, &mm1)
	fmt.Printf("to ptr map map-> pmm2: %v\n", pmm2)

	// 3. struct map
	stm1 := map[string]model.Perm{"perm1": {Action: "POST", Label: "rest-post-method"}}
	stm2 := make(map[string]types.Perm)
	gocopy.Copy(&stm2, &stm1)
	// copier.Copy(&stm2, &stm1)
	fmt.Println("==============================")
	fmt.Printf("struct map-> stm2: %v\n", stm2)

	// 3.1 ptr struct map
	pstm1 := map[string]*model.Perm{"perm1": {Action: "POST", Label: "rest-post-method"}}
	pstm2 := make(map[string]*types.Perm)
	gocopy.Copy(&pstm2, &pstm1)
	// copier.Copy(&pstm2, &pstm1)
	fmt.Printf("ptr struct map-> pstm2: %v\n", pstm2)

	// 3.2 from ptr struct map
	stm2 = make(map[string]types.Perm)
	gocopy.Copy(&stm2, &pstm1)
	// copier.Copy(&stm2, &pstm1)
	fmt.Printf("from ptr struct map-> stm2: %v\n", stm2)

	// 3.3 to ptr struct map
	pstm2 = make(map[string]*types.Perm)
	gocopy.Copy(&pstm2, &stm1)
	// copier.Copy(&pstm2, &stm1)
	fmt.Printf("to ptr struct map-> pstm2: %v\n", pstm2)

	// 4. struct to struct
	roll := 100
	st1 := model.AccessRolePerms{
		Role: "角色",
		Roll: &roll,
		EmbedFields: model.EmbedFields{
			EmbedF1: "embedF1",
		},
		Actions: []string{"GET", "POST"},
		Perms:   []*model.Perm{{Action: "GET", Label: "rest-get-method"}},
		PermMap: map[string]*model.Perm{"perm": {Action: "PUT", Label: "rest-put-method"}},
	}
	st2 := types.AccessRolePerms{}
	gocopy.Copy(&st2, &st1)
	// copier.Copy(&st2, &st1)
	fmt.Println("==============================")
	fmt.Println("struct to struct->")
	// fmt.Printf("st2: %#v\n", st2)
	fmt.Printf("st2.Role: %v\n", *st2.Role)
	fmt.Printf("st2.Roll: %v\n", *st2.Roll)
	fmt.Printf("st2.Actions: %v\n", st2.Actions)

	for _, v := range st2.Perms {
		fmt.Printf("Perms: %#v\n", v)
	}
	for k, v := range st2.PermMap {
		fmt.Printf("PermMap k:%v v:%#v\n", k, v)
	}

	// append mode
	opts := gocopy.Option{Append: true}
	// 1. append simple slice
	as1 := []int{3, 4, 5}
	as2 := []int{1, 2}
	gocopy.CopyWithOption(&as2, &as1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append slice-> as2: %v\n", as2)

	// 2. append slice map
	asm1 := []map[string]int{{"k1": 1}}
	asm2 := []map[string]int{{"k2": 2}}
	gocopy.CopyWithOption(&asm2, &asm1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append slice map-> asm2: %v\n", asm2)

	// 2.1 append slice ptr map
	apsm1 := []*map[string]int{{"k1": 1}}
	apsm2 := []*map[string]int{{"k2": 2}}
	gocopy.CopyWithOption(&apsm2, &apsm1, &opts)
	fmt.Printf("append slice ptr map-> apsm2: %v\n", apsm2)

	// 2.2 append from slice ptr map
	asm2 = []map[string]int{{"k2": 2}}
	gocopy.CopyWithOption(&asm2, &apsm1, &opts)
	fmt.Printf("append from slice ptr map-> asm2: %v\n", asm2)

	// 2.3 append to slice ptr map
	apsm2 = []*map[string]int{{"k2": 2}}
	gocopy.CopyWithOption(&apsm2, &asm1, &opts)
	fmt.Printf("append from slice ptr map-> apsm2: %v\n", apsm2)

	// 3. append slice struct
	asst1 := []model.Perm{{Action: "GET", Label: "rest-get-method"}}
	asst2 := []types.Perm{{Action: "PUT", Label: "rest-put-method"}}
	gocopy.CopyWithOption(&asst2, &asst1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append slice struct-> asst2: %v\n", asst2)

	// 3.1 append slice ptr struct
	apsst1 := []*model.Perm{{Action: "GET", Label: "rest-get-method"}}
	apsst2 := []*types.Perm{{Action: "PUT", Label: "rest-put-method"}}
	gocopy.CopyWithOption(&apsst2, &apsst1, &opts)
	fmt.Printf("append slice ptr struct-> apsst2: %v\n", apsst2)

	// 3.2 append from slice ptr struct
	asst2 = []types.Perm{{Action: "PUT", Label: "rest-put-method"}}
	gocopy.CopyWithOption(&asst2, &apsst1, &opts)
	fmt.Printf("append from slice ptr struct-> asst2: %v\n", asst2)

	// 3.3 append to slice ptr struct
	apsst2 = []*types.Perm{{Action: "PUT", Label: "rest-put-method"}}
	gocopy.CopyWithOption(&apsst2, &asst1, &opts)
	fmt.Printf("append to slice ptr struct-> apsst2: %v\n", apsst2)

	// 4. append simple map
	am1 := map[string]int{"key1": 1, "key2": 2}
	am2 := map[string]int{"key0": 0, "key2": 3}
	gocopy.CopyWithOption(&am2, &am1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append simple map-> am2: %v\n", am2)

	// 5. append map slice
	ams1 := map[string][]int{"key1": {1}, "key2": {2}}
	ams2 := map[string][]int{"key0": {0}, "key2": {3}}
	gocopy.CopyWithOption(&ams2, &ams1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append map slice-> ams2: %v\n", ams2)

	// 5.1 append map ptr slice
	apms1 := map[string]*[]int{"key1": {1}, "key2": {2}}
	apms2 := map[string]*[]int{"key0": {0}, "key2": {3}}
	gocopy.CopyWithOption(&apms2, &apms1, &opts)
	fmt.Printf("append map ptr slice-> apms2: %v\n", apms2)

	// 5.2 append from map ptr slice
	ams2 = map[string][]int{"key0": {0}, "key2": {3}}
	gocopy.CopyWithOption(&ams2, &apms1, &opts)
	fmt.Printf("append from map ptr slice-> ams2: %v\n", ams2)

	// 5.3 append to map ptr slice
	gocopy.CopyWithOption(&apms2, &ams1, &opts)
	fmt.Printf("append to map ptr slice-> apms2: %v\n", apms2)

	// 6. append map map
	amm1 := map[string]map[string]int{"level1k": {"level2k": 1}}
	amm2 := map[string]map[string]int{"level1k": {"level2k_": 2}}
	gocopy.CopyWithOption(&amm2, &amm1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append map map-> amm2: %#v\n", amm2)

	// 6.1 append map ptr map
	apmm1 := map[string]*map[string]int{"level1k": {"level2k": 1}}
	apmm2 := map[string]*map[string]int{"level1k": {"level2k_": 2}}
	gocopy.CopyWithOption(&apmm2, &apmm1, &opts)
	fmt.Printf("append map ptr map-> apmm2: %#v\n", apmm2)
	// 6.2 append from map ptr map
	amm2 = map[string]map[string]int{"level1k": {"level2k_": 2}}
	gocopy.CopyWithOption(&amm2, &apmm1, &opts)
	fmt.Printf("append from map ptr map-> amm2: %#v\n", amm2)

	// 6.3 append to map ptr map
	apmm2 = map[string]*map[string]int{"level1k": {"level2k_": 2}}
	gocopy.CopyWithOption(&apmm2, &amm1, &opts)
	fmt.Printf("append to map ptr map-> apmm2: %#v\n", apmm2)

	// 7. append map struct
	amst1 := map[string]model.Perm{"perm1": {Action: "PUT", Label: "rest-put-method"}}
	amst2 := map[string]types.Perm{"perm2": {Action: "GET", Label: "rest-get-method"}}
	gocopy.CopyWithOption(&amst2, &amst1, &opts)
	fmt.Println("==============================")
	fmt.Printf("append map struct-> amst2: %v\n", amst2)

	// 7.1 append map ptr struct
	apmst1 := map[string]*model.Perm{"perm1": {Action: "PUT", Label: "rest-put-method"}}
	apmst2 := map[string]*types.Perm{"perm2": {Action: "GET", Label: "rest-get-method"}}
	gocopy.CopyWithOption(&apmst2, &apmst1, &opts)
	fmt.Printf("append map ptr struct-> apmst2: %v\n", apmst2)

	// 7.2 append from map ptr struct
	amst2 = map[string]types.Perm{"perm2": {Action: "GET", Label: "rest-get-method"}}
	gocopy.CopyWithOption(&amst2, &apmst1, &opts)
	fmt.Printf("append from map ptr struct-> amst2: %v\n", amst2)

	// 7.3 append to map ptr struct
	apmst2 = map[string]*types.Perm{"perm2": {Action: "GET", Label: "rest-get-method"}}
	gocopy.CopyWithOption(&apmst2, &amst1, &opts)
	fmt.Printf("append to map ptr struct-> apmst2: %v\n", apmst2)

	// 8. append struct map/slice field
	ast1 := model.AccessRolePerms{
		Actions: []string{"PUT", "DELETE"},
		Perms:   []*model.Perm{{Action: "PUT", Label: "rest-put-method"}},
		PermMap: map[string]*model.Perm{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
	}
	ast2 := types.AccessRolePerms{
		Actions: []string{"GET", "POST"},
		Perms:   []*types.Perm{{Action: "GET", Label: "rest-get-method"}},
		PermMap: map[string]*types.Perm{"get": {Action: "GET", Label: "rest-get-method"}},
	}
	gocopy.CopyWithOption(&ast2, &ast1, &opts)
	fmt.Println("==============================")
	fmt.Println("append struct map/slice field->")
	fmt.Printf("ast2.Actions: %v\n", ast2.Actions)
	for i, perm := range ast2.Perms {
		fmt.Printf("ast2.Perms[%v]: %#v\n", i, perm)
	}
	for i, pm := range ast2.PermMap {
		fmt.Printf("ast2.PermMap[%v]: %#v\n", i, pm)
	}

	// 9. from field to another field
	ost1 := model.AccessRolePerms{
		From: "fromto",
	}
	ost2 := types.AccessRolePerms{}
	opt := gocopy.Option{
		NameFromTo: map[string]string{"From": "To"},
	}
	gocopy.CopyWithOption(&ost2, &ost1, &opt)
	fmt.Println("==============================")
	fmt.Printf("from field to another field-> ost2.To: %v\n", string(ost2.To))

	// 10. objectId to string and vice versa
	from := model.AccessRolePerms{
		Id1:    bson.NewObjectId(),
		Id2:    primitive.NewObjectID(),
		Id1Hex: "61f04828eb37b662c8f3b085",
		Id2Hex: "61f04828eb37b662c8f3b085",
	}
	to := types.AccessRolePerms{
		Actions: []string{"GET", "POST"},
	}
	option := &gocopy.Option{
		ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
		StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
		Append:           true,
	}
	gocopy.CopyWithOption(&to, from, option)
	fmt.Println("==============================")
	fmt.Println("objectId to string and vice versa->")
	fmt.Printf("from.Id1: %v to.Id1: %v \n", from.Id1, *to.Id1)
	fmt.Printf("from.Id2: %v to.Id2: %v \n", from.Id2, to.Id2)
	fmt.Printf("from.Id1Hex: %v to.Id1Hex: %v \n", from.Id1Hex, to.Id1Hex)
	fmt.Printf("from.Id2Hex: %v to.Id2Hex: %v \n", from.Id2Hex, to.Id2Hex)

	// 11. time.Time to string
	from1 := model.AccessRolePerms{
		CreatedAt: time.Now(),
		UpdatedAt: "2022/02/11 15:04:05",
	}
	to1 := types.AccessRolePerms{}
	option1 := gocopy.Option{
		// TimeToString: map[string]map[string]string{"CreatedAt": nil},
		// StringToTime: map[string]map[string]string{"UpdatedAt": nil},
		TimeToString: map[string]map[string]string{"CreatedAt": {"layout": "2006-01-02", "loc": "America/New_York"}},
		StringToTime: map[string]map[string]string{"UpdatedAt": {"layout": "2006/01/02 15:04:05"}},
	}
	gocopy.CopyWithOption(&to1, from1, &option1)
	fmt.Println("==============================")
	fmt.Printf("time.Time to string-> to1.CreatedAt: %v\n", to1.CreatedAt)
	fmt.Printf("string to time.Time-> to1.UpdatedAt: %v\n", to1.UpdatedAt)

	fromst := model.AccessRolePerms{
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
		Perms:     []*model.Perm{{Action: "PUT", Label: "rest-put-method"}},
		PermMap:   map[string]*model.Perm{"delete": {Action: "DELETE", Label: "rest-delete-method"}},
		Child: &model.AccessRolePerms{
			Id1Hex: bson.NewObjectId().Hex(),
			Role:   "embedstruct",
		},
	}
	// toBM := mbson.M{}
	toBM := bson.M{
		"perms": []*model.Perm{{Action: "GET", Label: "rest-get-method"}},
		"permMap": map[string]*model.Perm{"get": {Action: "GET", Label: "rest-get-method"}},
		// "permMap": bson.M{"get": model.Perm{Action: "GET", Label: "rest-get-method"}},
		// "child": map[string]string{"from":"child"},
		"child": &model.AccessRolePerms{From: "child"},
	}
	// toBM := make(map[interface{}]interface{})
	// v := []string{"GET"}
	// toBM := map[interface{}]interface{}{
	// 	"actions": &v,
	// 	"permMap": map[string]*model.Perm{"put": {Action: "PUT", Label: "rest-put-method"}},
	// }
	// gocopy.Copy(&toBM, from)
	gocopy.CopyWithOption(&toBM, fromst, &gocopy.Option{
		Append:           true,
		NameFromTo:       map[string]string{"From": "to", "Id1": "_id"},
		ObjectIdToString: map[string]string{"Id1": "mgo", "Id2": "official"},       // Id1: bson.ObjectId, Id2: primitive.ObjectId
		StringToObjectId: map[string]string{"Id1Hex": "mgo", "Id2Hex": "official"}, // Id1Hex: bson.ObjectId.Hex(), Id2Hex: primitive.ObjectId.Hex()
		TimeToString:     map[string]map[string]string{"CreatedAt": {"layout": "2006-01-02", "loc": "America/New_York"}},
		StringToTime:     map[string]map[string]string{"UpdatedAt": {"layout": "2006/01/02"}},
		IgnoreZero:       true,
	})
	fmt.Println("==============================")
	fmt.Println("copy struct to map->")
	fmt.Printf("toBM[\"createdAt\"]: %v\n", toBM["createdAt"])
	fmt.Printf("toBM[\"updatedAt\"]: %v\n", toBM["updatedAt"])
	fmt.Printf("toBM[\"id1\"]: %v\n", toBM["id1"])
	fmt.Printf("toBM[\"id2\"]: %v\n", toBM["id2"])
	fmt.Printf("toBM[\"_id\"]: %v\n", toBM["_id"])
	fmt.Printf("toBM[\"id1Hex\"]: %v\n", toBM["id1Hex"])
	fmt.Printf("toBM[\"id2Hex\"]: %v\n", toBM["id2Hex"])
	fmt.Printf("toBM[\"role\"]: %v\n", toBM["role"])
	fmt.Printf("toBM[\"roll\"]: %v\n", *toBM["roll"].(*int))
	fmt.Printf("toBM[\"to\"]: %v\n", toBM["to"])
	fmt.Printf("toBM[\"actions\"]: %v\n", toBM["actions"])
	fmt.Printf("toBM[\"perms\"]: %#v\n", toBM["perms"])
	fmt.Printf("toBM[\"permMap\"]: %#v\n", toBM["permMap"])
	fmt.Printf("toBM[\"child\"]: %#v\n", toBM["child"])
	// fmt.Printf("toBM: %v\n", toBM)
}
