package schema

// ReusableDef is the configuration construct that lets nested types be specified. For example,
// group.yml defines a Group type, while user.yml defines a User type. The spec asks for a Group
// field within the User type, and specifically to re-use the existing Group type definition,
// as opposed to UserGroup being a new type.
type ReusableDef struct {
	TopLevel bool     `config:"top_level" json:"top_level,omitempty" yaml:"top_level,omitempty" mapstructure:"top_level,omitempty"`
	Expected []string `config:"expected" json:"expected,omitempty" yaml:"expected,omitempty" mapstructure:"expected,omitempty"`
}
