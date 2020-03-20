package schema

import (
	"strings"
)

// Object represents a Go type that will be defined by the schema generator.
type Object struct {
	Namespace *Namespace `json:"-" yaml:"-" mapstructure:"-"`

	// ID is the identifier of the field
	ID Identifier `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`

	// Fields are direct fields
	Fields map[Identifier]*Field `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`

	// Refs are fields that reference other types
	Refs map[Identifier]*Object `json:"refs,omitempty" yaml:"refs,omitempty" mapstructure:"refs,omitempty"`

	// the fields that can nest underneath
	Nested map[Identifier]*Object

	// Source keeps a reference to the parsed TypeDef YAML configuration.
	Source *TypeDef `json:"-" yaml:"-" mapstructure:"-"`
}

// FindOrCreateField will locate a field of a given identifier for an object,
// or create it based on the provided field definition.
func (o *Object) FindOrCreateField(fd *FieldDef) *Field {
	// normalize the field name: client.hash.sha1 => sha1
	fieldName := fd.Name
	if strings.Contains(fieldName, ".") {
		pathParts := strings.Split(fieldName, ".")
		fieldName = pathParts[len(pathParts)-1]
	}

	// convert to identifier
	fieldIdent := NewIdentifier(fieldName)

	// check to see if this field already exists for this object
	if field, found := o.Fields[fieldIdent]; found {
		// make sure we track all the actual field def references
		if _, exists := field.Sources[fd.Name]; !exists {
			field.Sources[fd.Name] = fd
		}
		return field
	}

	// create the new field
	field := &Field{
		Object:    o,
		Namespace: o.Namespace,
		ID:        fieldIdent,
		Sources:   map[string]*FieldDef{},
	}

	// link the source (there might be multiple)
	field.Sources[fd.Name] = fd

	o.Namespace.Loader.logger.Debugf("Type %s: New Field %s", o.ID.Pascal(), fieldIdent.Pascal())

	// Save the field in the objects map
	o.Fields[fieldIdent] = field

	return field
}

// FindOrCreateRef is used to either lookup a reference field of an object, or if it doesn't exist,
// create it. This is useful for constructing chained references on the fly.
func (o *Object) FindOrCreateRef(id Identifier) *Object {
	// see if the named reference already exists
	if obj, exists := o.Refs[id]; exists {
		return obj
	}

	// Ref didn't exist - attempt to resolve the global type, which if it doesn't exist,
	// will create it.
	obj := o.Namespace.FindType(id)
	o.Namespace.Loader.logger.Debugf("Type %s: New Ref: %s", o.ID.Pascal(), obj.ID.Pascal())

	// link the global type reference to this object as a field
	o.Refs[id] = obj

	// return it
	return obj
}

// FindOrCreateNested is used to resolve nested references of an object. These are very similar to
// regular references, but we separate them out so that we can manage their construction
// separately.
func (o *Object) FindOrCreateNested(id Identifier) *Object {
	if obj, exists := o.Nested[id]; exists {
		return obj
	}

	// lookup type, almost certainly has been created by this point
	obj := o.Namespace.FindType(id)
	o.Namespace.Loader.logger.Debugf("Type %s: New Nested: %s", o.ID.Pascal(), obj.ID.Pascal())

	// link the nested reference
	o.Nested[id] = obj
	return obj
}
