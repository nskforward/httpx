package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/httpx/gzipx"
)

type Store struct {
	baseDir    string
	defaultTTL time.Duration
}

func NewStore(baseDir string, defaultTTL time.Duration) *Store {
	_, err := os.Stat(baseDir)
	if err != nil {
		panic(fmt.Errorf("cache base dir not found: %w", err))
	}
	return &Store{
		baseDir:    baseDir,
		defaultTTL: defaultTTL,
	}
}

func (s *Store) metadataPath(path string) string {
	return s.fullPath(path) + ".meta"
}

func (s *Store) archivePath(path string) string {
	return s.fullPath(path) + ".gz"
}

func (s *Store) fullPath(path string) string {
	return filepath.Join(s.baseDir, path)
}

func (s *Store) Get(path string) (*Entry, error) {
	metaPath := s.metadataPath(path)
	f, err := os.Open(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var entry Entry
	err = json.NewDecoder(f).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Store) Del(path string) {
	entry, err := s.Get(path)
	if err != nil || entry == nil {
		return
	}
	os.Remove(entry.FullPath)
	os.Remove(s.metadataPath(path))
	s.removeParentIfEmpty(entry.FullPath)
}

func (s *Store) Set(path string, ttl time.Duration, src io.Reader, headers map[string]string) error {
	metaPath := s.metadataPath(path)
	parent := filepath.Dir(metaPath)
	os.MkdirAll(parent, os.ModePerm)

	f1, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer f1.Close()

	entity := Entry{
		FullPath:   s.archivePath(path),
		Expiration: time.Now().Add(ttl),
		Header:     headers,
	}

	defer func() {
		json.NewEncoder(f1).Encode(&entity)
	}()

	f2, err := os.Create(entity.FullPath)
	if err != nil {
		return err
	}
	defer f2.Close()

	if headers["Content-Encoding"] != "" || !gzipx.IsSupportedContentType(headers["Content-Type"]) {
		_, err = io.Copy(f2, src)
		return err
	}

	entity.Header["Content-Encoding"] = "gzip"
	gz := gzipx.NewGZWriter(f2)
	defer gz.Close()

	_, err = io.Copy(gz, src)
	return err
}

func (s *Store) removeParentIfEmpty(path string) {
	parent := filepath.Dir(path)
	for {
		if parent == s.baseDir {
			break
		}
		err := os.Remove(parent)
		if err != nil {
			break
		}
		parent = filepath.Dir(parent)
	}
}
