package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	srvrcfg "github.com/kiryu-dev/server-config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"virtualization-technologies/conf"
	"virtualization-technologies/internal/adapters"
	"virtualization-technologies/internal/adapters/pgrepo"
	"virtualization-technologies/internal/ports/httpchi"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync() // flushes buffer, if any
	}()
	sugar := logger.Sugar()
	cfgPath := flag.String("config", "./http_server.yml", "http server configuration file")
	flag.Parse()
	if cfgPath == nil || *cfgPath == "" {
		sugar.Fatal("config file is not specified")
	}
	serverCfg, err := srvrcfg.LoadYamlCfg(*cfgPath)
	if err != nil {
		sugar.Fatalf("failed to load yaml config in path %s: %v", *cfgPath, err)
	}
	sugar.Infof("SERVER CONFIG: %#v\n", serverCfg)
	dbCfg := conf.NewDbConfig()
	sugar.Infof("DB CONFIG: %#v\n", dbCfg)
	conn, err := pgrepo.NewPgConnection(dbCfg)
	if err != nil {
		sugar.Fatalf("unable to create postgres connection pool: %v", err)
	}
	defer func() {
		sugar.Info("closing db connection...")
		if err := conn.Close(context.Background()); err != nil {
			sugar.Infof("failed to close db connection: %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		select {
		case s := <-sigChan:
			return errors.Errorf("captured signal: %v", s)
		}
	})
	var (
		repo   = adapters.New(conn)
		server = httpchi.NewHTTPServer(serverCfg, repo, sugar)
	)
	go func() {
		sugar.Info("http server is starting...")
		if err := server.ListenAndServe(); err != nil {
			sugar.Info(err)
		}
	}()
	if err := errGroup.Wait(); err != nil {
		sugar.Infof("gracefully shutting down the server: %v", err)
	}
	if err := server.Shutdown(context.Background()); err != nil {
		sugar.Infof("failed to shutdown http server: %v", err)
	}
}
