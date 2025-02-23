package manifest

type BooleanFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=boolean"`
	// The value returned from an unsuccessful flag evaluation
	DefaultValue bool `json:"defaultValue,omitempty"`
}

type StringFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=string"`
	// The value returned from an unsuccessful flag evaluation
	DefaultValue string `json:"defaultValue,omitempty"`
}

type IntegerFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=integer"`
	// The value returned from an unsuccessful flag evaluation
	DefaultValue int `json:"defaultValue,omitempty"`
}

type FloatFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=float"`
	// The value returned from an unsuccessful flag evaluation
	DefaultValue float64 `json:"defaultValue,omitempty"`
}

type ObjectFlag struct {
	BaseFlag
	Type string `json:"flagType,omitempty" jsonschema:"enum=object"`
	// The value returned from an unsuccessful flag evaluation
	DefaultValue any `json:"defaultValue,omitempty"`
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
	Flags map[string]any `json:"flags,omitempty" jsonschema:"title=Flags,required"`
}
