# Using ddcheck

`ddcheck` validates a governed data dictionary and, optionally, a JSON record against that dictionary.

## Install

### From the current main branch

```bash
go install github.com/jonathanselby/data-dictionary-first/cmd/ddcheck@main
```

The binary is installed into your Go binary directory (normally `$GOBIN` or `$GOPATH/bin`).

### From a tagged release

Once release `v0.2.0` is tagged:

```bash
go install github.com/jonathanselby/data-dictionary-first/cmd/ddcheck@v0.2.0
```

Check the installed version:

```bash
ddcheck -version
```

Expected output:

```text
ddcheck 0.2.0
```

## Validate a dictionary

```bash
ddcheck -dictionary examples/customer/dictionary.json
```

A valid dictionary returns:

```json
{
  "valid": true
}
```

## Validate a record

```bash
ddcheck \
  -dictionary examples/customer/dictionary.json \
  -record examples/customer/valid-record.json
```

## Validate another domain

```bash
ddcheck \
  -dictionary examples/supplier/dictionary.json \
  -record examples/supplier/valid-record.json
```

This is intentionally the same command shape. The validator is domain-neutral; the governed semantics live in the dictionary.

## Exit codes

`ddcheck` has a stable exit-code contract for scripts, CI/CD and ingestion workflows.

| Code | Meaning | Typical action |
|---:|---|---|
| `0` | dictionary/record is valid, or `-version` succeeded | continue |
| `1` | dictionary or record failed governed validation | reject/quarantine/fail quality gate |
| `2` | CLI usage error | correct invocation |
| `3` | file I/O or malformed JSON/runtime input error | fix file/path/input before retry |

Example shell gate:

```bash
if ddcheck -dictionary dictionary.json -record record.json; then
  echo "record accepted"
else
  case $? in
    1) echo "record rejected by governed contract" ;;
    2) echo "invalid ddcheck invocation" ;;
    3) echo "record could not be evaluated" ;;
  esac
fi
```

## Dictionary JSON Schema

The formal machine-readable dictionary contract is:

```text
schema/dictionary.schema.json
```

It uses JSON Schema Draft 2020-12 and can be consumed by editors, generators, documentation tools and independent validators.

`ddcheck` remains the reference runtime validator for this repository. The JSON Schema exists so other tooling can understand the contract without importing Go code.

## Suggested adoption pattern

For a real engineering team:

1. copy one of the example dictionaries;
2. define the dataset grain first;
3. assign data owner, steward and source of truth;
4. define each material field in plain English;
5. encode deterministic constraints (`type`, `required`, `nullable`, `allowed_values`);
6. commit the dictionary with the consuming code or data product;
7. run `ddcheck` in CI and/or ingestion validation;
8. treat dictionary changes as contract changes;
9. escalate to graph/ontology semantics only when relationship/inference requirements justify it.

## What ddcheck does not do

It is not a metadata catalogue, lineage platform, schema registry, ontology engine or data-quality platform.

It is deliberately smaller: a governed semantic contract plus deterministic validation that can be embedded into normal engineering workflows.
