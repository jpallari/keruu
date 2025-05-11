package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"go.lepovirta.org/keruu/internal/aggregation"
	"go.lepovirta.org/keruu/internal/feed"
	"go.lepovirta.org/keruu/internal/fetch"
	"go.lepovirta.org/keruu/internal/file"
	"gopkg.in/yaml.v3"
)

// TemplateYAML is a YAML template for the config
const TemplateYAML = `feeds:
  - name: <NAME>
    url: <URL>
    exclude:
      - <REGEX>
      ...
    include:
      - <REGEX>
      ...
  - name: <NAME>
    url: <URL>
  - ...
fetch:
  httpTimeout: <DURATION>
  propagateErrors <BOOLEAN>
aggregation:
  title: <STRING>
  description: <STRING>
  maxPosts: <NUMBER>
  css: <STRING>
  timeout: <DURATION>
links:
  - name: <NAME>
    url: <URL PATTERN>
`

// Config contains the configuration for the entire feed fetching, aggregation, and rendering process
type Config struct {
	Feeds       []feed.Config      `yaml:"feeds"`
	Fetch       fetch.Config       `yaml:"fetch,omitempty"`
	Aggregation aggregation.Config `yaml:"aggregation,omitempty"`
	Links       []feed.Linker      `yaml:"links,omitempty"`
}

// Init formats the config with default values
func (c *Config) Init() {
	c.Feeds = nil
	c.Fetch.Init()
	c.Aggregation.Init()
}

// FromYAMLFile reads the configuration from a YAML formatted file
func (c *Config) FromYAMLFile(filename string) (err error) {
	return file.WithFileReader(filename, func(r io.Reader) error {
		return c.FromYAML(r)
	})
}

// FromSTDIN reads the configuration from a STDIN in YAML format
func (c *Config) FromSTDIN() (err error) {
	return c.FromYAML(bufio.NewReader(os.Stdin))
}

// FromYAML reads the configuration in YAML format
func (c *Config) FromYAML(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

// ToYAML converts the configuration to YAML format
func (c *Config) ToYAML(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	errs := make([]error, 0, len(c.Feeds)+10)
	for _, feed := range c.Feeds {
		if err := feed.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(c.Feeds) == 0 {
		errs = append(errs, fmt.Errorf("no feeds provided"))
	}
	if err := c.Fetch.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := c.Aggregation.Validate(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
