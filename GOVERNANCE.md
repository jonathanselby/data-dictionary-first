# Governance Model

The repository uses a deliberately small governance model because the point is to make accountability concrete rather than reproduce a data-governance framework.

## Dataset accountabilities

### Data owner

The business role accountable for the meaning and acceptable use of the dataset.

### Data steward

The role responsible for maintaining the dictionary, resolving definition disputes and coordinating changes.

### Source of truth

The authoritative system or process for the record. This is explicit because two systems can contain the same field while only one is authoritative.

### Grain

A plain statement of what one record represents. Grain is mandatory because otherwise identical field names can be aggregated or interpreted at incompatible levels.

## Field governance

Each field must have a definition. Type information alone is not accepted as meaning.

Where semantics are finite and deterministic, encode them as constraints. Examples:

```text
required
nullable
allowed_values
```

Where semantics are contextual or explanatory, retain natural language. Examples:

```text
definition
source
example
```

## Change rule

A dictionary change should be treated like a contract change.

Before changing it, ask:

1. Is the meaning changing or only the representation?
2. Is the change backwards compatible for producers and consumers?
3. Does the authoritative source change?
4. Does the grain change?
5. Do runtime validators need to change at the same time?
6. Does this expose a relationship problem that belongs outside the dictionary?

## Unknown fields

This reference implementation makes the policy explicit:

- `reject` — undeclared fields fail validation;
- `allow` — undeclared fields are tolerated.

Neither policy is universally correct. The important point is that the contract says which behaviour applies.

## Runtime consequence

A governed rule should lead to an explicit runtime choice:

```text
invalid data
    |
    +--> reject
    +--> quarantine
    +--> warn
    +--> transform under an explicit rule
```

`ddcheck` chooses reject because it is a reference validator, not a production ingestion platform.
