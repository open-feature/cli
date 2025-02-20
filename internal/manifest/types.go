package manifest

type BooleanFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=boolean"`
	// The default value used if a flag evaluation fails.
	CodeDefault bool `json:"codeDefault,omitempty"`
}

type StringFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=string"`
	// The default value used if a flag evaluation fails.
	CodeDefault string `json:"codeDefault,omitempty"`
}

type IntegerFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=integer"`
	// The default value used if a flag evaluation fails.
	CodeDefault string `json:"codeDefault,omitempty"`
}

type FloatFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=float"`
	// The default value used if a flag evaluation fails.
	CodeDefault string `json:"codeDefault,omitempty"`
}

type ObjectFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=object"`
	// The default value used if a flag evaluation fails.
	CodeDefault any `json:"codeDefault,omitempty"`
}

type BaseFlag struct {
	// The type of feature flag (e.g., boolean, string, integer, float)
	Type string `json:"flagType,omitempty" jsonschema:"required"`
	// A concise description of this feature flag's purpose.
	Description string `json:"description,omitempty"`
}

// Feature flag manifest for the OpenFeature CLI
type Manifest struct {
	// Collection of feature flag definitions
	Flags map[string]any `json:"flags,omitempty" jsonschema:"required"`
}
