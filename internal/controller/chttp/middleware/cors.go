package middleware

import (
	"context"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type (
	CORSOption struct {
		AllowedOrigin, AllowedMethod, AllowedHeader []string
		ExposedHeader                               []string
		AllowCredentials                            bool
		SkipOrigin                                  bool
		MaxAge                                      int
	}

	ctxCORSOption struct{}
)

func CORS(fn func(opt *CORSOption)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if fn == nil {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			opt, ok := ctx.Value(ctxCORSOption{}).(*CORSOption)
			if opt == nil {
				opt = new(CORSOption)
				opt.AllowedMethod = []string{http.MethodOptions}
			}
			if !ok {
				ctx = context.WithValue(ctx, ctxCORSOption{}, opt)
			}
			fn(opt)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CORSProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rHeader, wHeader := r.Header, w.Header()

		opt, ok := ctx.Value(ctxCORSOption{}).(*CORSOption)
		if !ok || opt == nil {
			next.ServeHTTP(w, r)
			return
		}

		origin := rHeader.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		acrm := rHeader.Get("Access-Control-Request-Method")
		preflight := r.Method == http.MethodOptions && acrm != ""

		acao := origin
		if opt.SkipOrigin {
			acao = "*"
		} else {
			if !slices.ContainsFunc(opt.AllowedOrigin, func(s string) bool {
				if pattern, ok := strings.CutPrefix(s, "regex:"); ok {
					match, err := regexp.MatchString(pattern, origin)
					if err != nil {
						return false
					}
					return match
				}
				return s == origin
			}) {
				next.ServeHTTP(w, r)
				return
			}
		}

		if preflight {
			if !slices.Contains(opt.AllowedMethod, strings.ToUpper(acrm)) {
				next.ServeHTTP(w, r)
				return
			}

			acrhList := []string{}
			for _, h := range strings.Split(rHeader.Get("Access-Control-Request-Headers"), ",") {
				acrhList = append(acrhList, http.CanonicalHeaderKey(strings.TrimSpace(h)))
			}
			for _, h := range acrhList {
				if !slices.Contains(opt.AllowedHeader, h) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		wHeader.Set("Access-Control-Allow-Origin", acao)
		if acao != "*" {
			wHeader.Add("Vary", "Origin")
		}

		if preflight {
			if len(opt.AllowedMethod) > 0 {
				wHeader.Set("Access-Control-Allow-Methods", strings.Join(opt.AllowedMethod, ", "))
				wHeader.Add("Vary", "Access-Control-Request-Method")
			}

			if len(opt.AllowedHeader) > 0 {
				wHeader.Set("Access-Control-Allow-Headers", strings.Join(opt.AllowedHeader, ", "))
				wHeader.Add("Vary", "Access-Control-Request-Headers")
			}

			if opt.MaxAge > 0 {
				wHeader.Set("Access-Control-Max-Age", strconv.Itoa(opt.MaxAge))
			}
		} else {
			if len(opt.ExposedHeader) > 0 {
				wHeader.Set("Access-Control-Expose-Headers", strings.Join(opt.ExposedHeader, ", "))
			}
		}

		if opt.AllowCredentials {
			wHeader.Set("Access-Control-Allow-Credentials", "true")
		}

		if preflight {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
