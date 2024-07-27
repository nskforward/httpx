package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	/*
		Dir is the path where the cache file wiil be sored on the disc.
		Default is "cache" dir near the executable file.
	*/
	Dir string

	/*
		TotalMaxSize limits the max size of cache dir with the files.
		Default 0 - unlimited size.
	*/
	TotalMaxSize Size
}

func ValidateSettings(settings Settings) (Settings, error) {

	// Cache
	if settings.Dir != "" {
		fi, err := os.Stat(settings.Dir)
		if err != nil {
			return settings, fmt.Errorf("wrong cache dir: %w", err)
		}
		if !fi.IsDir() {
			return settings, fmt.Errorf("cache dir is a file: %s", settings.Dir)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return settings, err
		}
		settings.Dir = filepath.Join(dir, "cache")
		err = os.Mkdir(settings.Dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return settings, err
		}
	}
	return settings, nil
}
