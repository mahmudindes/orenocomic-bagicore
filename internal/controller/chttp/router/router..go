package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Mux struct {
	chiMux *chi.Mux
}

func NewMux() Mux {
	return Mux{chi.NewMux()}
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.chiMux.ServeHTTP(w, r)
}

func (mux Mux) NoMethodHandle(handler http.HandlerFunc) {
	mux.chiMux.MethodNotAllowed(handler)
}

func (mux Mux) NoMethodHandler() http.HandlerFunc {
	return mux.chiMux.MethodNotAllowedHandler()
}

func (mux Mux) NotFoundHandle(handler http.HandlerFunc) {
	mux.chiMux.NotFound(handler)
}

func (mux Mux) NotFoundHandler() {
	mux.chiMux.NotFoundHandler()
}

func (mux Mux) Pre(middlewares ...func(http.Handler) http.Handler) {
	mux.chiMux.Use(middlewares...)
}

func (mux Mux) Sub(pattern string, fn func(mux Mux)) Mux {
	sub := Mux{chiMux: chi.NewRouter()}
	mux.chiMux.Mount(pattern, sub.chiMux)
	if fn != nil {
		fn(sub)
	}
	return sub
}

func (mux Mux) Group(fn func(mux Mux)) Mux {
	group := Mux{chiMux: mux.chiMux.With().(*chi.Mux)}
	if fn != nil {
		fn(group)
	}
	return group
}

func (mux Mux) Mount(pattern string, h http.Handler) {
	mux.chiMux.Mount(pattern, h)
}

func (mux Mux) MethodDelete(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.chiMux.Delete(pattern, handler)
}

func (mux Mux) MethodGet(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.chiMux.Get(pattern, handler)
}

func (mux Mux) MethodPatch(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.chiMux.Patch(pattern, handler)
}

func (mux Mux) MethodPost(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.chiMux.Post(pattern, handler)
}

func (mux Mux) MethodPut(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.chiMux.Put(pattern, handler)
}

func (mux Mux) MultiMethod(methods []string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	for _, method := range methods {
		mux.chiMux.MethodFunc(method, pattern, handler)
	}
}

func (mux Mux) Underlying(middlewares ...func(http.Handler) http.Handler) chi.Router {
	return mux.chiMux.With(middlewares...)
}
