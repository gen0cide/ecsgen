package schema

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/elastic/go-ucfg/yaml"
	"github.com/gen0cide/genolog"
	"go.uber.org/multierr"
)

// Loader is used to laod the ECS YAML definitions and hold a data structure of definitions.
type Loader struct {
	Config *Config

	Defs map[string]TypeDef

	logger genolog.Logger
}

// NewLoader creates a new loader from a defined configuration.
func NewLoader(c *Config) (*Loader, error) {
	if c == nil {
		return nil, errors.New("cannot create loader - config was nil")
	}

	if c.logger == nil {
		return nil, errors.New("cannot create loader - config had nil logger")
	}

	return &Loader{
		Config: c,
		logger: c.logger,
		Defs:   map[string]TypeDef{},
	}, nil
}

// RawFromFile is used to parse an ECS definition from YAML. The path to the
// YAML file should be provided by the function with src.
func RawFromFile(src string) (TypeDefs, error) {
	cfg, err := yaml.NewConfigWithFile(src)
	types := TypeDefs{}
	if err != nil {
		return nil, fmt.Errorf("error reading yaml definition for file %s: %v", src, err)
	}
	if err = cfg.Unpack(&types); err != nil {
		return nil, fmt.Errorf("error parsing yaml definition for file %s: %v", src, err)
	}

	return types, nil
}

// ParseDefinitions is used to enumerate the YAML files in the source directory and parse them
// into RawField types.
func (l *Loader) ParseDefinitions() error {
	locs, err := filepath.Glob(filepath.Join(l.Config.SourceDir, "*.yml"))
	if err != nil {
		return fmt.Errorf("unable to glob YAML (*.yml) files from the source directory: %v", err)
	}

	if len(locs) == 0 {
		return fmt.Errorf("no YAML definitions (.yml) were found in source directory (%s)", l.Config.SourceDir)
	}

	// we're going to process these raw YAML files concurrently (each in separate goroutine)
	// these channels are where the results (or the errors) will be sent when it's done
	results := make(chan TypeDefs, len(locs))
	errors := make(chan error, len(locs))

	// iterate over each file and parse the contents, sending the results and errors to their
	// respective channels
	for _, defpath := range locs {
		go func(loc string) {
			res, err := RawFromFile(loc)
			if err != nil {
				errors <- err
				return
			}

			results <- res
		}(defpath)
	}

	// use this to stack multiple errors that get returned
	var errs []error

	// we should get something (either an error or a result) for each file, so enumerate
	// the channels that number of times. Store the results and aggregate the errors.
	for i := 0; i < len(locs); i++ {
		select {
		case defs := <-results:
			for _, x := range defs {
				l.Defs[x.Name] = x
			}
		case err := <-errors:
			l.logger.Errorw("error parsing definition files", "error", err)
			errs = append(errs, err)
		}
	}

	// combine any errors into a multi-error and return.
	if len(errs) > 0 {
		return multierr.Combine(errs...)
	}

	return nil
}

// BuildNamespace is used to recursively parse through the definitions and create
// a logical representation of the ECS data structures.
func (l *Loader) BuildNamespace() (*Namespace, error) {
	// create an empty namespace
	n := NewNamespace(l)

	// enumerate each def, which should contain a top level object definition
	for typeName, typeDef := range l.Defs {
		typeIdent := NewIdentifier(typeName)
		obj := n.FindType(typeIdent)

		// set the source so we can back reference if we need to later on
		if obj.Source == nil {
			newTypeRef := typeDef
			obj.Source = &newTypeRef
		}

		// enumerate the fields for this type
		for _, f := range typeDef.Fields {
			// if a field is directly mapped, just create a finite Field reference and move on
			if !strings.Contains(f.Name, ".") {
				_ = obj.FindOrCreateField(f) // don't need a reference to the field
				continue
			}

			// the field name had a "." in it - we probably need to try and resolve
			// this field to a type. First thing is to split the field name by period
			pathParts := strings.Split(f.Name, ".")

			// since this field could have multiple periods (multiple referencing objects)
			// we need to walk through them one by one. Once we get to the final
			// object, we can then actually assign this field.
			var intermediate *Object
			intermediate = obj

			// child types get prefixed with the parent object identifier, so not to collide
			// namespaces. If we didn't do this, we'd end up merging fields in the folowing types:
			// tls.client.*
			// client.*
			// so this is unfortunate, but neccessary.
			parts := []string{typeName}
			// step through all but the last identifier (the field name)
			// EXAMPLE:
			// In the "tls" object, there is a field: "server.hash.sha1"
			// We need to create references from type TLS all the way down. This loop
			// enumerates each part and then creates a type identifier. The following shows
			// what the resulting Go types will be:
			// [0] "server" => TLSServer
			// [1] "server.hash" TLSServerHash
			//
			// Once we hit that last "sha1" part - we know we're at a field. So just
			// create a field off the TLSServerHash type.
			for _, intermediateID := range pathParts[:len(pathParts)-1] {
				// append the next value
				parts = append(parts, intermediateID)
				// generate the identifier
				intermediateIdent := NewIdentifier(strings.Join(parts, "_"))
				// Find or create the intermediate object type
				intermediate = intermediate.FindOrCreateRef(intermediateIdent)
			}

			// now that we've hit the end of the object chain, we can actually create the field
			_ = intermediate.FindOrCreateField(f)
			// we don't need the reference to it
		}
	}

	// now we we need to construct the nested types - this is awkward because while type defs
	// list all of their fields, even if they've got intermediate types, they *don't* list their
	// child types. Each type explicitly defines which types it can possibly be a nested type of.
	//
	// For our purposes, the only difference between a "nested" type and a "child" type is that nested
	// types are "type shared", where as child types that we resolved above are not.
	// Example:
	// tls.server.hash is typed as TLSServerHash and exists as a field inside a TLSServer type,
	// which exists as a field inside of a TLS type.
	//
	// Where as a nested type (user.group) for example has type Group as a field inside of type User.
	// Basically it's about re-use of types.
	for typeName, typeDef := range l.Defs {
		// if the type doesn't specify that it's a nested type of anything, move on
		if typeDef.Reusable == nil {
			continue
		}

		// get the nested type's Identifier
		typeIdent := NewIdentifier(typeName)

		for _, refName := range typeDef.Reusable.Expected {
			// Find the parent the nested type references
			parentIdent := NewIdentifier(refName)
			parent := n.FindType(parentIdent)
			// link the parent to the referenced
			parent.FindOrCreateNested(typeIdent)
		}
	}

	return n, nil
}
