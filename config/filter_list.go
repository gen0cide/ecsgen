package config

import (
	"fmt"
	"regexp"
)

// FilterList is a type alias that is used to evaluate possible
type FilterList []*regexp.Regexp

// Whitelist is used to generate a Whitelist from a given Config object.
func (c *Config) Whitelist() (FilterList, error) {
	if len(c.whitelist.Value()) == 0 {
		return FilterList{}, nil
	}

	ret := FilterList{}

	for _, v := range c.whitelist.Value() {
		rx, err := regexp.Compile(v)
		if err != nil {
			return ret, fmt.Errorf("error creating regexp for whitelist value \"%s\": %v", v, err)
		}
		ret = append(ret, rx)
	}

	return ret, nil
}

// Blacklist is used to generate a FilterList of ECS keys that should not be allowed inside the model
// and is populated from a given Config object.
func (c *Config) Blacklist() (FilterList, error) {
	if len(c.blacklist.Value()) == 0 {
		return FilterList{}, nil
	}

	ret := FilterList{}

	for _, v := range c.blacklist.Value() {
		rx, err := regexp.Compile(v)
		if err != nil {
			return ret, fmt.Errorf("error creating regexp for blacklist value \"%s\": %v", v, err)
		}
		ret = append(ret, rx)
	}

	return ret, nil
}

// Match is used to check an arbitrary string against a filter list.
func (w FilterList) Match(val string) bool {
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

// Empty is a convienience function that is used to check if a FilterList is populated.
func (w FilterList) Empty() bool {
	return len(w) == 0
}
