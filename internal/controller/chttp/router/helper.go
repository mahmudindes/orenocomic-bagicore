package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	URLParam            = chi.URLParamFromCtx
	WrapResponseWritter = middleware.NewWrapResponseWriter
)
