package aggregation

import (
	"io"
	"sort"
	"time"

	"go.lepovirta.org/keruu/internal/feed"
)

type PostGroup struct {
	Name  string
	Posts []feed.Post
}

type Aggregation struct {
	Config     *Config
	Time       time.Time
	PostGroups []PostGroup
}

func (a *Aggregation) Init(
	config *Config,
	posts []feed.Post,
) {
	sortPostsByTime(posts)
	posts = posts[0:config.MaxPosts]
	a.Config = config
	a.PostGroups = groupPosts(posts, config.groupFunc())
	a.Time = time.Now()
}

func sortPostsByTime(posts []feed.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].After(&posts[j])
	})
}

func (a *Aggregation) FormattedTime() string {
	return a.Time.Format("2006-01-02 15:04:05 -0700 MST")
}

func (a *Aggregation) ToHTML(w io.Writer) error {
	return renderHTML(w, a)
}
