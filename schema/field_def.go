package schema

// FieldDefs holds multiple FieldDef objects (basically the fields: array in the YAML.)
type FieldDefs []FieldDef

// FieldDef holds information about the field from the ECS YAML definition.
type FieldDef struct {
	Required              bool             `config:"required" json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`
	Index                 bool             `config:"index" json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`
	Name                  string           `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Level                 string           `config:"level" json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
	Type                  string           `config:"type" json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	Short                 string           `config:"short" json:"short,omitempty" yaml:"short,omitempty" mapstructure:"short,omitempty"`
	Description           string           `config:"description" json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	Group                 int              `config:"group" json:"group,omitempty" yaml:"group,omitempty" mapstructure:"group,omitempty"`
	ScalingFactor         int              `config:"scaling_factor" json:"scaling_factor,omitempty" yaml:"scaling_factor,omitempty" mapstructure:"scaling_factor,omitempty"`
	Format                string           `config:"format" json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
	ObjectType            string           `config:"object_type" json:"object_type,omitempty" yaml:"object_type,omitempty" mapstructure:"object_type,omitempty"`
	ObjectTypeMappingType string           `config:"object_type_mapping_type" json:"object_type_mapping_type,omitempty" yaml:"object_type_mapping_type,omitempty" mapstructure:"object_type_mapping_type,omitempty"`
	Normalize             []string         `config:"normalize" json:"normalize,omitempty" yaml:"normalize,omitempty" mapstructure:"normalize,omitempty"`
	MultiFields           []*MultiFieldDef `config:"multi_fields" json:"multi_fields,omitempty" yaml:"multi_fields,omitempty" mapstructure:"multi_fields,omitempty"`
	Example               interface{}      `config:"example" json:"example,omitempty" yaml:"example,omitempty" mapstructure:"example,omitempty"`
}

// AcceptedValueDef is used to define limitations on what a value can possibly be,
// and what another fields value might need to be relative to the value of the field. (rarely used, mostly in event.yml)
type AcceptedValueDef struct {
	Name               string   `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Description        string   `config:"description" json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	ExpectedEventTypes []string `config:"expected_event_types" json:"expected_event_types,omitempty" yaml:"expected_event_types,omitempty" mapstructure:"expected_event_types,omitempty"`
}

// MultiFieldDef is used for strange things that honestly don't make much sense yet. Need an Elasticsearch
// guru to help me here. It's unclear if some fields need to be sub nesting a "text" field in order
// to show that they must be full text indexed.
type MultiFieldDef struct {
	Type string `config:"type" json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	Name string `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
}
