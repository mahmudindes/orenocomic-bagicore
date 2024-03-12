package controller

import (
	"net/http"

	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp"
	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	Controller struct {
		http func() (http.Handler, error)
	}

	Config struct {
		HTTP chttp.Config `conf:"http"`
	}

	service chttp.Service
	oauth   chttp.OAuth
)

func New(svc service, oa oauth, cfg Config, log logger.Logger) Controller {
	controller := Controller{}

	var cHTTP *chttp.HTTP
	controller.http = func() (http.Handler, error) {
		var err error
		if cHTTP == nil {
			cHTTP, err = chttp.New(svc, oa, cfg.HTTP, log.WithName("HTTP"))
		}
		return cHTTP, err
	}

	return controller
}

func (ctr Controller) HTTP() (http.Handler, error) {
	return ctr.http()
}
