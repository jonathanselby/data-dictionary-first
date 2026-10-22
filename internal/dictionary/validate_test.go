package dictionary

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func goodDictionary() Dictionary {
	return Dictionary{
		Name: "customer", Version: "1.0.0", DataOwner: "Customer", DataSteward: "Data Office",
		SourceOfTruth: "CRM", Grain: "one record per customer", UnknownFields: "reject",
		Fields: []Field{
			{Name: "customer_id", Definition: "Stable enterprise customer identifier", Type: "string", Required: true},
			{Name: "status", Definition: "Current customer lifecycle status", Type: "string", Required: true, AllowedValues: []string{"active", "inactive"}},
			{Name: "credit_limit", Definition: "Approved credit limit", Type: "number", Nullable: true},
		},
	}
}

func TestCheckValid(t *testing.T) {
	if r := Check(goodDictionary()); !r.Valid {
		t.Fatalf("expected valid: %+v", r.Issues)
	}
}

func TestCheckDuplicateField(t *testing.T) {
	d := goodDictionary()
	d.Fields = append(d.Fields, d.Fields[0])
	if r := Check(d); r.Valid {
		t.Fatal("expected duplicate field failure")
	}
}

func TestValidateRecord(t *testing.T) {
	d := goodDictionary()
	r := ValidateRecord(d, map[string]any{"customer_id": "C-1", "status": "active", "credit_limit": 1000.0})
	if !r.Valid {
		t.Fatalf("expected valid: %+v", r.Issues)
	}
}

func TestMissingRequired(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id": "C-1"})
	if r.Valid {
		t.Fatal("expected missing required failure")
	}
}

func TestAllowedValues(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id": "C-1", "status": "paused"})
	if r.Valid {
		t.Fatal("expected allowed value failure")
	}
}

func TestRejectUnknown(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id": "C-1", "status": "active", "nickname": "J"})
	if r.Valid {
		t.Fatal("expected unknown field failure")
	}
}

func TestUnknownFieldIssuesAreDeterministic(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{
		"customer_id": "C-1",
		"status":      "active",
		"z_extra":     true,
		"a_extra":     true,
	})
	if r.Valid || len(r.Issues) != 2 {
		t.Fatalf("expected two unknown-field issues: %+v", r.Issues)
	}
	if r.Issues[0].Path != "a_extra" || r.Issues[1].Path != "z_extra" {
		t.Fatalf("unknown fields are not sorted: %+v", r.Issues)
	}
}

func TestNullable(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id": "C-1", "status": "active", "credit_limit": nil})
	if !r.Valid {
		t.Fatalf("expected nullable field valid: %+v", r.Issues)
	}
}

func TestLoadRejectsUnknownDictionaryFields(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"top-level": `{"name":"customer","version":"1","data_owner":"o","data_steward":"s","source_of_truth":"x","grain":"g","unknown_fields":"reject","fields":[{"name":"id","definition":"id","type":"string","required":true,"nullable":false}],"mystery":true}`,
		"field":     `{"name":"customer","version":"1","data_owner":"o","data_steward":"s","source_of_truth":"x","grain":"g","unknown_fields":"reject","fields":[{"name":"id","definition":"id","type":"string","required":true,"nullable":false,"mystery":true}]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name+".json")
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(path); err == nil {
				t.Fatal("unknown dictionary property was accepted")
			}
		})
	}
}

func TestLoadRecordPreservesLargeInteger(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.json")
	if err := os.WriteFile(path, []byte(`{"customer_id":"C-1","status":"active","sequence":9007199254740993}`), 0o600); err != nil {
		t.Fatal(err)
	}
	record, err := LoadRecord(path)
	if err != nil {
		t.Fatal(err)
	}
	n, ok := record["sequence"].(json.Number)
	if !ok || n.String() != "9007199254740993" {
		t.Fatalf("large integer lost precision: %#v", record["sequence"])
	}
	d := goodDictionary()
	d.Fields = append(d.Fields, Field{Name: "sequence", Definition: "Exact sequence", Type: "integer", Required: true})
	if r := ValidateRecord(d, record); !r.Valid {
		t.Fatalf("exact large integer rejected: %+v", r.Issues)
	}
}
