package regexp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const regexpYAML = `
- "^Sponsored Post:"
- \[ad\]$
`

var regexpList = []string{
	"^Sponsored Post:",
	"\\[ad\\]$",
}

func TestRegexpUnmarshallingToYAML(t *testing.T) {
	// setup expected regexps
	expectedRegexps := make([]RE, 0, len(regexpList))
	for _, reStr := range regexpList {
		re, err := Compile(reStr)
		if err != nil {
			t.Fatalf("failed to parse regexp: %s", err)
		}
		expectedRegexps = append(expectedRegexps, re)
	}

	// unmarshal YAML
	var actualREs []RE
	if err := yaml.Unmarshal([]byte(regexpYAML), &actualREs); err != nil {
		t.Fatalf("failed to unmarshal regexps: %s", err)
	}

	// Test
	assert.Equal(t, expectedRegexps, actualREs)
}
