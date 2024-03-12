package hstatic

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	static struct {
		hypermedia hypermedia
		logger     logger.Logger
	}

	hypermedia interface {
		Error(ctx context.Context, w http.ResponseWriter, error string, code int)
	}
)

func New(hpmd hypermedia, log logger.Logger) static {
	return static{hypermedia: hpmd, logger: log}
}

func (stc static) File(hfs http.FileSystem, name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := stc.logger.WithContext(ctx)

		f, err := hfs.Open(name)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			stc.hypermedia.Error(ctx, w, "", http.StatusNotFound)
			return
		case err != nil:
			stc.hypermedia.Error(ctx, w, "", http.StatusInternalServerError)
			log.ErrMessage(err, "File open failed.")
			return
		}
		defer f.Close()

		s, err := f.Stat()
		if err != nil {
			stc.hypermedia.Error(ctx, w, "", http.StatusInternalServerError)
			log.ErrMessage(err, "File stat failed.")
			return
		}

		if s.IsDir() {
			stc.hypermedia.Error(ctx, w, "", http.StatusNotFound)
			return
		}

		http.ServeContent(w, r, s.Name(), s.ModTime(), f)
	}
}

func (stc static) Directory(hfs http.FileSystem, prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := stc.logger.WithContext(ctx)

		fpath := r.URL.Path
		fpath = path.Clean(strings.TrimPrefix(fpath, prefix))

		if strings.HasSuffix(fpath, "/") {
			stc.hypermedia.Error(ctx, w, "", http.StatusNotFound)
			return
		}

		f, err := hfs.Open(fpath)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			stc.hypermedia.Error(ctx, w, "", http.StatusNotFound)
			return
		case err != nil:
			stc.hypermedia.Error(ctx, w, "", http.StatusInternalServerError)
			log.ErrMessage(err, "File open failed.")
			return
		}
		defer f.Close()

		s, err := f.Stat()
		if err != nil {
			stc.hypermedia.Error(ctx, w, "", http.StatusInternalServerError)
			log.ErrMessage(err, "File stat failed.")
			return
		}

		if s.IsDir() {
			stc.hypermedia.Error(ctx, w, "", http.StatusNotFound)
			return
		}

		http.ServeContent(w, r, s.Name(), s.ModTime(), f)
	}
}
