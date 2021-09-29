package errors

import (
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
