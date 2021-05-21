package generator

import (
	"fmt"
	"sync"

	"github.com/urfave/cli/v2"
)

// Registry is a type that is used to hold all the output generators. This is similar
// to common plugin registry patterns in other go code.
type Registry interface {
	// Register is used to add a Generator to the Registry.
	Register(Generator) error

	// Get is used to retrieve a Generator from the Registry.
	Get(string) (Generator, error)

	// All is used to retrieve all of the generators in the registry.
	All() []Generator

	// CLIFlags returns a flat set of CLI flags for all plugins registered.
	CLIFlags() []cli.Flag
}

// registry implements the Registry interface in a concurrent safe manner.
type registry struct {
	sync.RWMutex

	store   map[string]Generator
	ordered []Generator
}

// NewRegistry is used to initialize a new registry.
func NewRegistry() Registry {
	return &registry{
		store:   map[string]Generator{},
		ordered: []Generator{},
	}
}

// Get implements the generator.Registry interface.
func (r *registry) Get(id string) (Generator, error) {
	r.Lock()
	defer r.Unlock()

	if generator, found := r.store[id]; found {
		return generator, nil
	}

	return nil, fmt.Errorf("generator %s could not be found in the registry", id)
}

// Register implements the generator.Registry interface.
func (r *registry) Register(g Generator) error {
	r.Lock()
	defer r.Unlock()

	if _, found := r.store[g.ID()]; found {
		return fmt.Errorf("generator %s has already been registered", g.ID())
	}

	r.store[g.ID()] = g
	r.ordered = append(r.ordered, g)
	return nil
}

// All implements the generator.Registry interface.
func (r *registry) All() []Generator {
	r.Lock()
	defer r.Unlock()

	// make a copy, don't just return our slice
	ret := make([]Generator, len(r.ordered))
	copy(ret, r.ordered)

	return ret
}

// CLIFlags implements the generator.Registry interface.
func (r *registry) CLIFlags() []cli.Flag {
	// short circuit if we have no plugins registered
	if len(r.ordered) == 0 {
		return []cli.Flag{}
	}

	r.Lock()
	defer r.Unlock()

	res := []cli.Flag{}

	for _, generator := range r.ordered {
		flags := shimCLIFlags(generator)
		res = append(res, flags...)
	}

	return res
}
