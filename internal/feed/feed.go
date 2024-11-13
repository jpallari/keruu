package feed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"gitlab.com/lepovirta/keruu/internal/regexp"
	"gitlab.com/lepovirta/keruu/internal/url"
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

func (c *Config) PostFromGoFeedItem(
	linkers []Linker,
	parsedFeed *gofeed.Feed,
	feedUrl url.URL,
	item *gofeed.Item,
) *Post {
	feedName := c.Name
	if feedName == "" {
		feedName = parsedFeed.Title
	}

	return &Post{
		FeedName: feedName,
		FeedLink: feedUrl.ResolveURL(parsedFeed.Link),
		Title:    item.Title,
		Link:     feedUrl.ResolveURL(item.Link),
		Time:     timeFromGoFeedItem(item),
		ExtLinks: goFeedItemToExtLinks(linkers, feedUrl, item),
	}
}

type Post struct {
	FeedName string
	FeedLink string
	Title    string
	Link     string
	Time     *time.Time
	ExtLinks []ExtLink
}

func timeFromGoFeedItem(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	return nil
}

func (p *Post) FormattedTime() string {
	if p.Time == nil {
		return ""
	}
	return p.Time.Format("2006-01-02")
}

func (p *Post) After(other *Post) bool {
	// No timestamp means it's considered older
	if p.Time == nil {
		return false
	} else if other.Time == nil {
		return true
	}

	// Newest first
	return p.Time.After(*other.Time)
}
