package httpchi

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type createUserRequest struct {
	Name  string `json:"name" validate:"required,max=50"`
	Email string `json:"email" validate:"required,email,max=100"`
}

type createUserResponse struct {
	Id int `json:"id"`
}

type updateUserRequest struct {
	Id    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,max=50"`
	Email string `json:"email" validate:"required,email,max=100"`
}

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = jsoniter.NewEncoder(w).Encode(errorResponse{
		ErrorCode: statusCode,
		Message:   err.Error(),
	})
}

func writeInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = jsoniter.NewEncoder(w).Encode(errorResponse{
		ErrorCode: http.StatusInternalServerError,
		Message:   "internal server error",
	})
}

func writeBadRequestError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_ = jsoniter.NewEncoder(w).Encode(errorResponse{
		ErrorCode: http.StatusBadRequest,
		Message:   "invalid request data",
	})
}
