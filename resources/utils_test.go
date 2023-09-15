package resources

import (
	"testing"
)

func TestSliceDelta(t *testing.T) {
	ref := &[]string{"b", "c", "d"}
	compare := &[]string{"a", "b", "c"}

	// expected to receive, "a" / "d" / nil
	news, olds := SliceDelta(ref, compare)
	if len(news) != 1 {
		t.Errorf("Expected 1 new record, got %d", len(news))
	}
	if len(olds) != 1 {
		t.Errorf("Expected 1 missing erecord, got %d", len(olds))
	}
}

func TestStringMapDelta(t *testing.T) {

	ref := map[string]string{
		"new":     "new",
		"changed": "new value",
	}
	comp := map[string]string{
		"changed":  "old value",
		"obsolete": "to be deleted",
	}

	n, u, d := StringMapDelta(ref, comp)

	if len(n) != 1 {
		t.Errorf("Expected 1 new record, got %d", len(n))
	}
	if _, ok := n["new"]; !ok {
		t.Errorf("Expected 1 new record with key 'new' and value 'new'")
	}
	if len(u) != 1 {
		t.Errorf("Expected 1 record to update, got %d", len(u))
	}
	if k, ok := u["changed"]; !ok || k != "new value" {
		t.Errorf("Expected 1 record to update with key 'changed' and value 'new value'")
	}
	if len(d) != 1 {
		t.Errorf("Expected 1 record to delete, got %d", len(d))
	}
	if _, ok := d["obsolete"]; !ok {
		t.Errorf("Expected 1 record to delete with key 'obsolete'")
	}
}

func TestValidateACLTopicException(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Topic missing @syntax", args{topic: "this/is/missing/the/syntax"}, true},
		{"Topic missing syntax", args{topic: "this/is/missing/the/syntax@"}, true},
		{"Topic invalid syntax", args{topic: "this/is/missing/the/syntax@sdfdsf"}, true},

		{"syntax wrong case", args{topic: " this/is/a/invalid/syntax@SMF"}, true},

		{"whitespaces at start", args{topic: " this/is/a/invalid/topic@smf"}, true},
		{"whitespaces at end", args{topic: "this/is/a/invalid/topic@smf "}, true},
		{"whitespaces in between", args{topic: "this/is/a/invalid topic@smf "}, true},

		{"emtpy level", args{topic: "this/is//invalid@smf"}, true},

		{"start with empty level", args{topic: "/this/is/invalid@smf"}, true},

		{"valid smf topic", args{topic: "this/is/a/valid/topic@smf"}, false},
		{"valid mqtt topic", args{topic: "this/is/a/valid/topic@mqtt"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateACLTopicException(tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("validateACLTopicException() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
