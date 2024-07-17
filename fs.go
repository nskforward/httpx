package httpx

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nskforward/httpx/types"
)

func ServeFile(filePath string) types.Handler {
	fi, err := os.Stat(filePath)
	if err != nil {
		panic(fmt.Errorf("cannot find the file: %w", err))
	}
	if fi.IsDir() {
		panic(fmt.Errorf("file cannot be a dir: %s", filePath))
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		http.ServeFile(w, r, filePath)
		return nil
	}
}

func ServeDir(dir string) types.Handler {
	fi, err := os.Stat(dir)
	if err != nil {
		panic(dir)
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("directory cannot be a file: %s", dir))
	}

	fserver := http.FileServer(http.FS(ProtectedFS{fs: os.DirFS(dir)}))

	return func(w http.ResponseWriter, r *http.Request) error {
		fserver.ServeHTTP(w, r)
		return nil
	}
}

type ProtectedFS struct {
	fs fs.FS
}

func (pfs ProtectedFS) Open(name string) (fs.File, error) {
	if strings.Contains(name, "/..") {
		return nil, fs.ErrNotExist
	}
	f, err := pfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if s.IsDir() {
		f.Close()
		f, err = pfs.fs.Open(filepath.Join(name, "index.html"))
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
