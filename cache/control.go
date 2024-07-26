package cache

import (
	"strconv"
	"strings"
	"time"
)

type Control struct {
	Private              bool
	Public               bool
	NoCache              bool
	MustRevalidate       bool
	ProxyRevalidate      bool
	NoStore              bool
	MustUnderstand       bool
	NoTransform          bool
	Immutable            bool
	MaxAge               time.Duration
	SMaxAge              time.Duration
	StaleWhileRevalidate time.Duration
	StaleIfError         time.Duration
}

func NewCacheControl(cacheControlHeader string) Control {
	var controlHeader Control
	items := strings.Split(cacheControlHeader, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		switch item {
		case "private":
			controlHeader.Private = true
		case "public":
			controlHeader.Public = true
		case "no-cache":
			controlHeader.NoCache = true
		case "must-revalidate":
			controlHeader.MustRevalidate = true
		case "proxy-revalidate":
			controlHeader.ProxyRevalidate = true
		case "no-store":
			controlHeader.NoStore = true
		case "must-understand":
			controlHeader.MustUnderstand = true
		case "no-transform":
			controlHeader.NoTransform = true
		case "immutable":
			controlHeader.Immutable = true
		}
		if !strings.Contains(item, "=") {
			continue
		}
		parts := strings.Split(item, "=")
		if len(parts) != 2 {
			continue
		}
		seconds, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}
		switch parts[0] {
		case "max-age":
			controlHeader.MaxAge = time.Duration(seconds) * time.Second
		case "s-maxage":
			controlHeader.SMaxAge = time.Duration(seconds) * time.Second
		case "stale-while-revalidate":
			controlHeader.StaleWhileRevalidate = time.Duration(seconds) * time.Second
		case "stale-if-error":
			controlHeader.StaleIfError = time.Duration(seconds) * time.Second
		}
	}
	return controlHeader
}
