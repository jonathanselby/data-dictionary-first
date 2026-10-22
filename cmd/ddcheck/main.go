package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jonathanselby/data-dictionary-first/internal/dictionary"
)

const version = "0.2.0"

const (
	exitOK         = 0
	exitValidation = 1
	exitUsage      = 2
	exitRuntime    = 3
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ddcheck", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dictPath := fs.String("dictionary", "", "path to dictionary JSON")
	recordPath := fs.String("record", "", "optional path to JSON record to validate")
	showVersion := fs.Bool("version", false, "print ddcheck version")

	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *showVersion {
		fmt.Fprintf(stdout, "ddcheck %s\n", version)
		return exitOK
	}
	if *dictPath == "" {
		fmt.Fprintln(stderr, "-dictionary is required")
		fs.Usage()
		return exitUsage
	}

	d, err := dictionary.Load(*dictPath)
	if err != nil {
		writeError(stdout, err)
		return exitRuntime
	}
	result := dictionary.Check(d)
	if *recordPath != "" && result.Valid {
		record, err := dictionary.LoadRecord(*recordPath)
		if err != nil {
			writeError(stdout, err)
			return exitRuntime
		}
		result = dictionary.ValidateRecord(d, record)
	}

	writeJSON(stdout, result)
	if !result.Valid {
		return exitValidation
	}
	return exitOK
}

func writeJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func writeError(w io.Writer, err error) {
	writeJSON(w, map[string]string{"error": err.Error()})
}
