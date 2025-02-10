package feed

import (
	"fmt"

	"github.com/jpallari/keruu/internal/regexp"
	"github.com/jpallari/keruu/internal/url"
)

// Config contains the details of a single feed
type Config struct {
	Name    string      `yaml:"name"`
	URL     url.URL     `yaml:"url"`
	Exclude []regexp.RE `yaml:"exclude,omitempty"`
	Include []regexp.RE `yaml:"include,omitempty"`
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("no name specified for feed")
	}
	if c.URL.URL == nil {
		return fmt.Errorf("no URL specified for feed '%s'", c.Name)
	}
	return nil
}

func (c *Config) IsIncluded(s string) bool {
	// Always exclude first
	for _, filter := range c.Exclude {
		if filter.MatchString(s) {
			return false
		}
	}

	for _, filter := range c.Include {
		if filter.MatchString(s) {
			return true
		}
	}

	// No match on filters?
	// Only include it, if there was no include filters set
	return len(c.Include) == 0
}

