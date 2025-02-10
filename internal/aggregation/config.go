package aggregation

import (
	"fmt"
	"html/template"
)

// Config contains the feed aggregation related configurations
type Config struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	MaxPosts    int    `yaml:"maxPosts,omitempty"`
	CSSString   string `yaml:"css,omitempty"`
	Grouping    string `yaml:"grouping,omitempty"`
}

func (c *Config) Init() {
	c.Title = "Keruu"
	c.Description = "Aggregation of posts"
	c.MaxPosts = 1000
	c.CSSString = defaultCSS
	c.Grouping = defaultGrouping
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.MaxPosts <= 0 {
		return fmt.Errorf("no point in limiting result size to 0")
	}
	if !isValidGrouping(c.Grouping) {
		return fmt.Errorf("invalid grouping '%s'", c.Grouping)
	}
	return nil
}

// CSS provides the CSS data in HTML template compatible format
func (c *Config) CSS() template.CSS {
	return template.CSS(c.CSSString)
}

func (c *Config) groupFunc() GroupFunc {
	return groupingStringToFunc(c.Grouping)
}
