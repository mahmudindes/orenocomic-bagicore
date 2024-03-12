package hhypermedia

import (
	"html/template"
	"net/http"
	"os"

	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type Hypermedia struct {
	tIndex *template.Template
	tError *template.Template
	logger logger.Logger
}

func New(dir string, log logger.Logger) (*Hypermedia, error) {
	hype, dirFS := &Hypermedia{logger: log}, os.DirFS(dir)

	tIndex := template.New("index")
	if _, err := tIndex.ParseFS(dirFS, "index.html"); err != nil {
		return nil, err
	}
	hype.tIndex = tIndex

	tError, err := template.Must(tIndex.Clone()).Funcs(template.FuncMap{
		"statusText": http.StatusText,
	}).ParseFS(dirFS, "index-error.html")
	if err != nil {
		return nil, err
	}
	hype.tError = tError

	return hype, nil
}
