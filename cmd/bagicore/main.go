package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
	"github.com/mahmudindes/orenocomic-bagicore/internal/auth"
	"github.com/mahmudindes/orenocomic-bagicore/internal/config"
	"github.com/mahmudindes/orenocomic-bagicore/internal/controller"
	"github.com/mahmudindes/orenocomic-bagicore/internal/datastore"
	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
	"github.com/mahmudindes/orenocomic-bagicore/internal/server"
	"github.com/mahmudindes/orenocomic-bagicore/internal/service"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

var StartTime = time.Now()

func main() {
	godotenv.Load()

	exitCode := mainRun()
	os.Exit(int(exitCode))
}

func mainRun() exitCode {
	log := logger.New()

	log.Message("Starting service.", "version", bagicore.Version)
	defer func() {
		log.Message("Service stopped.", "uptime", time.Since(StartTime))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cfg, err := config.New()
	if err != nil {
		log.ErrMessage(err, "Config initialization failed.")
		return exitError
	}

	ds, err := datastore.New(ctx, cfg.Datastore)
	if err != nil {
		log.ErrMessage(err, "Datastore initialization failed.")
		return exitError
	}
	defer func() {
		if err := ds.Stop(); err != nil {
			log.ErrMessage(err, "Datastore stop failed.")
		}
	}()

	au, err := auth.New(ctx, ds.Redis, cfg.Auth, log)
	if err != nil {
		log.ErrMessage(err, "Auth initialization failed.")
		return exitError
	}

	svc := service.New(ds.Database, au.OAuth)

	ctr := controller.New(svc, au.OAuth, cfg.General.Controller, log)

	svr, err := server.New(ctr, cfg.Server, log.WithName("Server"))
	if err != nil {
		log.ErrMessage(err, "Server initialization failed.")
		svr.Shutdown()
		return exitError
	}
	defer svr.Shutdown()

	select {
	case <-ctx.Done():
		log.Message("Stopping service.")
	case <-svr.ListenAndServe():
		log.Message("Unexpected server stopped.")
	}

	return exitOK
}
