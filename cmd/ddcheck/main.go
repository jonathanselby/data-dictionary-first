package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jonathanselby/data-dictionary-first/internal/dictionary"
)

func main() {
	dictPath := flag.String("dictionary", "", "path to dictionary JSON")
	recordPath := flag.String("record", "", "optional path to JSON record to validate")
	flag.Parse()
	if *dictPath == "" {
		fmt.Fprintln(os.Stderr, "-dictionary is required")
		os.Exit(2)
	}
	d, err := dictionary.Load(*dictPath)
	if err != nil { fail(err) }
	result := dictionary.Check(d)
	if *recordPath != "" && result.Valid {
		b, err := os.ReadFile(*recordPath)
		if err != nil { fail(err) }
		var record map[string]any
		if err := json.Unmarshal(b, &record); err != nil { fail(err) }
		result = dictionary.ValidateRecord(d, record)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
	if !result.Valid { os.Exit(1) }
}

func fail(err error) {
	_ = json.NewEncoder(os.Stdout).Encode(map[string]string{"error": err.Error()})
	os.Exit(1)
}
