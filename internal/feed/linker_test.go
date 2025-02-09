package feed

import (
	"net/url"
	"testing"

	urlext "github.com/jpallari/keruu/internal/url"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestExtLinkPatternToURL(t *testing.T) {
	linker1 := Linker{
		Name:       "DuckDuckGo",
		URLPattern: "https://duckduckgo.com/?q=$TITLE",
	}
	linker2 := Linker{
		Name:       "reddit",
		URLPattern: "https://old.reddit.com/submit?url=$URL",
	}
	post1 := gofeed.Item{
		Title: "Hello World",
		Link:  "https://example.org/hello-world",
	}
	post2 := gofeed.Item{
		Title: "Hello World",
		Link:  "/hello-world",
	}
	feedUrl_, err := url.Parse("https://example.org/")
	assert.Nil(t, err)
	feedUrl := urlext.URL{URL: feedUrl_}

	url1 := linker1.goFeedItemToExtLink(feedUrl, &post1)
	url2 := linker2.goFeedItemToExtLink(feedUrl, &post1)
	url3 := linker2.goFeedItemToExtLink(feedUrl, &post2)

	assert.Equal(t, "https://duckduckgo.com/?q=Hello+World", url1.Link)
	assert.Equal(t, "https://old.reddit.com/submit?url=https%3A%2F%2Fexample.org%2Fhello-world", url2.Link)
	assert.Equal(t, "https://old.reddit.com/submit?url=https%3A%2F%2Fexample.org%2Fhello-world", url3.Link)
}
