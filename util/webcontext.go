package util

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/auth"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type WebContext struct {
	c *gin.Context
}

func NewWebContext(g *gin.Context) *WebContext {
	return &WebContext{g}
}

func (w *WebContext) GetAppId() string {
	return auth.GetAppId(w.c)
}

func (w *WebContext) GetResourceId() string {
	return w.Param("id")
}

func (w *WebContext) SaveFileTo(file *multipart.FileHeader, dst string) error {
	return w.c.SaveUploadedFile(file, dst)
}

// for future
func (w *WebContext) Compress(header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return errors.New("could not open the file header")
	}

	buff := make([]byte, header.Size)
	if _, err = file.Read(buff); err != nil {
		return errors.New("could not read from the file")
	}

	var b bytes.Buffer
	writer := gzip.NewWriter(&b)

	writer.Name = "Demo.zip"

	if _, err = writer.Write(buff); err != nil {
		return errors.New("could not write to the file")
	}
	err = writer.Close()
	if err != nil {
		return errors.New("could not write to the file")
	}

	err = ioutil.WriteFile("./__uploads__/zipped.gz", buff, os.ModePerm)

	if err != nil {
		return errors.New("could not write to the file")
	}

	return nil
}

func (w *WebContext) FormFile() (*multipart.FileHeader, error) {
	return w.c.FormFile("file")
}

func (w *WebContext) Form() url.Values {
	return w.c.Request.Form
}

func (w *WebContext) Query() url.Values {
	return w.c.Request.URL.Query()
}

func (w *WebContext) QueryParam(key string) string {
	query := w.Query()
	if len(query[key]) == 0 {
		return ""
	}
	return query[key][0]
}

func (w *WebContext) Param(key string) string {
	return w.c.Param(key)
}

func (w *WebContext) Ok(data interface{}) {
	w.Respond(http.StatusOK, data)
}

func (w *WebContext) NotFound(data interface{}) {
	w.Respond(http.StatusNotFound, data)
}

func (w *WebContext) UnprocessableEntity(data interface{}) {
	w.Respond(http.StatusUnprocessableEntity, data)
}

func (w *WebContext) InternalServerError(data interface{}) {
	w.Respond(http.StatusInternalServerError, data)
}

func (w *WebContext) Respond(status int, data interface{}) {
	w.c.JSON(status, data)
}

func (w *WebContext) RespondWithFile(filePath, filename string) {
	w.c.FileAttachment(filePath, filename)
}
