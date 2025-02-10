package url

import (
	"net/url"
	"strings"
)

// URL wraps the standard library URL to provide extra functionality.
// Specifically, it provides YAML unmarshalling.
type URL struct {
	*url.URL
}

func QueryEscape(s string) string {
	return url.QueryEscape(s)
}

func Parse(rawURL string) (URL, error) {
	u, err := url.Parse(rawURL)
	return URL{u}, err
}

// UnmarshalYAML parses an URL from a YAML formatted string
func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = url
	return nil
}

func (u *URL) ResolveURL(link string) string {
	if strings.HasPrefix(link, "/") {
		linkUrl, _ := url.Parse(link)
		return u.ResolveReference(linkUrl).String()
	}
	return link
}

func (u *URL) String() string {
	return u.URL.String()
}
