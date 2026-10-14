# Data Dictionary First — Engineering Design Document

## 1. Purpose

The repository contains a deliberately small Go utility, `ddcheck`, that treats a data dictionary as an executable contract.

It has two responsibilities:

1. validate that the dictionary itself contains the minimum governed semantics expected by this project;
2. optionally validate a JSON record against the dictionary's deterministic runtime rules.

V0.2 adds the engineer-adoption layer around that core: a formal JSON Schema, stable CLI exit codes, version reporting, a second business-domain example, and CI acceptance tests.

It is not a metadata platform or ontology engine.

## 2. Design principles

1. **Meaning before modelling.** Establish clear field definitions and record grain before introducing richer semantic machinery.
2. **Deterministic before inferential.** Enforce rules that can be expressed exactly before solving relationship/inference problems.
3. **Governance must have runtime consequences.** Important constraints should be testable.
4. **Keep the core small.** The implementation uses only the Go standard library.
5. **Do not turn the dictionary into an ontology.** Escalate rather than accumulating graph semantics into the wrong abstraction.
6. **Engineer handoff is part of the design.** Contract shape, installation and failure behaviour must be recoverable from the repository without author explanation.

## 3. Component model

```text
+----------------------------+
| dictionary.schema.json     |
| formal external contract   |
+-------------+--------------+
              |
              | describes
              v
+-------------------+
|   dictionary.json |
+---------+---------+
          |
          v
+-------------------+
|      ddcheck      |
| CLI / file load   |
| version/exit code |
+---------+---------+
          |
          v
+-------------------+
| Dictionary.Check  |
+---------+---------+
          |
     valid dictionary?
       /          \
      no           yes
      |             |
      v             +--------------------+
FAILED ISSUES                            |
 exit 1                                 v
                              optional record.json
                                        |
                                        v
                              ValidateRecord(...)
                                        |
                                  +-----+-----+
                                  |           |
                                valid       invalid
                                  |           |
                                  v           v
                              SUCCESS     FAILED ISSUES
                               exit 0       exit 1
```

Usage/parse/file failures sit outside governed validation and return their own exit codes.

## 4. Dictionary model

Dataset-level fields:

```text
name
version
description
data_owner
data_steward
source_of_truth
grain
unknown_fields
fields[]
```

Each field contains:

```text
name
definition
type
required
nullable
allowed_values[]
source
example
```

The supported runtime types are intentionally small:

```text
string
integer
number
boolean
object
array
```

The formal representation of this contract is `schema/dictionary.schema.json` using JSON Schema Draft 2020-12.

## 5. Dictionary validation

`Check()` validates the dictionary before runtime data is considered.

It checks:

- dictionary name and version;
- owner and steward;
- source of truth;
- record grain;
- unknown-field policy (`allow` or `reject`);
- presence of fields;
- unique field names;
- field definition;
- supported type;
- current constraint that `allowed_values` applies to strings.

A bad dictionary is not used to validate records.

The JSON Schema is an integration artefact for editors, generators and independent validators; `Check()` remains the reference runtime validation behaviour for this repository.

## 6. Record validation

`ValidateRecord()` applies the validated dictionary to one decoded JSON object.

For each declared field it checks:

- required presence;
- nullability;
- runtime type;
- governed values.

If `unknown_fields` is `reject`, undeclared fields also fail validation.

The validator does not coerce values. A semantic contract should not silently turn an invalid producer into a valid one.

## 7. CLI contract

`ddcheck` supports:

```text
-dictionary <path>
-record <path>       optional
-version
```

The exit-code contract is stable and deliberately separates governed validation from execution failures:

```text
0  valid dictionary/record, or successful version command
1  governed validation failure
2  invalid CLI usage
3  file I/O, malformed JSON or runtime-input failure
```

This makes the binary suitable for normal CI/CD and data-ingestion gates: an invalid business record is distinguishable from an unevaluable record.

## 8. Failure model

Semantic failures are returned as structured issues:

```json
{
  "path": "status",
  "code": "allowed_value",
  "message": "value is outside governed allowed_values"
}
```

Execution failures are returned as structured error JSON and exit `3`.

CLI usage errors are written to stderr and exit `2`.

## 9. Examples and reuse proof

V0.2 ships two domain examples:

- `customer` — identity, lifecycle status and optional numeric credit limit;
- `supplier` — identity, governed risk tier, integer contract count and boolean strategic flag.

The same binary and dictionary model validate both. The examples are acceptance-tested in GitHub Actions so documentation and executable behaviour cannot drift silently.

## 10. State and persistence

There is no persistent runtime state.

The dictionary and record are read from files into process memory, validated, printed as JSON and discarded.

## 11. Security boundary

- file paths are user-controlled local inputs;
- the tool does not make network calls;
- no credentials are accepted;
- no data is persisted by the tool;
- input data may still be sensitive, so shell/output handling remains an operator responsibility.

## 12. Operational behaviour

- single-process CLI;
- Go standard library only;
- deterministic validation;
- structured JSON output;
- explicit version output;
- stable automation-friendly exit codes;
- no retries because there are no remote dependencies;
- installable with normal `go install` semantics;
- GitHub Actions runs unit tests, builds the binary, checks version behaviour, runs both valid examples and asserts both invalid examples fail with exit code `1`.

## 13. Deliberate limitations

V0.2 does not implement:

- nested field schemas beyond primitive `object`/`array` type checks;
- patterns or numeric ranges;
- cross-field rules;
- referential integrity;
- relationship semantics;
- lineage;
- inference;
- compatibility/version migration rules;
- remote registry/catalogue APIs.

Some of those could be sensible future dictionary features. Others are signals that the problem has moved into a different abstraction.

## 14. Architecture boundary

The key boundary is:

```text
DATA CONTRACT QUESTION
        |
        v
dictionary + deterministic validation
        |
        +-- sufficient -----------------> stop
        |
        +-- relationships / inference
             are now the problem --------> graph / ontology evaluation
```

This boundary is more important than the validator itself.

## 15. Release boundary

V0.1 proved the thesis and core validator.

V0.2 exists because a good engineering artefact should be usable without its author present. The release therefore focuses on adoption mechanics rather than adding semantic complexity.
