package httpchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	srvcfg "github.com/kiryu-dev/server-config"
	"go.uber.org/zap"
	"virtualization-technologies/internal/entity/user"
)

func NewHTTPServer(cfg *srvcfg.ServerConfig, repo user.Repository, logger *zap.SugaredLogger) *http.Server {
	mux := chi.NewMux()
	appRouter(mux, repo, logger)
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
