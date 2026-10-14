# Data Dictionary First — Product Requirements Document

## Version

V0.2.0 — engineer adoption release.

## Problem

Data dictionaries are often passive documentation. They describe fields but do not reliably govern runtime data, while enterprise semantic programmes can introduce richer modelling before basic field meaning and authority are stable.

A second problem appears once the idea works: a reference implementation can still depend on its author to explain installation, contract shape and failure behaviour. That is not enough for engineering adoption.

## Product objective

Provide a small, engineer-usable utility showing how a governed data dictionary can carry both human-readable semantics and deterministic runtime constraints.

The product should help answer a practical architectural question:

> **What can we solve with a dictionary before we need an ontology?**

And it should be usable without author guidance:

> **Can an engineer install it, understand the contract, run it in automation and apply it to another domain from repository evidence alone?**

## Users

- enterprise architects;
- data architects and engineers;
- application engineers;
- data stewards;
- teams defining canonical events, APIs or datasets.

## In scope

- dictionary-level ownership and authority metadata;
- record grain;
- field definitions;
- primitive runtime types;
- required and nullable behaviour;
- governed string values;
- policy for undeclared fields;
- deterministic validation of the dictionary;
- deterministic validation of JSON records;
- structured validation failures;
- formal JSON Schema contract for independent tooling;
- installable Go CLI with version command;
- stable exit-code contract for CI/CD and ingestion automation;
- worked examples across more than one business domain;
- explicit escalation criteria for ontology/graph semantics;
- automated acceptance tests around the documented examples.

## Out of scope

- metadata catalogue;
- ontology engine;
- RDF/OWL processing;
- schema registry service;
- automatic semantic inference;
- lineage engine;
- master-data management;
- data-quality profiling;
- database connectivity;
- automatic schema generation;
- enterprise canonical model;
- network service or hosted control plane.

## Acceptance criteria

- repository contains a clear dictionary-first thesis and explicit boundary with ontology use cases;
- dictionary defines owner, steward, source of truth and grain;
- every field requires a human definition and supported runtime type;
- duplicate fields and structurally invalid dictionaries fail validation;
- runtime record validation checks required fields, nullability, type, allowed values and unknown-field policy;
- validation returns structured issues rather than free-form logs;
- formal JSON Schema describes the dictionary contract;
- `ddcheck -version` reports the shipped version;
- CLI distinguishes validation failure from usage and runtime-input errors through stable exit codes;
- engineer-facing installation and usage instructions exist;
- customer and supplier examples both pass when valid and fail when invalid;
- tests cover dictionary, runtime and CLI behaviour;
- implementation uses Go standard library only;
- GitHub Actions builds the CLI, runs unit tests and exercises both documented domains.
