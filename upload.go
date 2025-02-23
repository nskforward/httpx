package httpx

import (
	"io"
	"mime/multipart"
	"net/http"
)

type Files struct {
	r *http.Request
}

func (ctx *Ctx) LoadFiles(maxMemory int64) (*Files, error) {
	err := ctx.Request().ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}
	return &Files{ctx.Request()}, nil
}

func (files *Files) GetFile(name string, handler func(*multipart.FileHeader, io.Reader) error) error {
	f, h, err := files.r.FormFile(name)
	if err != nil {
		return err
	}
	defer f.Close()
	err = handler(h, f)
	if err != nil {
		return err
	}
	return nil
}
