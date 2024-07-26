package cache

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	dir     string
	buckets sync.Map
	tags    *RadixNode
}

func NewStore(dir string) *Store {
	return &Store{
		dir:  dir,
		tags: &RadixNode{},
	}
}

func (s *Store) ChangeDir(dir string) {
	s.dir = dir
}

func (s *Store) GetBucket(key string) *Bucket {
	res, ok := s.buckets.Load(key)
	if !ok {
		return nil
	}
	return res.(*Bucket)
}

func (s *Store) GetOrCreateBucket(key string) *Bucket {
	dir := filepath.Join(s.dir, Hash(key))
	res, loaded := s.buckets.LoadOrStore(key, &Bucket{key: key, dir: dir, store: s})
	if !loaded {
		fmt.Println("create bucket:", key)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			slog.Error("cannot create a cache bucket folder", "error", err)
		}
	}
	return res.(*Bucket)
}

func (s *Store) GetTag(name string) *Tag {
	node := s.tags.GetOrCreate([]byte(name))
	if node.tag == nil {
		node.tag = &Tag{}
	}
	return node.tag
}
