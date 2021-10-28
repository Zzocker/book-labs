package errors

type Code uint32

const (
	CodeInvalidArgument Code = 3
	CodeNotFound        Code = 5
	CodeAlreadyExists   Code = 6
	CodeInternal        Code = 13
	CodeUnauthenticated Code = 16
)
