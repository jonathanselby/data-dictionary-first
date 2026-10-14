package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"-version"}, &out, &errOut); code != exitOK {
		t.Fatalf("expected %d, got %d", exitOK, code)
	}
	if got := strings.TrimSpace(out.String()); got != "ddcheck 0.2.0" {
		t.Fatalf("unexpected version output %q", got)
	}
}

func TestUsageExitCode(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run(nil, &out, &errOut); code != exitUsage {
		t.Fatalf("expected %d, got %d", exitUsage, code)
	}
}

func TestRuntimeExitCode(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"-dictionary", "does-not-exist.json"}, &out, &errOut); code != exitRuntime {
		t.Fatalf("expected %d, got %d", exitRuntime, code)
	}
}

func TestValidationExitCode(t *testing.T) {
	dir := t.TempDir()
	dict := filepath.Join(dir, "dictionary.json")
	content := `{
	  "name":"customer",
	  "version":"1.0.0",
	  "description":"test dictionary",
	  "data_owner":"Customer",
	  "data_steward":"Data Office",
	  "source_of_truth":"CRM",
	  "grain":"one record per customer",
	  "unknown_fields":"reject",
	  "fields":[]
	}`
	if err := os.WriteFile(dict, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"-dictionary", dict}, &out, &errOut); code != exitValidation {
		t.Fatalf("expected %d, got %d; output=%s", exitValidation, code, out.String())
	}
}
