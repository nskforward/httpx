package cors

import (
	"fmt"
	"net/url"
	"strings"
)

type Origin struct {
	scheme   string
	host     string
	wildcard bool
}

func Parse(s string) (Origin, error) {
	u, err := url.Parse(s)
	if err != nil {
		return Origin{}, fmt.Errorf("origin has bad format: %w", err)
	}

	if u.Host == "" {
		return Origin{}, fmt.Errorf("origin has bad format: %w", err)
	}

	if u.Scheme == "" {
		return Origin{}, fmt.Errorf("origin must contain the scheme")
	}

	host := u.Host
	switch u.Scheme {
	case "http":
		host = strings.TrimSuffix(host, ":80")
	case "https":
		host = strings.TrimSuffix(host, ":443")
	}

	if host == "" {
		return Origin{}, fmt.Errorf("origin has bad format: %w", err)
	}

	wildcard := host[0] == '*'

	if wildcard && strings.Count(host, ".") < 2 {
		return Origin{}, fmt.Errorf("origins with leading wildcard must a third level domain or higher")
	}

	return Origin{
		scheme:   u.Scheme,
		host:     host,
		wildcard: wildcard,
	}, nil
}

func (origin Origin) Equal(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if origin.scheme != u.Scheme {
		return false
	}
	host := u.Host
	switch u.Scheme {
	case "http":
		host = strings.TrimSuffix(host, ":80")
	case "https":
		host = strings.TrimSuffix(host, ":443")
	}
	return host == origin.host
}
