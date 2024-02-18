package httpchi

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"virtualization-technologies/internal/entity/user"
)

func appRouter(mux *chi.Mux, repo user.Repository, logger *zap.SugaredLogger) {
	mux.Get("/users", getAllUsers(repo, logger))
	mux.Get("/users/{id}", getUser(repo, logger))
	mux.Post("/users", createUser(repo, logger))
	mux.Put("/users", updateUser(repo, logger))
	mux.Delete("/users/{id}", deleteUser(repo, logger))
}
