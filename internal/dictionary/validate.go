package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"sort"
	"strings"
)

var validTypes = map[string]bool{
	"string":  true,
	"integer": true,
	"number":  true,
	"boolean": true,
	"object":  true,
	"array":   true,
}

func Load(path string) (Dictionary, error) {
	f, err := os.Open(path)
	if err != nil {
		return Dictionary{}, err
	}
	defer f.Close()

	var d Dictionary
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	dec.UseNumber()
	if err := decodeOne(dec, &d); err != nil {
		return Dictionary{}, err
	}
	return d, nil
}

func LoadRecord(path string) (map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var record map[string]any
	dec := json.NewDecoder(f)
	dec.UseNumber()
	if err := decodeOne(dec, &record); err != nil {
		return nil, err
	}
	return record, nil
}

func decodeOne(dec *json.Decoder, v any) error {
	if err := dec.Decode(v); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}

func Check(d Dictionary) ValidationResult {
	var issues []ValidationIssue
	add := func(path, code, msg string) {
		issues = append(issues, ValidationIssue{Path: path, Code: code, Message: msg})
	}
	if strings.TrimSpace(d.Name) == "" {
		add("name", "required", "dictionary name is required")
	}
	if strings.TrimSpace(d.Version) == "" {
		add("version", "required", "dictionary version is required")
	}
	if strings.TrimSpace(d.DataOwner) == "" {
		add("data_owner", "required", "data owner is required")
	}
	if strings.TrimSpace(d.DataSteward) == "" {
		add("data_steward", "required", "data steward is required")
	}
	if strings.TrimSpace(d.SourceOfTruth) == "" {
		add("source_of_truth", "required", "source of truth is required")
	}
	if strings.TrimSpace(d.Grain) == "" {
		add("grain", "required", "record grain is required")
	}
	if d.UnknownFields != "allow" && d.UnknownFields != "reject" {
		add("unknown_fields", "invalid", "must be allow or reject")
	}
	if len(d.Fields) == 0 {
		add("fields", "required", "at least one field is required")
	}

	seen := map[string]bool{}
	for i, f := range d.Fields {
		base := fmt.Sprintf("fields[%d]", i)
		name := strings.TrimSpace(f.Name)
		if name == "" {
			add(base+".name", "required", "field name is required")
			continue
		}
		if seen[name] {
			add(base+".name", "duplicate", "field name must be unique")
		}
		seen[name] = true
		if strings.TrimSpace(f.Definition) == "" {
			add(base+".definition", "required", "field definition is required")
		}
		if !validTypes[f.Type] {
			add(base+".type", "invalid", "unsupported field type")
		}
		if len(f.AllowedValues) > 0 && f.Type != "string" {
			add(base+".allowed_values", "invalid", "allowed_values currently requires type string")
		}
	}
	return ValidationResult{Valid: len(issues) == 0, Issues: issues}
}

func ValidateRecord(d Dictionary, record map[string]any) ValidationResult {
	base := Check(d)
	if !base.Valid {
		return base
	}
	var issues []ValidationIssue
	fields := map[string]Field{}
	for _, f := range d.Fields {
		fields[f.Name] = f
	}

	for _, f := range d.Fields {
		v, exists := record[f.Name]
		if !exists {
			if f.Required {
				issues = append(issues, ValidationIssue{Path: f.Name, Code: "required", Message: "required field is missing"})
			}
			continue
		}
		if v == nil {
			if !f.Nullable {
				issues = append(issues, ValidationIssue{Path: f.Name, Code: "null", Message: "field is not nullable"})
			}
			continue
		}
		if !matchesType(f.Type, v) {
			issues = append(issues, ValidationIssue{Path: f.Name, Code: "type", Message: fmt.Sprintf("expected %s", f.Type)})
			continue
		}
		if len(f.AllowedValues) > 0 {
			s := v.(string)
			ok := false
			for _, allowed := range f.AllowedValues {
				if s == allowed {
					ok = true
					break
				}
			}
			if !ok {
				issues = append(issues, ValidationIssue{Path: f.Name, Code: "allowed_value", Message: "value is outside governed allowed_values"})
			}
		}
	}
	if d.UnknownFields == "reject" {
		unknown := make([]string, 0)
		for k := range record {
			if _, ok := fields[k]; !ok {
				unknown = append(unknown, k)
			}
		}
		sort.Strings(unknown)
		for _, k := range unknown {
			issues = append(issues, ValidationIssue{Path: k, Code: "unknown_field", Message: "field is not defined in the dictionary"})
		}
	}
	return ValidationResult{Valid: len(issues) == 0, Issues: issues}
}

func matchesType(t string, v any) bool {
	switch t {
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "number":
		switch v.(type) {
		case json.Number, float64:
			return true
		default:
			return false
		}
	case "integer":
		switch n := v.(type) {
		case json.Number:
			r, ok := new(big.Rat).SetString(n.String())
			return ok && r.IsInt()
		case float64:
			return math.Trunc(n) == n
		default:
			return false
		}
	default:
		return false
	}
}
