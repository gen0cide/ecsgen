package main

import (
	"fmt"
	"sort"

	"github.com/gen0cide/ecsgen/schema"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	genConfig = schema.NewEmptyConfig()

	generateCommand = &cli.Command{
		Name:        "generate",
		Aliases:     []string{"g"},
		Usage:       "Use to translate ECS YAML definitions into a Go package.",
		Description: "Takes the input YAML definitions and translates into Idiomatic Go code.",
		Flags:       genConfig.ToCLIFlags(),
		Action:      generate,
		Before: func(c *cli.Context) error {
			genConfig.SetLogger(logger)
			return nil
		},
	}
)

type basic struct {
	Type       string   `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	Fields     []string `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`
	References []string `json:"references,omitempty" yaml:"references,omitempty" mapstructure:"references,omitempty"`
	Nested     []string `json:"nested,omitempty" yaml:"nested,omitempty" mapstructure:"nested,omitempty"`
}

func generate(c *cli.Context) error {
	logger.Info("Running Generator")
	g, err := schema.NewLoader(genConfig)
	if err != nil {
		return err
	}

	err = g.ParseDefinitions()
	if err != nil {
		return err
	}

	logger.Infof("Found %d definition files", len(g.Defs))

	ns, err := g.BuildNamespace()
	if err != nil {
		return err
	}

	_ = ns

	names := []string{}
	namemap := map[string]*schema.Object{}
	for typeIdent, obj := range ns.Types {
		names = append(names, typeIdent.Snake())
		namemap[typeIdent.Snake()] = obj
	}

	sort.Strings(names)

	results := []*basic{}

	for _, x := range names {
		typeIdent := schema.NewIdentifier(x)
		typeDef, found := namemap[x]
		if !found {
			return fmt.Errorf("could not find type def for type %s", typeIdent.Pascal())
		}

		obj := &basic{
			Type: typeIdent.Pascal(),
		}

		for name := range typeDef.Fields {
			obj.Fields = append(obj.Fields, name.Pascal())
		}

		for ref := range typeDef.Refs {
			obj.References = append(obj.References, ref.Pascal())
		}

		for nests := range typeDef.Nested {
			obj.Nested = append(obj.Nested, nests.Pascal())
		}

		results = append(results, obj)
	}

	show, err := yaml.Marshal(results)
	if err != nil {
		return err
	}

	fmt.Printf("========== RESULTS =========\n%s\n", string(show))
	return nil
}
