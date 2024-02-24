package httpchi

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"virtualization-technologies/internal/entity"
	"virtualization-technologies/internal/entity/user"
)

var validate = validator.New()

func getAllUsers(repo user.Repository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := parsePaginationInfo(r.URL.Query())
		if err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		users, err := repo.GetAll(r.Context(), p.offset, p.count)
		if err != nil {
			logger.Error(err)
			writeInternalError(w)
			return
		}
		if err := jsoniter.NewEncoder(w).Encode(users); err != nil {
			logger.Error(err)
			writeInternalError(w)
		}
	}
}

func getUser(repo user.Repository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		u, err := repo.Get(r.Context(), userId)
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			logger.Warn(err)
			writeError(w, entity.ErrUserNotFound, http.StatusNotFound)
		case err != nil:
			logger.Error(err)
			writeInternalError(w)
		default:
			if err := jsoniter.NewEncoder(w).Encode(u); err != nil {
				logger.Error(err)
				writeInternalError(w)
			}
		}
	}
}

func createUser(repo user.Repository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := new(createUserRequest)
		if err := jsoniter.NewDecoder(r.Body).Decode(u); err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		if err := validate.Struct(u); err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		id, err := repo.Create(r.Context(), user.User{
			Name:  u.Name,
			Email: u.Email,
		})
		switch {
		case errors.Is(err, entity.ErrEmailIsAlreadyTaken):
			logger.Warn(err)
			writeError(w, entity.ErrEmailIsAlreadyTaken, http.StatusBadRequest)
		case err != nil:
			logger.Error(err)
			writeInternalError(w)
		default:
			err := jsoniter.NewEncoder(w).Encode(createUserResponse{Id: id})
			if err != nil {
				logger.Error(err)
				writeInternalError(w)
			}
		}
	}
}

func updateUser(repo user.Repository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := new(updateUserRequest)
		if err := jsoniter.NewDecoder(r.Body).Decode(u); err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		if err := validate.Struct(u); err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		err := repo.Update(r.Context(), user.User{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			logger.Warn(err)
			writeError(w, entity.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, entity.ErrEmailIsAlreadyTaken):
			logger.Warn(err)
			writeError(w, entity.ErrEmailIsAlreadyTaken, http.StatusBadRequest)
		case err != nil:
			logger.Error(err)
			writeInternalError(w)
		}
	}
}

func deleteUser(repo user.Repository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Warn(err)
			writeBadRequestError(w)
			return
		}
		u, err := repo.Delete(r.Context(), userId)
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			logger.Warn(err)
			writeError(w, entity.ErrUserNotFound, http.StatusNotFound)
		case err != nil:
			logger.Error(err)
			writeInternalError(w)
		default:
			if err := jsoniter.NewEncoder(w).Encode(u); err != nil {
				logger.Error(err)
				writeInternalError(w)
			}
		}
	}
}

type pagination struct {
	offset uint64
	count  uint64
}

func parsePaginationInfo(queryParams url.Values) (pagination, error) {
	result := pagination{
		offset: 0,
		count:  math.MaxUint64,
	}
	startStr := strings.TrimSpace(queryParams.Get("start"))
	endStr := strings.TrimSpace(queryParams.Get("end"))
	if startStr == "" && endStr == "" {
		return result, nil
	}
	if startStr != "" {
		start, err := strconv.ParseUint(startStr, 10, 64)
		if err != nil {
			return pagination{}, errors.WithMessage(err, "parse uint")
		}
		result.offset = start
	}
	if endStr != "" {
		end, err := strconv.ParseUint(endStr, 10, 64)
		if err != nil {
			return pagination{}, errors.WithMessage(err, "parse uint")
		}
		if result.offset > end {
			return pagination{}, errors.New("pagination info: start greater than end")
		}
		result.count = end - result.offset + 1
	}
	return result, nil
}
