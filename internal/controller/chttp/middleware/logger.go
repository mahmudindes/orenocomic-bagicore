package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/router"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller/chttp/utilb"
	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			log = log.WithContext(r.Context())

			requestURI := r.URL.RequestURI()
			log.Message(
				"HTTP Request "+r.Method+" "+requestURI,
				"httpRequest", map[string]any{
					"url":       utilb.GetScheme(r) + "://" + r.Host + requestURI,
					"method":    r.Method,
					"path":      r.URL.Path,
					"proto":     r.Proto,
					"userAgent": r.UserAgent(),
					"referer":   r.Referer(),
				},
			)

			ww := router.WrapResponseWritter(w, r.ProtoMajor)
			defer func() {
				status := ww.Status()
				log.Message(
					"HTTP Response "+strconv.FormatInt(int64(status), 10)+" "+http.StatusText(status),
					"httpResponse", map[string]any{
						"status":  status,
						"size":    ww.BytesWritten(),
						"elapsed": time.Since(now).Nanoseconds(),
					},
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
