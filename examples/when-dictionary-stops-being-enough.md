# Example — when the dictionary stops being enough

Suppose the customer dictionary is working well. `customer_id`, `legal_name`, `status` and `credit_limit` all have stable meanings and runtime contracts.

A new requirement arrives:

> Find every natural person who controls an organisation that owns another organisation which is a guarantor for an account related to this customer.

This is no longer mainly a field-definition problem.

The material semantics are relationships:

```text
PERSON
  | controls
  v
ORGANISATION A
  | owns
  v
ORGANISATION B
  | guarantees
  v
ACCOUNT
  | relates_to
  v
CUSTOMER
```

Now further requirements appear:

- `controls` can arise through several relationship types;
- indirect ownership may imply control above a threshold;
- organisations and people appear in several domain systems;
- consumers need to traverse the relationship network;
- some relationships imply additional classifications.

A dictionary should still define the identifiers and properties used in those records. But it is no longer sufficient to express the business problem cleanly.

That is the point to evaluate a graph/ontology model.

The sequence is therefore:

```text
stable data elements
      |
      v
reliable identifiers and field semantics
      |
      v
relationships become the limiting problem
      |
      v
add relationship semantics deliberately
```

Dictionary first did not eliminate the ontology. It made the ontology start from cleaner primitives.
