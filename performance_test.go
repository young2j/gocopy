package gocopy

import "testing"

type plainCopyValue struct {
	ID      int64
	Name    string
	Enabled bool
	Scores  [4]uint32
}

func BenchmarkCopyPlainStruct(b *testing.B) {
	source := plainCopyValue{ID: 42, Name: "value", Enabled: true, Scores: [4]uint32{1, 2, 3, 4}}
	destination := plainCopyValue{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}

func BenchmarkCopyMap2Struct(b *testing.B) {
	source := map[string]interface{}{
		"role":    "admin",
		"actions": []interface{}{"read", "write"},
	}
	var destination struct {
		Role    string
		Actions []string
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}

func BenchmarkCopyMap2StructInterfaceKeys(b *testing.B) {
	source := map[interface{}]interface{}{"a": "1", "b": "2", "c": "3", "d": "4", "e": "5", "f": "6", "g": "7", "h": "8"}
	var destination struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}

func BenchmarkCopyStruct2Map(b *testing.B) {
	source := AccessRolePerms1{Role: "admin", Actions: []string{"read", "write"}}
	destination := map[string]interface{}{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}

func BenchmarkCopySliceNested(b *testing.B) {
	source := []*Perm1{{Action: "GET", Label: "read"}, {Action: "POST", Label: "write"}}
	destination := []*Perm2{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}

func BenchmarkCopyMapNested(b *testing.B) {
	source := map[string]*Perm1{"read": {Action: "GET", Label: "read"}}
	destination := map[string]*Perm2{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(&destination, source)
	}
}
