# Release guide

This repository is prepared for release `v0.2.0`.

The GitHub connector used to build the repository can write code, documentation and branches but cannot create Git tags/releases. Creating the release is therefore the only manual repository-administration step.

## Release v0.2.0

Create a GitHub release/tag:

```text
Tag: v0.2.0
Target: main
Title: v0.2.0 — Engineer adoption release
```

Suggested release notes:

```text
V0.2 turns Data Dictionary First from a reference implementation into an engineer-usable utility.

The core dictionary-first thesis is unchanged. This release closes the handoff gaps around it:

- formal JSON Schema Draft 2020-12 contract for dictionary files;
- installable ddcheck CLI with version reporting;
- stable exit codes for CI/CD and ingestion automation;
- unit-tested CLI behaviour;
- engineer-facing installation and integration guide;
- second worked domain (supplier master) alongside customer;
- CI acceptance tests that build ddcheck, check the version contract, run both valid examples and assert invalid examples fail as governed validation errors.

The semantic scope remains deliberately small. V0.2 improves adoption, not ontology complexity.
```

## Verify before tagging

GitHub Actions should be green on the target commit.

Local checks, if desired:

```bash
go test ./...
go build -o ddcheck ./cmd/ddcheck
./ddcheck -version
./ddcheck -dictionary examples/customer/dictionary.json -record examples/customer/valid-record.json
./ddcheck -dictionary examples/supplier/dictionary.json -record examples/supplier/valid-record.json
```

Expected version:

```text
ddcheck 0.2.0
```

## Install after tagging

```bash
go install github.com/jonathanselby/data-dictionary-first/cmd/ddcheck@v0.2.0
```
