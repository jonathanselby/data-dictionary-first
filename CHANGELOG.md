# Changelog

## 0.2.0

Adoption release: move from a usable reference implementation to a tool another engineer can install, integrate and operate without author guidance.

Why this version exists:

V0.1 proved the dictionary-first thesis and executable validation model. V0.2 closes the engineer-handoff gaps: a formal machine-readable contract, stable CLI behaviour, install/version guidance, a second domain, and acceptance tests around the examples.

- add JSON Schema Draft 2020-12 contract for the dictionary format;
- add `ddcheck -version`;
- define stable CLI exit codes: `0` valid, `1` governed validation failure, `2` usage error, `3` I/O or malformed JSON/runtime input error;
- refactor CLI behaviour so exit codes are unit-tested;
- add engineer-facing `USAGE.md` with installation and CI/integration guidance;
- add supplier-master worked example to demonstrate cross-domain reuse, including string, integer and boolean semantics;
- add CI acceptance runs for customer and supplier examples;
- document tagged `go install` path for `v0.2.0` once the GitHub tag is created.

## 0.1.0

Initial release.

- define the dictionary-first thesis and ontology escalation boundary;
- add governed dataset and field model;
- add deterministic dictionary validation;
- add JSON record validation for required fields, nullability, primitive types, governed values and undeclared-field policy;
- add customer worked example with valid and invalid records;
- add governance and standards notes;
- add unit tests and GitHub Actions.
