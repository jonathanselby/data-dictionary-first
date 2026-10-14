package dictionary

type Dictionary struct {
	Name          string  `json:"name"`
	Version       string  `json:"version"`
	Description   string  `json:"description"`
	DataOwner     string  `json:"data_owner"`
	DataSteward   string  `json:"data_steward"`
	SourceOfTruth string  `json:"source_of_truth"`
	Grain         string  `json:"grain"`
	UnknownFields string  `json:"unknown_fields"`
	Fields        []Field `json:"fields"`
}

type Field struct {
	Name          string   `json:"name"`
	Definition    string   `json:"definition"`
	Type          string   `json:"type"`
	Required      bool     `json:"required"`
	Nullable      bool     `json:"nullable"`
	AllowedValues []string `json:"allowed_values,omitempty"`
	Source        string   `json:"source,omitempty"`
	Example       any      `json:"example,omitempty"`
}

type ValidationIssue struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues,omitempty"`
}
