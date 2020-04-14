package loader

import (
	"errors"
	"fmt"

	"ecsgen"
	"ecsgen/config"
	"github.com/elastic/go-ucfg/yaml"
)

// Loader is used to load the ECS schema definitions into a valid data structure.
type Loader struct {
	config *config.Config
	root   *ecsgen.Root
}

// NewLoader is used to create a new Loader for ECS schema definition parsing.
func NewLoader(c *config.Config) (*Loader, error) {
	if c == nil {
		return nil, errors.New("cannot create loader - config was nil")
	}

	err := c.Validate()
	if err != nil {
		return nil, fmt.Errorf("cannot create loader - invalid config: %v", err)
	}

	return &Loader{
		config: c,
		root:   ecsgen.NewRoot(),
	}, nil
}

// Load attempts to load the YAML configuration into an ecsgen definition tree.
func (l *Loader) Load() error {
	var data map[string]*ecsgen.Definition

	config, err := yaml.NewConfigWithFile(l.config.SourceFile)
	if err != nil {
		return fmt.Errorf("error reading ECS YAML: %v", err)
	}

	err = config.Unpack(&data)
	if err != nil {
		return fmt.Errorf("error marshaling YAML into ecsgen definitions: %v", err)
	}

	whitelist, err := l.config.Whitelist()
	if err != nil {
		return fmt.Errorf("error creating ecs key whitelist: %v", err)
	}

	blacklist, err := l.config.Blacklist()
	if err != nil {
		return fmt.Errorf("error creating ecs key blacklist: %v", err)
	}

	// enumerate the parsed map and create the structure
	for id, def := range data {
		// if whitelist is empty, all values will pass
		// if not, only specific values will pass
		if !whitelist.Empty() && !whitelist.Match(id) {
			continue
		}

		// if the blacklist has elements *and* they match
		// the id, skip to the next field
		if !blacklist.Empty() && blacklist.Match(id) {
			continue
		}

		def.ID = id
		node := l.root.Branch(id)
		node.Definition = def
	}

	return nil
}

// Root is used to get the loader's root definition tree.
func (l *Loader) Root() *ecsgen.Root {
	return l.root
}
