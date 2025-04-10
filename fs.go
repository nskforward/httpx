package httpx

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ProtectedFS struct {
	fs        fs.FS
	indexFile string
}

func (pfs ProtectedFS) Open(name string) (fs.File, error) {
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
		f, err = pfs.fs.Open(filepath.Join(name, pfs.indexFile))
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

/*
	 	ServeDir servers dir files dynamically based on request URL path:

		file_to_serve = dir_path + url_path
		file_to_serve = dir_path + StripPrefix(url_path, prefix)

		Route("/", httpx.ServeDir("/data/static"))
		GET	/1.html	 -->  200 OK  /data/static/1.html
		GET /2.html	 -->  200 OK  /data/static/2.html

		Route("/data/static", httpx.ServeDir("/data/static"))
		GET	/data/static/1.html	 -->  404 Not Found  /data/static/data/static/1.html
		GET /data/static/2.html	 -->  404 Not Found  /data/static/data/static/2.html

		r.Route("/static", httpx.ServeDir("/data/static"), middleware.StripPrefix("/static"))
		GET	/static/1.html  -->  200 OK  /data/static/1.html
		GET	/static/2.html  -->  200 OK  /data/static/2.html
*/
func ServeFolder(dir, indexFile, stripPrefix string) Handler {

	path, err := resolvePath(dir)
	if err != nil {
		panic(dir)
	}

	fserver := http.FileServer(http.FS(ProtectedFS{fs: os.DirFS(path), indexFile: indexFile}))

	return func(req *http.Request, resp *Response) error {
		if stripPrefix != "" {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, stripPrefix)
			req.URL.RawPath = strings.TrimPrefix(req.URL.RawPath, stripPrefix)
		}
		fserver.ServeHTTP(resp.w, req)
		return nil
	}
}

func resolvePath(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return filepath.EvalSymlinks(path)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("directory cannot be a file: %s", path)
	}

	return path, nil
}
