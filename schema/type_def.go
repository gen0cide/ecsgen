package schema

// TypeDefs is an alias to a slice of TypeDef. Basically, every YAML file in the ECS spec contains
// a YAML representation of a []TypeDef where it's length is 1. This alias is simply convenience for
// unmarshaling the YAML.
type TypeDefs []TypeDef

// TypeDef holds the majority of information about a given "top level" ECS type. ECS, to it's fault,
// blends implicit and explicit typing - User or Host are explicitly typed within the YAML definitions,
// but types like "tls.server.hash" are simply implied. The loader attempts to resolve these issues and
// it starts with this top level Type definition (and it's subsequent Fields ([]*FieldDef)).
type TypeDef struct {
	Name        string       `config:"name" json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Root        bool         `config:"root" json:"root,omitempty" yaml:"root,omitempty" mapstructure:"root,omitempty"`
	Title       string       `config:"title" json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title,omitempty"`
	Group       int          `config:"group" json:"group,omitempty" yaml:"group,omitempty" mapstructure:"group,omitempty"`
	Description string       `config:"description" json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	Type        string       `config:"type" json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	Reusable    *ReusableDef `config:"reusable" json:"reusable,omitempty" yaml:"reusable,omitempty" mapstructure:"reusable,omitempty"`

	Fields []*FieldDef `config:"fields" json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`
}
