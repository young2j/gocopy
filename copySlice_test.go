/*
 * File: copySlice.go
 * Created Date: 2022-01-26 06:14:15
 * Author: ysj
 * Description:  slice to slice copy
 */

package gocopy

import (
	"testing"
)

func Test_copySlice(t *testing.T) {
	type args struct {
		to   interface{}
		from interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copyslice-simple",
			args: args{
				from: []int{3, 4, 5},
				to:   make([]int, 0),
			},
		},
		{
			name: "copyslice-map",
			args: args{
				from: []map[string]int{{"key1": 1, "key2": 2}},
				to:   make([]map[string]int, 0),
			},
		},
		{
			name: "copyslice-ptrmap",
			args: args{
				from: []*map[string]int{{"key1": 1, "key2": 2}},
				to:   make([]*map[string]int, 0),
			},
		},
		{
			name: "copyslice-from-ptrmap",
			args: args{
				from: []*map[string]int{{"key1": 1, "key2": 2}},
				to:   make([]map[string]int, 0),
			},
		},
		{
			name: "copyslice-to-ptrmap",
			args: args{
				from: []map[string]int{{"key1": 1, "key2": 2}},
				to:   make([]*map[string]int, 0),
			},
		},
		{
			name: "copyslice-struct",
			args: args{
				from: []Perm1{{Action: "GET", Label: "rest-get-method"}},
				to:   make([]Perm2, 0),
			},
		},
		{
			name: "copyslice-ptrStruct",
			args: args{
				from: []*Perm1{{Action: "GET", Label: "rest-get-method"}},
				to:   make([]*Perm2, 0),
			},
		},
		{
			name: "copyslice-from-ptrStruct",
			args: args{
				from: []*Perm1{{Action: "GET", Label: "rest-get-method"}},
				to:   make([]Perm2, 0),
			},
		},
		{
			name: "copyslice-to-ptrStruct",
			args: args{
				from: []Perm1{{Action: "GET", Label: "rest-get-method"}},
				to:   make([]*Perm2, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copyslice-simple":
				from, ok := tt.args.from.([]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					t.Logf("from[%d] == to[%d] = %v\n", i, i, from[i] == to[i])
					if from[i] != to[i] {
						t.Fail()
					}
				}
			case "copyslice-map":
				from, ok := tt.args.from.([]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := to[i]
					for k, v := range toelem {
						t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, from[i][k] == v)
						if from[i][k] != v {
							t.Fail()
						}
					}
				}
			case "copyslice-ptrmap":
				from, ok := tt.args.from.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := *to[i]
					fromelem := *from[i]
					for k, v := range toelem {
						t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, fromelem[k] == v)
						if fromelem[k] != v {
							t.Fail()
						}
					}
				}
			case "copyslice-from-ptrmap":
				from, ok := tt.args.from.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := to[i]
					fromelem := *from[i]
					for k, v := range toelem {
						t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, fromelem[k] == v)
						if fromelem[k] != v {
							t.Fail()
						}
					}
				}
			case "copyslice-to-ptrmap":
				from, ok := tt.args.from.([]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := *to[i]
					fromelem := from[i]
					for k, v := range toelem {
						t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, fromelem[k] == v)
						if fromelem[k] != v {
							t.Fail()
						}
					}
				}
			case "copyslice-struct":
				from, ok := tt.args.from.([]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					elem := to[i]
					t.Logf("to.Action==from.Action: %v\n", elem.Action == from[i].Action)
					t.Logf("to.Label==from.Label: %v\n", elem.Label == from[i].Label)
					if elem.Action != from[i].Action {
						t.Fail()
					}
					if elem.Label != from[i].Label {
						t.Fail()
					}
				}
			case "copyslice-ptrStruct":
				from, ok := tt.args.from.([]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := *to[i]
					fromelem := *from[i]
					t.Logf("to.Action==from.Action: %v\n", toelem.Action == fromelem.Action)
					t.Logf("to.Label==from.Label: %v\n", toelem.Label == fromelem.Label)
					if toelem.Action != fromelem.Action {
						t.Fail()
					}
					if toelem.Label != fromelem.Label {
						t.Fail()
					}
				}
			case "copyslice-from-ptrStruct":
				from, ok := tt.args.from.([]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := to[i]
					fromelem := *from[i]
					t.Logf("to.Action==from.Action: %v\n", toelem.Action == fromelem.Action)
					t.Logf("to.Label==from.Label: %v\n", toelem.Label == fromelem.Label)
					if toelem.Action != fromelem.Action {
						t.Fail()
					}
					if toelem.Label != fromelem.Label {
						t.Fail()
					}
				}
			case "copyslice-to-ptrStruct":
				from, ok := tt.args.from.([]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					toelem := *to[i]
					fromelem := from[i]
					t.Logf("to.Action==from.Action: %v\n", toelem.Action == fromelem.Action)
					t.Logf("to.Label==from.Label: %v\n", toelem.Label == fromelem.Label)
					if toelem.Action != fromelem.Action {
						t.Fail()
					}
					if toelem.Label != fromelem.Label {
						t.Fail()
					}
				}
			}
		})
	}
}

func Test_copySliceWithOption(t *testing.T) {
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
			name: "copyslice-simple",
			args: args{
				from: []int{3, 4, 5},
				to:   []int{1, 2},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copyslice-map",
			args: args{
				from: []map[string]int{{"k1": 1}},
				to:   []map[string]int{{"k2": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copyslice-ptrmap",
			args: args{
				from: []*map[string]int{{"k1": 1}},
				to:   []*map[string]int{{"k2": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copyslice-from-ptrmap",
			args: args{
				from: []*map[string]int{{"k1": 1}},
				to:   []map[string]int{{"k2": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copyslice-to-ptrmap",
			args: args{
				from: []map[string]int{{"k1": 1}},
				to:   []*map[string]int{{"k2": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copyslice-simple":
				from, ok := tt.args.from.([]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]int)
				if !ok {
					t.Fail()
				}
				fromto := append(to, from...)
				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					if to[i] != fromto[i] {
						t.Fail()
					}
				}
			case "copyslice-map":
				from, ok := tt.args.from.([]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]map[string]int)
				if !ok {
					t.Fail()
				}
				fromto := append(to, from...)
				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					fromtoElem := fromto[i]
					toElem := to[i]
					for k, v := range fromtoElem {
						tov, ok := toElem[k]
						if !ok {
							t.Fail()
						}
						if tov != v {
							t.Fail()
						}
					}
				}
			case "copyslice-ptrmap":
				from, ok := tt.args.from.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				fromto := append(to, from...)
				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					fromtoElem := *fromto[i]
					toElem := *to[i]
					for k, v := range fromtoElem {
						tov, ok := toElem[k]
						if !ok {
							t.Fail()
						}
						if tov != v {
							t.Fail()
						}
					}
				}
			case "copyslice-from-ptrmap":
				from, ok := tt.args.from.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]map[string]int)
				if !ok {
					t.Fail()
				}
				var fromto []map[string]int
				fromto = append(fromto, to...)
				for _, v := range from {
					fromto = append(fromto, *v)
				}
				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					fromtoElem := fromto[i]
					toElem := to[i]
					for k, v := range fromtoElem {
						tov, ok := toElem[k]
						if !ok {
							t.Fail()
						}
						if tov != v {
							t.Fail()
						}
					}
				}
			case "copyslice-to-ptrmap":
				from, ok := tt.args.from.([]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.([]*map[string]int)
				if !ok {
					t.Fail()
				}
				var fromto []map[string]int
				for _, v := range to {
					fromto = append(fromto, *v)
				}
				fromto = append(fromto, from...)

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for i := 0; i < len(to); i++ {
					fromtoElem := fromto[i]
					toElem := *to[i]
					for k, v := range fromtoElem {
						tov, ok := toElem[k]
						if !ok {
							t.Fail()
						}
						if tov != v {
							t.Fail()
						}
					}
				}
			}
		})
	}
}
