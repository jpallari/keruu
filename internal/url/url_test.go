package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const urlYAML = `
- http://example.org/
- https://lepovirta.org/posts/index.html
- https://duckduckgo.com/?q=helloworld
`

var urlList = []string{
	"http://example.org/",
	"https://lepovirta.org/posts/index.html",
	"https://duckduckgo.com/?q=helloworld",
}

func TestURLUnmarshallingToYAML(t *testing.T) {
	// setup expected URLs
	expectedURLs := make([]URL, 0, len(urlList))
	for _, urlStr := range urlList {
		res, err := Parse(urlStr)
		if err != nil {
			t.Fatalf("failed to parse URL: %s", err)
		}
		expectedURLs = append(expectedURLs, res)
	}

	// unmarshal YAML
	var actualURLs []URL
	if err := yaml.Unmarshal([]byte(urlYAML), &actualURLs); err != nil {
		t.Fatalf("failed to unmarshal URLs: %s", err)
	}

	// Test
	assert.Equal(t, expectedURLs, actualURLs)
}
