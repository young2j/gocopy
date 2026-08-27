package gocopy

import (
	"testing"
	"time"
)

func Test_copyMap2Struct(t *testing.T) {
	type child struct {
		Name string
	}
	type target struct {
		Role   string
		Count  *int
		Tags   []string
		Labels []string
		Attrs  map[string]string
		Child  child
		Hidden string
	}

	count := 3
	result := target{Hidden: "keep"}
	Copy(&result, map[string]interface{}{
		"role":   "admin",
		"Count":  count,
		"Tags":   []string{"a", "b"},
		"Labels": []interface{}{"x", "y"},
		"Attrs":  map[string]interface{}{"env": "test"},
		"Child":  map[string]interface{}{"name": "nested"},
	})

	if result.Role != "admin" || result.Count == nil || *result.Count != count {
		t.Fatalf("basic fields not copied: %#v", result)
	}
	if len(result.Tags) != 2 || result.Tags[1] != "b" || len(result.Labels) != 2 || result.Labels[1] != "y" || result.Attrs["env"] != "test" || result.Child.Name != "nested" {
		t.Fatalf("nested fields not copied: %#v", result)
	}
	if result.Hidden != "keep" {
		t.Fatalf("missing map key should preserve destination: %#v", result)
	}
}

func Test_copyMap2StructWithOption(t *testing.T) {
	type target struct {
		To        string
		CreatedAt time.Time
		UpdatedAt string
	}
	result := target{}
	CopyWithOption(&result, map[string]interface{}{
		"from":      "mapped",
		"createdAt": "2024/01/02",
		"UpdatedAt": time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}, &Option{
		NameFromTo:   map[string]string{"from": "To"},
		StringToTime: map[string]map[string]string{"createdAt": {"loc": "UTC", "layout": "2006/01/02"}},
		TimeToString: map[string]map[string]string{"UpdatedAt": {"loc": "UTC", "layout": "2006-01-02"}},
	})

	if result.To != "mapped" || result.CreatedAt.Format("2006/01/02") != "2024/01/02" || result.UpdatedAt != "2024-01-02" {
		t.Fatalf("options not applied: %#v", result)
	}
}
