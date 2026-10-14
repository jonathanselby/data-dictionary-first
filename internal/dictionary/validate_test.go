package dictionary

import "testing"

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
	if r := Check(goodDictionary()); !r.Valid { t.Fatalf("expected valid: %+v", r.Issues) }
}

func TestCheckDuplicateField(t *testing.T) {
	d := goodDictionary(); d.Fields = append(d.Fields, d.Fields[0])
	if r := Check(d); r.Valid { t.Fatal("expected duplicate field failure") }
}

func TestValidateRecord(t *testing.T) {
	d := goodDictionary()
	r := ValidateRecord(d, map[string]any{"customer_id":"C-1", "status":"active", "credit_limit":1000.0})
	if !r.Valid { t.Fatalf("expected valid: %+v", r.Issues) }
}

func TestMissingRequired(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id":"C-1"})
	if r.Valid { t.Fatal("expected missing required failure") }
}

func TestAllowedValues(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id":"C-1", "status":"paused"})
	if r.Valid { t.Fatal("expected allowed value failure") }
}

func TestRejectUnknown(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id":"C-1", "status":"active", "nickname":"J"})
	if r.Valid { t.Fatal("expected unknown field failure") }
}

func TestNullable(t *testing.T) {
	r := ValidateRecord(goodDictionary(), map[string]any{"customer_id":"C-1", "status":"active", "credit_limit":nil})
	if !r.Valid { t.Fatalf("expected nullable field valid: %+v", r.Issues) }
}
