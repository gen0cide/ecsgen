package config

import (
	"fmt"
	"regexp"
)

// Whitelist is a type alias that is used to evaluate possible
type Whitelist []*regexp.Regexp

// Whitelist is used to generate a Whitelist from a given Config object.
func (c *Config) Whitelist() (Whitelist, error) {
	if len(c.whitelist.Value()) == 0 {
		return Whitelist{}, nil
	}

	ret := Whitelist{}

	for _, v := range c.whitelist.Value() {
		rx, err := regexp.Compile(v)
		if err != nil {
			return ret, fmt.Errorf("error creating regexp for whitelist value \"%s\": %v", v, err)
		}
		ret = append(ret, rx)
	}

	return ret, nil
}

// Allowed is used to check an arbitrary string against the whitelist.
func (w Whitelist) Allowed(val string) bool {
	// short circuit for an empty whitelist - allow all
	if len(w) == 0 {
		return true
	}

	passed := false
	for _, matcher := range w {
		if matcher.MatchString(val) {
			passed = true
			break
		}
	}

	return passed
}
