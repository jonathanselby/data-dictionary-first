# Standards context

This repository is an opinionated engineering method, not an implementation of these standards. The standards are useful for locating the method on the semantic spectrum.

## ISO/IEC 11179 — metadata registries

ISO/IEC 11179-1:2023 provides the framework for metadata registries. ISO describes metadata in this family as descriptions of data.

Reference:
https://www.iso.org/standard/78914.html

This supports the basic idea that precise descriptions of data elements are a legitimate semantic layer in their own right.

## RDF — graph representation

W3C RDF represents information as a graph of subject-predicate-object statements. That is a different problem shape from validating one record against field definitions.

Reference:
https://www.w3.org/TR/rdf-concepts/

## OWL — ontology and inference

OWL 2 is a formally defined ontology language. It represents classes, properties, individuals and data values and supports software reasoning over those semantics.

References:
https://www.w3.org/TR/owl-primer/
https://www.w3.org/OWL/

The dictionary-first argument does not reject this capability. It argues that formal ontology machinery should be introduced when those capabilities are required, not before the underlying data elements have stable meaning.

## JSON Schema — shape is not the whole dictionary

JSON Schema is useful for structural validation of JSON instances. A governed data dictionary overlaps with schema concerns but also needs organisational semantics such as definition, ownership, stewardship, authority and grain.

This repository therefore keeps its dictionary model explicit rather than presenting structural schema as a complete semantic solution.
