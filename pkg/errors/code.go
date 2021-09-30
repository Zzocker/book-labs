package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Code used by this project.
type Code int

// Mapping of internal code used by this project
// with gRPC status code.
const (
	CodeNotFound      Code = Code(codes.NotFound)
	CodeUnauthorized  Code = Code(codes.Unauthenticated)
	CodeInvalidInput  Code = Code(codes.InvalidArgument)
	CodeUnexpected    Code = Code(codes.Internal)
	CodeAlreadyExists Code = Code(codes.AlreadyExists)
)

func ErrCodeToHTTP(c Code) int {
	switch c {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeUnexpected:
		return http.StatusInternalServerError
	case CodeAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
