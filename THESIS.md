# Thesis — Data Dictionary First

## Proposition

An enterprise data model should not begin with the richest semantic mechanism available. It should begin with the smallest mechanism that makes data meaning explicit, governed and executable.

For many systems, that mechanism is a good data dictionary.

The recommended sequence is:

```text
1. define the record grain
2. define the fields
3. define ownership and authority
4. define deterministic constraints
5. use the definitions in running systems
6. observe where relationships or inference become the limiting problem
7. add graph/ontology semantics only where they earn their complexity
```

## What problem this is trying to solve

Enterprise data programmes often encounter genuine ambiguity:

- the same word means different things in different systems;
- the same concept has several names;
- required fields are not actually required everywhere;
- enumerations drift;
- systems disagree about authority;
- data governance exists in documents but not at runtime.

A common response is to design a richer enterprise semantic layer.

That may ultimately be correct. But if the underlying data elements are still unstable, a richer model starts by inheriting their instability.

The first engineering question should therefore be:

> **Can we remove the ambiguity with explicit definitions, ownership and deterministic constraints before introducing a richer semantic system?**

## What a dictionary must contain

A useful dictionary is more than name + type.

### Dataset-level semantics

- name and version;
- plain-English description;
- record grain;
- data owner;
- data steward;
- authoritative source;
- policy for undeclared fields.

### Field-level semantics

- stable field name;
- unambiguous definition;
- runtime type;
- required behaviour;
- nullability;
- governed values where finite;
- source mapping;
- representative example.

The implementation in this repository deliberately keeps the schema small. The point is to make the minimum useful semantics obvious, not to build another metadata platform.

## Why executable semantics matter

A definition that cannot influence a running system is easy to ignore.

Consider:

```text
status = "Current lifecycle status of the customer relationship"
```

Useful for a reader, but operationally weak.

Now add:

```text
type            string
required        true
nullable        false
allowed_values  prospect | active | inactive | closed
```

The same definition now becomes a contract that producers, consumers, pipelines and tests can enforce.

The engineering principle is:

> **Keep natural-language meaning for what machines cannot infer; make deterministic rules executable wherever possible.**

## Dictionary first does not mean dictionary forever

The approach has an explicit stopping point.

Do not force a dictionary to solve a graph problem.

Escalate when the material problem becomes one of:

- typed relationships between entities;
- cross-domain identity and equivalence;
- class/subclass semantics;
- multi-hop questions;
- formal inference;
- reusable enterprise vocabulary spanning many contracts;
- reasoning over incomplete explicit facts.

At that point RDF/OWL or another graph/semantic approach may be the smaller solution to the real problem.

## What to avoid

### The spreadsheet graveyard

A dictionary that lives in a spreadsheet, has no owner and is never validated at runtime is reference material, not control.

### The universal model

Do not force every local system to adopt one giant canonical shape merely because the enterprise wants consistent semantics. Shared meaning does not require identical storage.

### Ontology by aspiration

Do not introduce formal semantic machinery because the organisation wants to become "data driven". Add it when specific relationship or reasoning requirements justify it.

### Schema mistaken for meaning

A JSON Schema, SQL DDL or protobuf definition can tell you shape and type. It does not automatically establish business meaning, ownership, authority or grain.

### Governance without failure behaviour

If a field violates a governed rule, decide what happens. Reject, quarantine, warn or transform. A rule without runtime consequence is only advisory.

## Completion test

A dictionary is doing useful work when an engineer unfamiliar with the source system can answer:

- what one record represents;
- what each important field means;
- who owns it;
- where truth originates;
- what values and types are valid;
- what is required versus optional;
- what happens when new or invalid data arrives;
- which semantic questions remain unsolved.

If those questions cannot be answered, adding a richer semantic layer is unlikely to fix the foundation.
