# Data Dictionary First

**Define the data before you model the universe.**

`data-dictionary-first` is a practical argument for establishing governed data definitions and executable runtime semantics before committing to an enterprise ontology.

The thesis is deliberately narrower than "ontologies are bad":

> **Start with the smallest semantic structure that can make the data unambiguous and operable. Escalate to an ontology when relationships, inference or graph semantics are genuinely part of the problem.**

V0.2 moves the repository from a reference implementation to an engineer-usable utility: the dictionary format now has a formal JSON Schema, the CLI has a stable exit-code contract and version command, there are examples from two domains, and CI exercises the examples as acceptance tests.

A useful data dictionary is not a spreadsheet of column names and types. It should establish, at minimum:

- what one record represents;
- who owns and stewards the data;
- which system is authoritative;
- what each field means;
- type, required/nullable behaviour and governed values;
- what happens when runtime data violates those rules.

This repository turns that proposition into a small Go tool rather than leaving it as architecture prose.

## Why dictionary first

Many organisations jump from inconsistent data directly to an enterprise semantic model. That can reverse the dependency.

Before modelling relationships between concepts, you need to know whether the underlying concepts are stable enough to relate.

```text
RAW DATA
   |
   v
DEFINE THE RECORD
   |
   v
DEFINE EACH FIELD
   |
   v
ASSIGN OWNER + SOURCE OF TRUTH
   |
   v
MAKE CONSTRAINTS EXECUTABLE
   |
   v
RUN IT IN REAL SYSTEMS
   |
   +-- semantics sufficient? ------> KEEP IT SMALL
   |
   +-- relationships/inference
       now material? --------------> ESCALATE TO GRAPH / ONTOLOGY
```

The order matters because an ontology can formalise ambiguity just as effectively as it can formalise clarity.

## What this repository demonstrates

The included `ddcheck` utility validates a governed dictionary and then validates JSON records against it.

The dictionary carries both **human semantics** and **machine-checkable behaviour**:

```json
{
  "name": "status",
  "definition": "Current lifecycle status of the customer relationship.",
  "type": "string",
  "required": true,
  "nullable": false,
  "allowed_values": ["prospect", "active", "inactive", "closed"],
  "source": "CRM.status",
  "example": "active"
}
```

That is already enough to answer several operational questions without an ontology:

```text
What does status mean?           -> definition
Who owns the record?             -> data_owner
Where is truth maintained?       -> source_of_truth / source
Can status be null?              -> nullable
Must status exist?               -> required
What values are legitimate?      -> allowed_values
What happens to undeclared data? -> unknown_fields policy
```

## Install

Requires Go 1.23+.

Install the current main branch:

```bash
go install github.com/jonathanselby/data-dictionary-first/cmd/ddcheck@main
```

Once the `v0.2.0` GitHub tag is created, install the pinned release with:

```bash
go install github.com/jonathanselby/data-dictionary-first/cmd/ddcheck@v0.2.0
```

Check the binary:

```bash
ddcheck -version
```

Expected output:

```text
ddcheck 0.2.0
```

See [`USAGE.md`](USAGE.md) for the full engineer adoption guide.

## Run it

Check the dictionary itself:

```bash
ddcheck -dictionary examples/customer/dictionary.json
```

Validate a good customer record:

```bash
ddcheck \
  -dictionary examples/customer/dictionary.json \
  -record examples/customer/valid-record.json
```

Result:

```json
{
  "valid": true
}
```

Validate a bad customer record:

```bash
ddcheck \
  -dictionary examples/customer/dictionary.json \
  -record examples/customer/invalid-record.json
```

Result:

```json
{
  "valid": false,
  "issues": [
    {
      "path": "legal_name",
      "code": "required",
      "message": "required field is missing"
    },
    {
      "path": "status",
      "code": "allowed_value",
      "message": "value is outside governed allowed_values"
    },
    {
      "path": "nickname",
      "code": "unknown_field",
      "message": "field is not defined in the dictionary"
    }
  ]
}
```

The supplier example exercises the same validator against a different business domain and adds integer and boolean fields:

```bash
ddcheck \
  -dictionary examples/supplier/dictionary.json \
  -record examples/supplier/valid-record.json
```

## Stable exit codes

`ddcheck` is designed to sit in normal engineering automation.

| Code | Meaning |
|---:|---|
| `0` | validation passed or version command succeeded |
| `1` | governed dictionary/record validation failed |
| `2` | invalid CLI usage |
| `3` | file I/O, malformed JSON or other runtime-input failure |

This distinction matters in pipelines: invalid business data is not the same failure as being unable to read or parse the input.

## Formal dictionary contract

The dictionary format is formally described in [`schema/dictionary.schema.json`](schema/dictionary.schema.json), using JSON Schema Draft 2020-12.

That allows editors, generators, documentation systems and independent validation tooling to understand the contract without importing the Go implementation.

`ddcheck` remains the reference runtime validator for the repository.

## The escalation test

A dictionary is usually enough when the main questions are about a dataset or contract:

- What does this field mean?
- Which values are valid?
- Is it required?
- Who owns it?
- What is the system of record?
- Can a producer or consumer validate it deterministically?

An ontology becomes more useful when the important questions are about **relationships and inference across concepts**, for example:

- a person can act in several roles across several organisations;
- the same real-world entity appears in multiple domain vocabularies;
- meaning depends on class hierarchies or relationship types;
- software must infer facts that are not directly stored;
- graph traversal is itself part of the business question;
- multiple teams need a formally shared conceptual model beyond one data contract.

See [`references/dictionary-vs-ontology.md`](references/dictionary-vs-ontology.md).

## Standards context

This approach is consistent with, but deliberately smaller than, formal metadata and semantic-web standards.

- ISO/IEC 11179 defines a framework for metadata registries and describes metadata as descriptions of data.
- RDF represents information as graphs of subject-predicate-object statements.
- OWL adds formally defined ontology semantics, including classes, properties, individuals and machine reasoning.

Those are not competing ideas. They solve different levels of semantic complexity. See [`references/standards.md`](references/standards.md).

## Repository structure

```text
data-dictionary-first/
├── README.md
├── USAGE.md
├── THESIS.md
├── PRD.md
├── EDD.md
├── GOVERNANCE.md
├── architecture.txt
├── VERSION
├── CHANGELOG.md
├── LICENSE
├── go.mod
├── schema/
│   └── dictionary.schema.json
├── cmd/
│   └── ddcheck/
│       ├── main.go
│       └── main_test.go
├── internal/
│   └── dictionary/
│       ├── types.go
│       ├── validate.go
│       └── validate_test.go
├── examples/
│   ├── customer/
│   │   ├── dictionary.json
│   │   ├── valid-record.json
│   │   └── invalid-record.json
│   ├── supplier/
│   │   ├── dictionary.json
│   │   ├── valid-record.json
│   │   └── invalid-record.json
│   └── when-dictionary-stops-being-enough.md
└── references/
    ├── dictionary-vs-ontology.md
    └── standards.md
```

## Design rule

The dictionary is not documentation next to the system.

**If a semantic rule matters operationally and can be expressed deterministically, make it executable.**

That is the difference between a glossary and an engineering artefact.

## Licence

MIT. See [`LICENSE`](LICENSE).
