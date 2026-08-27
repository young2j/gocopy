package gocopy

import (
	"strings"
	"sync"
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

func TestCopyWithOptionCanBeReused(t *testing.T) {
	type source struct {
		Keep int
		Skip int
	}
	type destination struct {
		Keep int
		Skip int
	}

	opt := &Option{IgnoreFields: []string{"Skip"}}
	input := source{Keep: 1, Skip: 2}
	var first, second destination
	CopyWithOption(&first, input, opt)
	CopyWithOption(&second, input, opt)

	if len(opt.IgnoreFields) != 1 || first.Keep != 1 || first.Skip != 0 || second.Keep != 1 || second.Skip != 0 {
		t.Fatalf("reused option changed behavior: option=%#v first=%#v second=%#v", opt, first, second)
	}
}

func TestCopyWithOptionConcurrentReuse(t *testing.T) {
	type source struct {
		Keep int
		Skip int
	}
	type destination struct {
		Keep int
		Skip int
	}

	opt := &Option{IgnoreFields: []string{"Skip"}}
	input := source{Keep: 1, Skip: 2}
	var wait sync.WaitGroup
	for i := 0; i < 8; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var output destination
			CopyWithOption(&output, input, opt)
			if output.Keep != 1 || output.Skip != 0 {
				t.Errorf("concurrent option reuse changed behavior: %#v", output)
			}
		}()
	}
	wait.Wait()
}

func TestMapToStructLargeInterfaceKeyMap(t *testing.T) {
	type destination struct {
		A string
		B string
		C string
		D string
		E string
	}
	input := map[interface{}]interface{}{"a": "1", "b": "2", "c": "3", "d": "4", "e": "5"}
	var output destination
	Copy(&output, input)
	if output.A != "1" || output.B != "2" || output.C != "3" || output.D != "4" || output.E != "5" {
		t.Fatalf("interface-key map was not copied: %#v", output)
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
