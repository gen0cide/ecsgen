package schema

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/gen0cide/genolog"
	"github.com/urfave/cli"
)

var (
	// ErrInvalidSourceDir is thrown when the source directory cannot be resolved
	ErrInvalidSourceDir = errors.New("source directory was either blank or did not exist")

	// ErrNoDefinitionsInSourceDir is thrown when the source directory does not contain any valid
	// ECS definitions.
	ErrNoDefinitionsInSourceDir = errors.New("source directory does not contain any valid ecs definitions")

	// ErrInvalidPackageName is thrown when a Go package name is either not specified or is not valid.
	ErrInvalidPackageName = errors.New("package name was either empty or an invalid go package identifier")

	// ErrInvalidOutputDir is thrown when the output directory does not exist.
	ErrInvalidOutputDir = errors.New("output directory was either blank or did not exist")
)

// Config holds the parameters needed for proper generation of Go code.
type Config struct {
	SourceDir   string `json:"source_dir" yaml:"source_dir" toml:"source_dir"`
	PackageName string `json:"package_name" yaml:"package_dir" toml:"package_dir"`
	OutputDir   string `json:"output_dir" yaml:"output_dir" toml:"output_dir"`

	logger genolog.Logger
}

// NewEmptyConfig is a constructor for an empty Config object.
func NewEmptyConfig() *Config {
	return &Config{}
}

// SetLogger is used to override the logger.
func (c *Config) SetLogger(logger genolog.Logger) {
	c.logger = logger
}

// ToCLIFlags is a helper to automatically set fields within a Config object
// based on CLI flags using the github.com/urfave/cli framework.
func (c *Config) ToCLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Usage:       "Path to directory containing ECS YAML definitions.",
			EnvVars:     []string{"ECSGEN_SOURCE"},
			Required:    true,
			Destination: &c.SourceDir,
		},
		&cli.StringFlag{
			Name:        "package-name",
			Usage:       "Name of the Go package for the generated code.",
			EnvVars:     []string{"ECSGEN_PACKAGE_NAME"},
			Required:    true,
			Destination: &c.PackageName,
		},
		&cli.StringFlag{
			Name:        "output",
			Usage:       "Path to the directory where the generated code should be written.",
			EnvVars:     []string{"ECSGEN_OUTPUT"},
			Required:    true,
			Destination: &c.OutputDir,
		},
	}
}

// Validate validates that the configuration has expected parameters.
func (c *Config) Validate() error {
	// Check the Source Directory
	// is it assigned?
	if c.SourceDir == "" {
		return ErrInvalidSourceDir
	}

	// Is it a valid path?
	dir, err := os.Stat(c.SourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrInvalidSourceDir
		}

		return fmt.Errorf("error locating specified source directory: %v", err)
	}

	// is it a valid directory?
	if !dir.IsDir() {
		return fmt.Errorf("specified source directory was a path to a file, not a directory")
	}

	// Check the Output Directory
	// is it assigned?
	if c.OutputDir == "" {
		return ErrInvalidOutputDir
	}

	// Is it a valid path?
	dir, err = os.Stat(c.OutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrInvalidOutputDir
		}

		return fmt.Errorf("error locating specified output directory: %v", err)
	}

	// is it a valid directory?
	if !dir.IsDir() {
		return fmt.Errorf("specified output directory was a path to a file, not a directory")
	}

	// while Go maintains STRONG guidance on package naming conventions,
	// it doesn't actually seem to enforce a whole lot. Keeping it basic for now.
	pkgRegex := regexp.MustCompile(`^[a-zA-Z0-9\_]{1,64}$`)
	if !pkgRegex.MatchString(c.PackageName) {
		return ErrInvalidPackageName
	}

	return nil
}
