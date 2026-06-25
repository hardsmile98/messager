package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gateway/internal/validation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RequestError(w http.ResponseWriter, err error) {
	if fields := validation.FieldErrors(err); len(fields) > 0 {
		JSON(w, http.StatusBadRequest, map[string]any{
			"error":  "validation failed",
			"fields": fields,
		})
		return
	}

	JSON(w, http.StatusBadRequest, map[string]string{
		"error": "invalid request body",
	})
}

func GRPCError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		JSON(w, http.StatusGatewayTimeout, map[string]string{
			"error": "request timeout",
		})
		return
	}

	st, ok := status.FromError(err)
	if !ok {
		JSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	httpStatus := http.StatusInternalServerError

	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.DeadlineExceeded:
		httpStatus = http.StatusGatewayTimeout
	}

	JSON(w, httpStatus, map[string]string{
		"error": st.Message(),
	})
}
