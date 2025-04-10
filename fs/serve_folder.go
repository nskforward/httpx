package fs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nskforward/httpx"
)

func ServeFolder(dir, indexFile, stripPrefix string) httpx.Handler {

	return func(req *http.Request, resp *httpx.Response) error {
		resolved, err := resolveSymlink(dir)
		if err != nil {
			return err
		}

		target := filepath.Join(resolved, strings.TrimPrefix(req.URL.Path, stripPrefix))

		info, err := os.Stat(target)
		if os.IsNotExist(err) {
			return resp.Text(http.StatusNotFound, "file not found")
		}
		if err != nil {
			resp.Logger().Error(err.Error())
			return resp.Text(http.StatusMisdirectedRequest, "file exists but server cannot serve it")
		}

		if info.IsDir() {
			if indexFile != "" {
				http.ServeFile(resp.ResponseWriter(), req, filepath.Join(target, indexFile))
				return nil
			}
			return resp.Text(http.StatusNotFound, "file not found")
		}

		http.ServeFile(resp.ResponseWriter(), req, target)

		return nil
	}
}

func resolveSymlink(path string) (string, error) {
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
