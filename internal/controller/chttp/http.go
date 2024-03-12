package chttp

import (
	"fmt"
	"net/http"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/hhypermedia"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/hstatic"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/middleware"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/rapi"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/router"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/utilb"
	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	HTTP struct {
		mux http.Handler
	}

	Config struct {
		CORSOrigins []string `conf:"cors_origins"`
	}

	Service interface {
		rapi.Service
	}

	OAuth interface {
		middleware.AuthOAuth
		rapi.OAuth
	}
)

func New(svc Service, oa OAuth, cfg Config, log logger.Logger) (*HTTP, error) {
	mux0 := router.NewMux()

	hpmd, err := hhypermedia.New("./web/template", log.WithName("Hypermedia"))
	if err != nil {
		return nil, fmt.Errorf("initialize hypermedia failed: %w", err)
	}

	mux0.NoMethodHandle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux0.NotFoundHandle(func(w http.ResponseWriter, r *http.Request) {
		switch utilb.HeaderAcceptFirst(r.Header.Get("Accept"), "text/html") {
		case "text/html":
			hpmd.Error(r.Context(), w, "", http.StatusNotFound)
		default:
			utilb.ResponseErr404(w)
		}
	})

	mux0.Pre(middleware.Logger(log), func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wHeader := w.Header()
			wHeader.Set("Server", bagicore.Project)
			wHeader.Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")
			wHeader.Set("X-Frame-Options", "DENY")
			next.ServeHTTP(w, r)
		})
	}, middleware.CORS(func(opt *middleware.CORSOption) {
		opt.AllowedMethod = append(opt.AllowedMethod, http.MethodGet)
		opt.AllowedHeader = append(opt.AllowedHeader, "Content-Type")
		opt.SkipOrigin = true
	}))

	mux0.Group(func(mux router.Mux) {
		mux.Pre(middleware.CORSProcess)

		mux.MethodGet("/", hpmd.IndexHandler)
	})

	mux0.Group(func(mux router.Mux) {
		mux.Pre(middleware.CORS(func(opt *middleware.CORSOption) {
			opt.AllowedMethod = append(opt.AllowedMethod, http.MethodHead)
		}), middleware.CORSProcess)

		stac := hstatic.New(hpmd, log.WithName("Static"))
		stacVal := []string{
			"favicon.ico",
			"humans.txt",
			"robots.txt",
		}
		stacMtd := []string{http.MethodGet, http.MethodHead}
		stacDir := http.Dir("./web/static")
		for _, val := range stacVal {
			mux.MultiMethod(stacMtd, "/"+val, stac.File(stacDir, val))
		}
		mux.MultiMethod(stacMtd, "/static/*", stac.Directory(stacDir, "/static"))
	})

	sapi, err := rapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("initialize swagger failed: %w", err)
	}

	mux0.Sub("/api", func(mux1 router.Mux) {
		log := log.WithName("API")

		mux1.Pre(middleware.CORS(func(opt *middleware.CORSOption) {
			opt.AllowedOrigin = cfg.CORSOrigins
			opt.AllowedMethod = append(opt.AllowedMethod, http.MethodPatch, http.MethodPost)
			opt.AllowedMethod = append(opt.AllowedMethod, http.MethodDelete)
			opt.AllowCredentials = true
			opt.SkipOrigin = false
		}), middleware.CORSProcess, middleware.Auth(oa))

		iapi := rapi.NewAPI(svc, oa, log)
		mapi := mux1.Underlying(rapi.Middleware(sapi, iapi.Authentication))
		rapi.HandlerFromMuxWithBaseURL(iapi, mapi, "/v0")
	})

	return &HTTP{mux: mux0}, nil
}

func (ctr HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctr.mux.ServeHTTP(w, r)
}
