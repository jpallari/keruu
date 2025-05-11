package feed

import (
	"time"

	"go.lepovirta.org/keruu/internal/url"
	"github.com/mmcdole/gofeed"
)

type Post struct {
	FeedName string
	FeedLink string
	Title    string
	Link     string
	Time     *time.Time
	ExtLinks []ExtLink
}

func (p *Post) FromGoFeedItem(
	feedName string,
	linkers []Linker,
	parsedFeed *gofeed.Feed,
	feedUrl url.URL,
	item *gofeed.Item,
) {
	p.FeedName = feedName
	if p.FeedName == "" {
		p.FeedName = parsedFeed.Title
	}
	p.FeedName = feedName
	p.FeedLink = feedUrl.ResolveURL(parsedFeed.Link)
	p.Title = item.Title
	p.Link = feedUrl.ResolveURL(item.Link)
	p.Time = timeFromGoFeedItem(item)
	p.ExtLinks = goFeedItemToExtLinks(linkers, feedUrl, item)
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
