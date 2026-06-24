package errors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func Unauthenticated(msg string) error {
	return status.Error(codes.Unauthenticated, msg)
}

func AlreadyExists(msg string) error {
	return status.Error(codes.AlreadyExists, msg)
}

func InternalError(ctx context.Context, msg string, err error) error {
	slog.ErrorContext(ctx, msg, "error", err)

	return status.Error(codes.Internal, msg)
}
