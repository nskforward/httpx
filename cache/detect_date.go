package cache

import (
	"net/http"
	"time"

	"github.com/nskforward/httpx/types"
)

func DetectDateMofified(header http.Header) time.Time {
	lastModifiedString := header.Get(types.LastModified)
	if lastModifiedString != "" {
		t, err := time.Parse(http.TimeFormat, lastModifiedString)
		if err == nil {
			return t
		}
	}

	dateString := header.Get(types.Date)
	if dateString != "" {
		t, err := time.Parse(http.TimeFormat, dateString)
		if err == nil {
			return t
		}
	}

	return time.Now()
}

func DetectDateExpiration(header http.Header, control Control) time.Time {
	expiresString := header.Get(types.Expires)
	if expiresString != "" {
		t, err := time.Parse(http.TimeFormat, expiresString)
		if err == nil {
			return t
		}
	}

	age := control.MaxAge
	if control.SMaxAge > 0 {
		age = control.SMaxAge
	}

	lastModifiedString := header.Get(types.LastModified)
	if lastModifiedString != "" {
		t, err := time.Parse(http.TimeFormat, lastModifiedString)
		if err == nil {
			return t.Add(age)
		}
	}

	dateString := header.Get(types.Date)
	if dateString != "" {
		t, err := time.Parse(http.TimeFormat, dateString)
		if err == nil {
			return t.Add(age)
		}
	}

	return time.Now().Add(age)
}
