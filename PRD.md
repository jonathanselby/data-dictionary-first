# Data Dictionary First — Product Requirements Document

## Problem

Data dictionaries are often passive documentation. They describe fields but do not reliably govern runtime data, while enterprise semantic programmes can introduce richer modelling before basic field meaning and authority are stable.

## Product objective

Provide a small reference implementation showing how a governed data dictionary can carry both human-readable semantics and deterministic runtime constraints.

The product should help answer a practical architectural question:

> **What can we solve with a dictionary before we need an ontology?**

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
- worked example;
- explicit escalation criteria for ontology/graph semantics.

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
- enterprise canonical model.

## Acceptance criteria

- repository contains a clear dictionary-first thesis and explicit boundary with ontology use cases;
- dictionary defines owner, steward, source of truth and grain;
- every field requires a human definition and supported runtime type;
- duplicate fields and structurally invalid dictionaries fail validation;
- runtime record validation checks required fields, nullability, type, allowed values and unknown-field policy;
- validation returns structured issues rather than free-form logs;
- valid worked example passes;
- invalid worked example fails for multiple material reasons;
- tests cover dictionary and runtime validation;
- implementation uses Go standard library only;
- GitHub Actions runs `go test ./...`.
