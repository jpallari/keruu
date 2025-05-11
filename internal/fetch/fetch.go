package fetch

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"go.lepovirta.org/keruu/internal/feed"
	"github.com/mmcdole/gofeed"
)

type Config struct {
	HTTPTimeout     time.Duration `yaml:"httpTimeout,omitempty"`
	PropagateErrors bool          `yaml:"propagateErrors"`
}

func (c *Config) Init() {
	c.HTTPTimeout = time.Second * 10
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP timeout can't be zero")
	}
	return nil
}

type state struct {
	config      *Config
	feeds       []feed.Config
	linkers     []feed.Linker
	feedParser  *gofeed.Parser
	errorsFound int32
}

func (s *state) init(config *Config, feeds []feed.Config, linkers []feed.Linker) {
	s.config = config
	s.feeds = feeds
	s.linkers = linkers
	s.feedParser = gofeed.NewParser()

	httpTranport := http.DefaultTransport.(*http.Transport).Clone()
	httpTranport.MaxIdleConns = 100
	httpTranport.MaxIdleConnsPerHost = 100
	httpTranport.MaxConnsPerHost = 100

	s.feedParser.Client = &http.Client{
		Transport: httpTranport,
		Timeout:   config.HTTPTimeout,
	}
}

func (s *state) run() ([]feed.Post, error) {
	var wg sync.WaitGroup
	allPosts := make([][]feed.Post, len(s.feeds))

	// Fetch feeds
	for i, f := range s.feeds {
		wg.Add(1)
		f := f
		go s.fetchFeed(&f, &wg, &allPosts[i])
	}

	// Wait for everything to finish
	wg.Wait()

	// Check for errors
	if s.errorsFound > 0 {
		return nil, fmt.Errorf("%d feed parsing errors found", s.errorsFound)
	}

	collectedPostsCount := 0
	for _, posts := range allPosts {
		collectedPostsCount += len(posts)
	}

	collectedPosts := make([]feed.Post, 0, collectedPostsCount)
	for _, posts := range allPosts {
		collectedPosts = append(collectedPosts, posts...)
	}

	return collectedPosts, nil
}

func (s *state) fetchFeed(
	f *feed.Config,
	wg *sync.WaitGroup,
	posts *[]feed.Post,
) {
	defer wg.Done()

	parsedFeed, err := s.feedParser.ParseURL(f.URL.String())
	if err != nil {
		log.Printf("error processing feed '%s': %s", &f.URL, err)
		if s.config.PropagateErrors {
			atomic.AddInt32(&s.errorsFound, 1)
		}
		return
	}
	if parsedFeed == nil {
		log.Printf("empty feed '%s'", &f.URL)
		return
	}

	*posts = make([]feed.Post, len(parsedFeed.Items))
	for i, item := range parsedFeed.Items {
		if f.IsIncluded(item.Title) {
			(*posts)[i].FromGoFeedItem(f.Name, s.linkers, parsedFeed, f.URL, item)
		}
	}
}

func Run(config *Config, feeds []feed.Config, linkers []feed.Linker) ([]feed.Post, error) {
	var state state
	state.init(config, feeds, linkers)
	collectedPosts, err := state.run()
	return collectedPosts, err
}
