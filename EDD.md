# Data Dictionary First — Engineering Design Document

## 1. Purpose

The repository contains a deliberately small Go utility, `ddcheck`, that treats a data dictionary as an executable contract.

It has two responsibilities:

1. validate that the dictionary itself contains the minimum governed semantics expected by this project;
2. optionally validate a JSON record against the dictionary's deterministic runtime rules.

It is not a metadata platform or ontology engine.

## 2. Design principles

1. **Meaning before modelling.** Establish clear field definitions and record grain before introducing richer semantic machinery.
2. **Deterministic before inferential.** Enforce rules that can be expressed exactly before solving relationship/inference problems.
3. **Governance must have runtime consequences.** Important constraints should be testable.
4. **Keep the core small.** The implementation uses only the Go standard library.
5. **Do not turn the dictionary into an ontology.** Escalate rather than accumulating graph semantics into the wrong abstraction.

## 3. Component model

```text
+-------------------+
|   dictionary.json |
+---------+---------+
          |
          v
+-------------------+
|      ddcheck      |
|  CLI / file load  |
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
                                        v
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
```

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

## 6. Record validation

`ValidateRecord()` applies the validated dictionary to one decoded JSON object.

For each declared field it checks:

- required presence;
- nullability;
- runtime type;
- governed values.

If `unknown_fields` is `reject`, undeclared fields also fail validation.

The validator does not coerce values. A semantic contract should not silently turn an invalid producer into a valid one.

## 7. Failure model

All semantic failures are returned as structured issues:

```json
{
  "path": "status",
  "code": "allowed_value",
  "message": "value is outside governed allowed_values"
}
```

The CLI exits non-zero for an invalid dictionary or record.

File and JSON parsing errors are returned as explicit execution errors.

## 8. State and persistence

There is no persistent runtime state.

The dictionary and record are read from files into process memory, validated, printed as JSON and discarded.

## 9. Security boundary

- file paths are user-controlled local inputs;
- the tool does not make network calls;
- no credentials are accepted;
- no data is persisted by the tool;
- input data may still be sensitive, so shell/output handling remains an operator responsibility.

## 10. Operational behaviour

- single-process CLI;
- Go standard library only;
- deterministic validation;
- structured JSON output;
- non-zero exit on validation failure;
- no retries because there are no remote dependencies;
- GitHub Actions executes the unit tests.

## 11. Deliberate limitations

V0.1 does not implement:

- nested field schemas beyond primitive `object`/`array` type checks;
- patterns or numeric ranges;
- cross-field rules;
- referential integrity;
- relationship semantics;
- lineage;
- inference;
- compatibility/version migration rules.

Some of those could be sensible future dictionary features. Others are signals that the problem has moved into a different abstraction.

## 12. Architecture boundary

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
