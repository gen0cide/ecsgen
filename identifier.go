package ecsgen

import (
	"strings"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/flect"
)

// Initialisms is a type alias for a slice of strings that is used to create capitalization
// identifiers for strings used.
type Initialisms []string

// DefaultInitialisms are the default initialisms we set. You can set more via init functions
// by calling the AddIdentifierInitialism function.
var DefaultInitialisms = Initialisms{
	"AMI",
	"PPID",
	"PID",
	"PGID",
	"MAC",
	"IP",
	"NAT",
	"IANA",
	"UID",
	"ECS",
	"AS",
	"DLL",
	"OS",
	"DNS",
	"HTTP",
	"PE",
	"JA3",
	"JA3s",
	"VLAN",
	"ISO",
	"MIME",
	"MD5",
	"SHA1",
	"SHA256",
	"SHA512",
	"UUID",
}

func init() {
	swag.AddInitialisms([]string(DefaultInitialisms)...)
}

// AddIdentifierInitialism is used to add custom initializations into the library.
func AddIdentifierInitialism(s ...string) {
	swag.AddInitialisms(s...)
}

// Identifier is used to create various cases and pluralities of the enums being generated.
type Identifier struct {
	original string
}

// NewIdentifier creates a new identifier from a provided string.
func NewIdentifier(s string) Identifier {
	return Identifier{
		original: s,
	}
}

// Snake returns the identifier's snake_case representation.
func (i Identifier) Snake() string {
	return swag.ToFileName(i.original)
}

// Pascal returns the identifiers PascalCase representation.
func (i Identifier) Pascal() string {
	return swag.ToGoName(i.original)
}

// Screaming returns the identifiers SCREAMING_CASE representation.
func (i Identifier) Screaming() string {
	return strings.ToUpper(i.Snake())
}

// Command returns the identifiers command-case representation.
func (i Identifier) Command() string {
	return swag.ToCommandName(i.original)
}

// Camel returns the identifiers camelCase representation.
func (i Identifier) Camel() string {
	return swag.ToVarName(i.original)
}

// Train returns the identifiers TRAIN-CASE representation.
func (i Identifier) Train() string {
	return strings.ToUpper(i.Command())
}

// Dotted returns the identifiers dotted.case representation.
func (i Identifier) Dotted() string {
	return strings.ReplaceAll(i.Command(), `-`, `.`)
}

// Ident returns the identifiers underlying flect structure.
func (i Identifier) Ident() flect.Ident {
	return flect.New(i.original)
}
