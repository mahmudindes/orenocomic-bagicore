package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	Server struct {
		starter []func()
		stopper []func()
	}

	Config struct {
		HTTP struct {
			Address         string        `conf:"address"`
			ReadTimeout     time.Duration `conf:"read_timeout"`
			WriteTimeout    time.Duration `conf:"write_timeout"`
			ShutdownTimeout time.Duration `conf:"shutdown_timeout"`
		} `conf:"http"`
	}

	controller interface {
		HTTP() (http.Handler, error)
	}
)

func New(ctr controller, cfg Config, log logger.Logger) (Server, error) {
	server := Server{starter: make([]func(), 0), stopper: make([]func(), 0)}

	chttp, err := ctr.HTTP()
	if err != nil {
		return server, fmt.Errorf("initialize http controller failed: %w", err)
	}
	shttp := &http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      chttp,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}
	server.starter = append(server.starter, func() {
		log.Message("HTTP server started.", "host", shttp.Addr)
		if err := shttp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ErrMessage(err, "HTTP server listen and serve failed.")
		}
	})
	server.stopper = append(server.stopper, func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := shttp.Shutdown(ctx); err != nil {
			log.ErrMessage(err, "HTTP server shutdown failed.")
		}
		log.Message("HTTP server stopped.")
	})

	return server, nil
}

func (svr Server) ListenAndServe() <-chan struct{} {
	ch := make(chan struct{}, 1)
	for _, startFn := range svr.starter {
		go func(fn func()) {
			fn()
			ch <- struct{}{}
		}(startFn)
	}
	return ch
}

func (svr Server) Shutdown() {
	var wg sync.WaitGroup
	for _, stopFn := range svr.stopper {
		wg.Add(1)
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(stopFn)
	}
	wg.Wait()
}
