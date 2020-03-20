package schema

// Namespace defines a state table to construct ECS types. It holds all the data
// regarding information about the possible ECS types.
type Namespace struct {
	Loader *Loader
	Types  map[Identifier]*Object
}

// NewNamespace creates a new, empty namespace to use.
func NewNamespace(l *Loader) *Namespace {
	return &Namespace{
		Loader: l,
		Types:  map[Identifier]*Object{},
	}
}

// FindType is used to locate an object of a given type. If that type does not exist within
// the namespace's type table, it will be created and returned.
func (n *Namespace) FindType(ident Identifier) *Object {
	if obj, found := n.Types[ident]; found {
		return obj
	}

	n.Loader.logger.Debugf("Created Type %s", ident.Pascal())

	// Create the object
	obj := &Object{
		Namespace: n,
		ID:        ident,
		Fields:    map[Identifier]*Field{},
		Refs:      map[Identifier]*Object{},
		Nested:    map[Identifier]*Object{},
	}

	// link it in the namespace's type table
	n.Types[ident] = obj
	return obj
}
