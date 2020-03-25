package ecsgen

// Definition represents the YAML definition in the ECS generated "ecs_flat.yml" schema definition.
// nolint:maligned
type Definition struct {
	ID               string          `config:"-" json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`
	AllowedValues    []*AllowedValue `config:"allowed_values" json:"allowed_values,omitempty" yaml:"allowed_values,omitempty" mapstructure:"allowed_values,omitempty"`
	DashedName       string          `config:"dashed_name" json:"dashed_name,omitempty" yaml:"dashed_name,omitempty" mapstructure:"dashed_name,omitempty"`
	Description      string          `config:"description" json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	DocValues        bool            `config:"doc_values" json:"doc_values,omitempty" yaml:"doc_values,omitempty" mapstructure:"doc_values,omitempty"`
	Example          interface{}     `config:"example" json:"example,omitempty" yaml:"example,omitempty" mapstructure:"example,omitempty"`
	FlatName         string          `config:"flat_name" json:"flat_name,omitempty" yaml:"flat_name,omitempty" mapstructure:"flat_name,omitempty"`
	Format           string          `config:"format" json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
	IgnoreAbove      int             `config:"ignore_above" json:"ignore_above,omitempty" yaml:"ignore_above,omitempty" mapstructure:"ignore_above,omitempty"`
	Index            bool            `config:"index" json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`
	InputFormat      string          `config:"input_format" json:"input_format,omitempty" yaml:"input_format,omitempty" mapstructure:"input_format,omitempty"`
	Level            string          `config:"level" json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
	MultiFields      []*MultiField   `config:"multi_fields" json:"multi_fields,omitempty" yaml:"multi_fields,omitempty" mapstructure:"multi_fields,omitempty"`
	Name             string          `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Normalize        []string        `config:"normalize" json:"normalize,omitempty" yaml:"normalize,omitempty" mapstructure:"normalize,omitempty"`
	Order            int             `config:"order" json:"order,omitempty" yaml:"order,omitempty" mapstructure:"order,omitempty"`
	OriginalFieldset string          `config:"original_fieldset" json:"original_fieldset,omitempty" yaml:"original_fieldset,omitempty" mapstructure:"original_fieldset,omitempty"`
	OutputFormat     string          `config:"output_format" json:"output_format,omitempty" yaml:"output_format,omitempty" mapstructure:"output_format,omitempty"`
	OutputPrecision  string          `config:"output_precision" json:"output_precision,omitempty" yaml:"output_precision,omitempty" mapstructure:"output_precision,omitempty"`
	Short            string          `config:"short" json:"short,omitempty" yaml:"short,omitempty" mapstructure:"short,omitempty"`
	Type             string          `config:"type" json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
}

// MultiField defines a multi-field setting within an ECS schema Definition.
type MultiField struct {
	FlatName string `config:"flat_name" json:"flat_name,omitempty" yaml:"flat_name,omitempty" mapstructure:"flat_name,omitempty"`
	Name     string `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Norms    bool   `config:"norms" json:"norms,omitempty" yaml:"norms,omitempty" mapstructure:"norms,omitempty"`
	Type     string `config:"type" json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
}

// AllowedValue defines a field level constrait setting within an ECS schema Definition.
type AllowedValue struct {
	Name               string   `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Description        string   `config:"description" json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	ExpectedEventTypes []string `config:"expected_event_types" json:"expected_event_types,omitempty" yaml:"expected_event_types,omitempty" mapstructure:"expected_event_types,omitempty"`
}
