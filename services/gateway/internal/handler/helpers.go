package handler

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeGrpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)

	if !ok {
		writeJson(w, http.StatusInternalServerError, map[string]string{
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
	}

	writeJson(w, httpStatus, map[string]string{
		"error": st.Message(),
	})
}
