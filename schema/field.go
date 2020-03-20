package schema

// Field represents a finite field within an ECS type
type Field struct {
	Object    *Object    `json:"-" yaml:"-" mapstructure:"-"`
	Namespace *Namespace `json:"-" yaml:"-" mapstructure:"-"`
	ID        Identifier `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`
	Type      string     `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	Sources   map[string]*FieldDef
}
