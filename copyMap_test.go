/*
 * File: copyMap.go
 * Created Date: 2022-01-26 06:14:44
 * Author: ysj
 * Description:  copy map to map
 */

package gocopy

import (
	"testing"
)

func Test_copyMap(t *testing.T) {
	type args struct {
		to   interface{}
		from interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copymap-simple",
			args: args{
				from: map[string]int{"key1": 1, "key2": 2},
				to:   make(map[string]int),
			},
		},
		{
			name: "copymap-slice",
			args: args{
				from: map[int][]string{1: {"a", "b", "c"}},
				to:   make(map[int][]string),
			},
		},
		{
			name: "copymap-ptrSlice",
			args: args{
				from: map[int]*[]string{1: {"a", "b", "c"}},
				to:   make(map[int]*[]string),
			},
		},
		{
			name: "copymap-from-ptrSlice",
			args: args{
				from: map[int]*[]string{1: {"a", "b", "c"}},
				to:   make(map[int][]string),
			},
		},
		{
			name: "copymap-to-ptrSlice",
			args: args{
				from: map[int][]string{1: {"a", "b", "c"}},
				to:   make(map[int]*[]string),
			},
		},
		{
			name: "copymap-map",
			args: args{
				from: map[string]map[string]int{"level1k": {"level2k": 1}},
				to:   make(map[string]map[string]int),
			},
		},
		{
			name: "copymap-ptrmap",
			args: args{
				from: map[string]*map[string]int{"level1k": {"level2k": 1}},
				to:   make(map[string]*map[string]int),
			},
		},
		{
			name: "copymap-from-ptrmap",
			args: args{
				from: map[string]*map[string]int{"level1k": {"level2k": 1}},
				to:   make(map[string]map[string]int),
			},
		},
		{
			name: "copymap-to-ptrmap",
			args: args{
				from: map[string]map[string]int{"level1k": {"level2k": 1}},
				to:   make(map[string]*map[string]int),
			},
		},
		{
			name: "copymap-struct",
			args: args{
				from: map[string]Perm1{"perm1": {Action: "POST", Label: "rest-post-method"}},
				to:   make(map[string]Perm2),
			},
		},
		{
			name: "copymap-ptrStruct",
			args: args{
				from: map[string]*Perm1{"perm1": {Action: "POST", Label: "rest-post-method"}},
				to:   make(map[string]*Perm2),
			},
		},
		{
			name: "copymap-from-ptrStruct",
			args: args{
				from: map[string]*Perm1{"perm1": {Action: "POST", Label: "rest-post-method"}},
				to:   make(map[string]Perm2),
			},
		},
		{
			name: "copymap-to-ptrStruct",
			args: args{
				from: map[string]Perm1{"perm1": {Action: "POST", Label: "rest-post-method"}},
				to:   make(map[string]*Perm2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copymap-simple":
				from, ok := tt.args.from.(map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, %v \n", k, k, v, toV == v)
					if toV != v {
						t.Fail()
					}
				}
			case "copymap-slice":
				from, ok := tt.args.from.(map[int][]string)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[int][]string)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, toV)
					for i := 0; i < len(v); i++ {
						if toV[i] != v[i] {
							t.Fail()
						}
					}
				}
			case "copymap-ptrSlice":
				from, ok := tt.args.from.(map[int]*[]string)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[int]*[]string)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, *toV)
					for i := 0; i < len(*v); i++ {
						if (*toV)[i] != (*v)[i] {
							t.Fail()
						}
					}
				}
			case "copymap-from-ptrSlice":
				from, ok := tt.args.from.(map[int]*[]string)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[int][]string)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, toV)
					for i := 0; i < len(*v); i++ {
						if toV[i] != (*v)[i] {
							t.Fail()
						}
					}
				}
			case "copymap-to-ptrSlice":
				from, ok := tt.args.from.(map[int][]string)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[int]*[]string)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, *toV)
					for i := 0; i < len(v); i++ {
						if (*toV)[i] != v[i] {
							t.Fail()
						}
					}
				}
			case "copymap-map":
				from, ok := tt.args.from.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for K, fromV := range from {
					toV, ok := to[K]
					if !ok {
						t.Fail()
					}
					if len(toV) != len(fromV) {
						t.Fail()
					}
					for k, fromv := range fromV {
						tov, ok := toV[k]
						if !ok {
							t.Fail()
						}
						if tov != fromv {
							t.Fail()
						}
					}
				}
			case "copymap-ptrmap":
				from, ok := tt.args.from.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for K, fromV := range from {
					toV, ok := to[K]
					if !ok {
						t.Fail()
					}
					if len(*toV) != len(*fromV) {
						t.Fail()
					}
					for k, fromv := range *fromV {
						tov, ok := (*toV)[k]
						if !ok {
							t.Fail()
						}
						if tov != fromv {
							t.Fail()
						}
					}
				}
			case "copymap-from-ptrmap":
				from, ok := tt.args.from.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for K, fromV := range from {
					toV, ok := to[K]
					if !ok {
						t.Fail()
					}
					if len(toV) != len(*fromV) {
						t.Fail()
					}
					for k, fromv := range *fromV {
						tov, ok := toV[k]
						if !ok {
							t.Fail()
						}
						if tov != fromv {
							t.Fail()
						}
					}
				}
			case "copymap-to-ptrmap":
				from, ok := tt.args.from.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				if len(to) != len(from) {
					t.Fail()
				}
				for K, fromV := range from {
					toV, ok := to[K]
					if !ok {
						t.Fail()
					}
					if len(*toV) != len(fromV) {
						t.Fail()
					}
					for k, fromv := range fromV {
						tov, ok := (*toV)[k]
						if !ok {
							t.Fail()
						}
						if tov != fromv {
							t.Fail()
						}
					}
				}
			case "copymap-struct":
				from, ok := tt.args.from.(map[string]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, toV)
					if toV.Action != v.Action {
						t.Fail()
					}
					if toV.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-ptrStruct":
				from, ok := tt.args.from.(map[string]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, *toV)
					if toV.Action != v.Action {
						t.Fail()
					}
					if toV.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-from-ptrStruct":
				from, ok := tt.args.from.(map[string]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, toV)
					if toV.Action != v.Action {
						t.Fail()
					}
					if toV.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-to-ptrStruct":
				from, ok := tt.args.from.(map[string]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*Perm2)
				if !ok {
					t.Fail()
				}
				Copy(&to, from)

				for k, v := range from {
					toV, ok := to[k]
					if !ok {
						t.Fail()
						break
					}
					t.Logf("to[%v]==from[%v]: %v, \n", k, k, *toV)
					if toV.Action != v.Action {
						t.Fail()
					}
					if toV.Label != v.Label {
						t.Fail()
					}
				}
			}
		})
	}
}

func Test_copyMapWithOption(t *testing.T) {
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
			name: "copymap-simple",
			args: args{
				from: map[string]int{"key1": 1, "key2": 2},
				to:   map[string]int{"key0": 0, "key2": 3},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-slice",
			args: args{
				from: map[string][]int{"key1": {3, 4}, "key2": {7, 8}},
				to:   map[string][]int{"key1": {1, 2}, "key3": {5, 6}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-ptrSlice",
			args: args{
				from: map[string]*[]int{"key1": {3, 4}, "key2": {7, 8}},
				to:   map[string]*[]int{"key1": {1, 2}, "key3": {5, 6}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-from-ptrSlice",
			args: args{
				from: map[string]*[]int{"key1": {3, 4}, "key2": {7, 8}},
				to:   map[string][]int{"key1": {1, 2}, "key3": {5, 6}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-to-ptrSlice",
			args: args{
				from: map[string][]int{"key1": {3, 4}, "key2": {7, 8}},
				to:   map[string]*[]int{"key1": {1, 2}, "key3": {5, 6}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-map",
			args: args{
				from: map[string]map[string]int{"level1k": {"level2k": 1}},
				to:   map[string]map[string]int{"level1k": {"level2k_": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-ptrmap",
			args: args{
				from: map[string]*map[string]int{"level1k": {"level2k": 1}},
				to:   map[string]*map[string]int{"level1k": {"level2k_": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-from-ptrmap",
			args: args{
				from: map[string]*map[string]int{"level1k": {"level2k": 1}},
				to:   map[string]map[string]int{"level1k": {"level2k_": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-to-ptrmap",
			args: args{
				from: map[string]map[string]int{"level1k": {"level2k": 1}},
				to:   map[string]*map[string]int{"level1k": {"level2k_": 2}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-struct",
			args: args{
				from: map[string]Perm1{"perm1": {Action: "PUT", Label: "rest-put-method"}},
				to:   map[string]Perm2{"perm2": {Action: "GET", Label: "rest-get-method"}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-ptrStruct",
			args: args{
				from: map[string]*Perm1{"perm1": {Action: "PUT", Label: "rest-put-method"}},
				to:   map[string]*Perm2{"perm2": {Action: "GET", Label: "rest-get-method"}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-from-ptrStruct",
			args: args{
				from: map[string]*Perm1{"perm1": {Action: "PUT", Label: "rest-put-method"}},
				to:   map[string]Perm2{"perm2": {Action: "GET", Label: "rest-get-method"}},
				opt: &Option{
					Append: true,
				},
			},
		},
		{
			name: "copymap-to-ptrStruct",
			args: args{
				from: map[string]Perm1{"perm1": {Action: "PUT", Label: "rest-put-method"}},
				to:   map[string]*Perm2{"perm2": {Action: "GET", Label: "rest-get-method"}},
				opt: &Option{
					Append: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "copymap-simple":
				from, ok := tt.args.from.(map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]int)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					fromto[k] = v
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range to {
					fromtoV, ok := fromto[k]
					if !ok {
						t.Fail()
					}
					if fromtoV != v {
						t.Fail()
					}
				}
			case "copymap-slice":
				from, ok := tt.args.from.(map[string][]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string][]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string][]int)
				for k, v := range to {
					fromto[k] = append(fromto[k], v...)
				}
				for k, v := range from {
					fromto[k] = append(fromto[k], v...)
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for i := 0; i < len(v); i++ {
						if tov[i] != fromto[k][i] {
							t.Fail()
						}
					}
				}
			case "copymap-ptrSlice":
				from, ok := tt.args.from.(map[string]*[]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*[]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string][]int)
				for k, v := range to {
					fromto[k] = append(fromto[k], *v...)
				}
				for k, v := range from {
					fromto[k] = append(fromto[k], *v...)
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for i := 0; i < len(v); i++ {
						if (*tov)[i] != fromto[k][i] {
							t.Fail()
						}
					}
				}
			case "copymap-from-ptrSlice":
				from, ok := tt.args.from.(map[string]*[]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string][]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string][]int)
				for k, v := range to {
					fromto[k] = append(fromto[k], v...)
				}
				for k, v := range from {
					fromto[k] = append(fromto[k], *v...)
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for i := 0; i < len(v); i++ {
						if tov[i] != fromto[k][i] {
							t.Fail()
						}
					}
				}
			case "copymap-to-ptrSlice":
				from, ok := tt.args.from.(map[string][]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*[]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string][]int)
				for k, v := range to {
					fromto[k] = append(fromto[k], *v...)
				}
				for k, v := range from {
					fromto[k] = append(fromto[k], v...)
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for i := 0; i < len(v); i++ {
						if (*tov)[i] != fromto[k][i] {
							t.Fail()
						}
					}
				}
			case "copymap-map":
				from, ok := tt.args.from.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]map[string]int)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					tov, ok := fromto[k]
					if !ok {
						fromto[k] = v
						continue
					}
					for kk, vv := range v {
						tov[kk] = vv
					}
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for kk, vv := range tov {
						fromtovv, ok := v[kk]
						if !ok {
							t.Fail()
							break
						}
						if fromtovv != vv {
							t.Fail()
						}
					}
				}
			case "copymap-ptrmap":
				from, ok := tt.args.from.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]map[string]int)
				for k, v := range to {
					fromto[k] = *v
				}
				for k, v := range from {
					tov, ok := fromto[k]
					if !ok {
						fromto[k] = *v
						continue
					}
					for kk, vv := range *v {
						tov[kk] = vv
					}
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for kk, vv := range *tov {
						fromtovv, ok := v[kk]
						if !ok {
							t.Fail()
							break
						}
						if fromtovv != vv {
							t.Fail()
						}
					}
				}
			case "copymap-from-ptrmap":
				from, ok := tt.args.from.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]map[string]int)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					tov, ok := fromto[k]
					if !ok {
						fromto[k] = *v
						continue
					}
					for kk, vv := range *v {
						tov[kk] = vv
					}
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for kk, vv := range tov {
						fromtovv, ok := v[kk]
						if !ok {
							t.Fail()
							break
						}
						if fromtovv != vv {
							t.Fail()
						}
					}
				}
			case "copymap-to-ptrmap":
				from, ok := tt.args.from.(map[string]map[string]int)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*map[string]int)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]map[string]int)
				for k, v := range to {
					fromto[k] = *v
				}
				for k, v := range from {
					tov, ok := fromto[k]
					if !ok {
						fromto[k] = v
						continue
					}
					for kk, vv := range v {
						tov[kk] = vv
					}
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					for kk, vv := range *tov {
						fromtovv, ok := v[kk]
						if !ok {
							t.Fail()
							break
						}
						if fromtovv != vv {
							t.Fail()
						}
					}
				}
			case "copymap-struct":
				from, ok := tt.args.from.(map[string]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]Perm2)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]Perm2)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					perm2 := Perm2{
						Action: v.Action,
						Label:  v.Label,
					}
					fromto[k] = perm2
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					if tov.Action != v.Action {
						t.Fail()
					}
					if tov.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-ptrStruct":
				from, ok := tt.args.from.(map[string]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*Perm2)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]*Perm2)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					perm2 := &Perm2{
						Action: v.Action,
						Label:  v.Label,
					}
					fromto[k] = perm2
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					if tov.Action != v.Action {
						t.Fail()
					}
					if tov.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-from-ptrStruct":
				from, ok := tt.args.from.(map[string]*Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]Perm2)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]Perm2)
				for k, v := range to {
					fromto[k] = v
				}
				for k, v := range from {
					perm2 := Perm2{
						Action: v.Action,
						Label:  v.Label,
					}
					fromto[k] = perm2
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					if tov.Action != v.Action {
						t.Fail()
					}
					if tov.Label != v.Label {
						t.Fail()
					}
				}
			case "copymap-to-ptrStruct":
				from, ok := tt.args.from.(map[string]Perm1)
				if !ok {
					t.Fail()
				}
				to, ok := tt.args.to.(map[string]*Perm2)
				if !ok {
					t.Fail()
				}

				fromto := make(map[string]Perm2)
				for k, v := range to {
					fromto[k] = *v
				}
				for k, v := range from {
					perm2 := Perm2{
						Action: v.Action,
						Label:  v.Label,
					}
					fromto[k] = perm2
				}

				CopyWithOption(&to, from, tt.args.opt)

				if len(to) != len(fromto) {
					t.Fail()
				}
				for k, v := range fromto {
					tov, ok := to[k]
					if !ok {
						t.Fail()
					}
					if tov.Action != v.Action {
						t.Fail()
					}
					if tov.Label != v.Label {
						t.Fail()
					}
				}
			}
		})
	}
}
