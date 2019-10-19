package feeds2html

import (
	"net/url"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

type aggregation struct {
	Config *AggregationConfig
	Time   time.Time
	Posts  []*feedPost
}

type feedPost struct {
	Title string
	Link  *url.URL
	Time  *time.Time
}

func goFeedItemToPost(item *gofeed.Item) (post *feedPost, err error) {
	link, err := url.Parse(item.Link)
	if err != nil {
		return
	}

	post = &feedPost{
		Title: item.Title,
		Link:  link,
		Time:  timeFromGoFeedItem(item),
	}

	return
}

func timeFromGoFeedItem(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	return nil
}

func (a *aggregation) push(post *feedPost) {
	a.Posts = append(a.Posts, post)
}

func (a *aggregation) finalize() {
	sortPostsByTime(a.Posts)
	a.Posts = a.Posts[0:a.Config.MaxItems]
	a.Time = time.Now()
}

func (a *aggregation) FormattedTime() string {
	return a.Time.Format("2006-01-02 15:04:05 -0700 MST")
}

func (p *feedPost) Hostname() string {
	return p.Link.Hostname()
}

func (p *feedPost) FormattedTime() string {
	if p.Time == nil {
		return ""
	}
	return p.Time.Format("2006-01-02")
}

func (p *feedPost) after(other *feedPost) bool {
	// No timestamp means it's considered older
	if p.Time == nil {
		return false
	} else if other.Time == nil {
		return true
	}

	// Newest first
	return p.Time.After(*other.Time)
}

func sortPostsByTime(posts []*feedPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].after(posts[j])
	})
}
