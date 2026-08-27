package gocopy

import (
	"strings"
	"testing"
	"time"
)

func TestStructTagsAreSymmetric(t *testing.T) {
	type tagged struct {
		UserName string `gocopy:"user_name"`
		JSONID   int    `json:"json_id"`
		Optional string `gocopy:",omitempty" json:"optional_name"`
		Skipped  string `gocopy:"-"`
		Empty    string `gocopy:"empty,omitempty"`
	}

	var value tagged
	if err := CopyE(&value, map[interface{}]interface{}{
		"user_name": "alice",
		"json_id":   float64(7),
		"Skipped":   "ignore",
		"empty":     "",
	}); err != nil {
		t.Fatalf("map to struct returned error: %v", err)
	}
	if value.UserName != "alice" || value.JSONID != 7 || value.Skipped != "" {
		t.Fatalf("tags were not applied: %#v", value)
	}

	result := map[string]interface{}{}
	Copy(&result, value)
	if result["user_name"] != "alice" || result["json_id"] != 7 {
		t.Fatalf("tagged fields were not emitted: %#v", result)
	}
	if _, ok := result["Skipped"]; ok {
		t.Fatalf("ignored field was emitted: %#v", result)
	}
	if _, ok := result["empty"]; ok {
		t.Fatalf("omitempty field was emitted: %#v", result)
	}
	if _, ok := result["optional_name"]; ok {
		t.Fatalf("fallback tag name was not honored: %#v", result)
	}
}

func TestCopyEReportsErrors(t *testing.T) {
	var destination int
	if err := CopyE(&destination, []string{"invalid"}); err == nil {
		t.Fatal("expected unsupported conversion error")
	}

	type target struct {
		CreatedAt time.Time
	}
	err := CopyWithOptionE(&target{}, map[string]interface{}{"CreatedAt": "not-a-time"}, &Option{
		StringToTime: map[string]map[string]string{"CreatedAt": nil},
	})
	if err == nil || !strings.Contains(err.Error(), "time") {
		t.Fatalf("expected time conversion error, got %v", err)
	}
}

func TestCopyPreservesPointerCyclesAndAliases(t *testing.T) {
	type sourceNode struct {
		Value int
		Next  *sourceNode
	}
	type destinationNode struct {
		Value int
		Next  *destinationNode
	}
	type source struct {
		First  *sourceNode
		Second *sourceNode
	}
	type destination struct {
		First  *destinationNode
		Second *destinationNode
	}

	shared := &sourceNode{Value: 9}
	shared.Next = shared
	input := source{First: shared, Second: shared}
	var output destination
	Copy(&output, input)

	if output.First == nil || output.First != output.Second || output.First.Next != output.First || output.First.Value != 9 {
		t.Fatalf("pointer cycle/alias was not preserved: %#v", output)
	}
}

func TestCopyPreservesMapCycles(t *testing.T) {
	type node struct {
		Name  string
		Child *node
	}

	input := map[string]interface{}{}
	input["Name"] = "root"
	input["Child"] = input
	var output node
	Copy(&output, input)

	if output.Name != "root" || output.Child != &output || output.Child.Name != "root" {
		t.Fatalf("map cycle was not preserved: %#v", output)
	}
}
