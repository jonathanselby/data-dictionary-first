# Dictionary vs ontology — decision guide

The choice is not ideological. Use the smaller abstraction that answers the material questions.

| Question | Data dictionary | Ontology / semantic graph |
|---|---|---|
| What does this field mean? | Strong fit | Can represent it, but usually more machinery than needed |
| What type is valid? | Strong fit | Possible, but not the primary reason to use ontology |
| Is the field required / nullable? | Strong fit | Usually better handled as contract/schema semantics |
| Which values are permitted? | Strong fit | Possible |
| Who owns the dataset? | Strong fit | Possible |
| What is the source of truth? | Strong fit | Possible |
| Are two entities related? | Limited | Strong fit |
| Are two concepts equivalent across domains? | Limited | Strong fit |
| Do class/subclass relationships matter? | Limited | Strong fit |
| Must software infer facts not explicitly stored? | Poor fit | Strong fit |
| Are multi-hop graph questions central? | Poor fit | Strong fit |
| Do many domains need a formally shared vocabulary? | Can help | Often strong fit |

## Rule of thumb

Use a dictionary when the unit of governance is primarily a **record, event, table, API payload or dataset**.

Evaluate ontology/graph semantics when the unit of meaning becomes a **network of entities and typed relationships**.

## Warning signs that the dictionary is being stretched too far

- fields begin encoding relationship triples in strings;
- the same entity must be represented differently depending on role;
- consumers repeatedly reconstruct the same relationship graph;
- equivalence and subclass rules are duplicated across applications;
- business questions require several relationship hops;
- teams need inference rather than validation.

At that point, adding more dictionary columns can be less simple than introducing the right semantic model.
