package httpx

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServeFile(filePath string) Handler {
	fi, err := os.Stat(filePath)
	if err != nil {
		panic(fmt.Errorf("cannot find the file: %w", err))
	}
	if fi.IsDir() {
		panic(fmt.Errorf("file cannot be a dir: %s", filePath))
	}

	return func(ctx *Context) error {
		http.ServeFile(ctx.w, ctx.req, filePath)
		return nil
	}
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
func ServeDir(dir, stripPrefix string) Handler {
	fi, err := os.Stat(dir)
	if err != nil {
		panic(dir)
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("directory cannot be a file: %s", dir))
	}

	fserver := http.FileServer(http.FS(ProtectedFS{fs: os.DirFS(dir)}))

	return func(ctx *Context) error {
		if stripPrefix != "" {
			ctx.req.URL.Path = strings.TrimPrefix(ctx.req.URL.Path, stripPrefix)
			ctx.req.URL.RawPath = strings.TrimPrefix(ctx.req.URL.RawPath, stripPrefix)
		}
		fserver.ServeHTTP(ctx.w, ctx.req)
		return nil
	}
}

type ProtectedFS struct {
	fs fs.FS
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
		f, err = pfs.fs.Open(filepath.Join(name, "index.html"))
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
