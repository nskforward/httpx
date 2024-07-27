package cache

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/nskforward/httpx/types"
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
		os.MkdirAll(dir, os.ModePerm)
	}
	return res.(*Bucket)
}

func (s *Store) GetTag(name string) *Tag {
	node := s.tags.GetOrCreate([]byte(name))
	if node.tag == nil {
		node.tag = &Tag{name: name}
	}
	return node.tag
}

func (s *Store) Inject(r *http.Request) *http.Request {
	return types.SetParam(r, "cache.store", s)
}

func GetStore(r *http.Request) *Store {
	res := types.GetParam(r, "cache.store")
	if res == nil {
		return nil
	}
	return res.(*Store)
}
