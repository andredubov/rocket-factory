package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents business error codes
type ErrorCode int64

const (
	NotFoundErrCode ErrorCode = iota
	InvalidArgumentErrCode
	AlreadyExistsErrCode
)

// businessError represents a structured business error
type businessError struct {
	code ErrorCode
	err  error
}

func (b *businessError) Error() string {
	if b.err != nil {
		return b.err.Error()
	}
	return "unknown business error"
}

func (b *businessError) Unwrap() error {
	return b.err
}

func (b *businessError) Code() ErrorCode {
	return b.code
}

// Constructor functions for different error types
func NewNotFoundError(err error) *businessError {
	return &businessError{
		code: NotFoundErrCode,
		err:  err,
	}
}

func NewInvalidArgumentError(err error) *businessError {
	return &businessError{
		code: InvalidArgumentErrCode,
		err:  err,
	}
}

func NewAlreadyExistsError(err error) *businessError {
	return &businessError{
		code: AlreadyExistsErrCode,
		err:  err,
	}
}

// GetBusinessError returns businessError if err is a business error, nil otherwise
func GetBusinessError(err error) *businessError {
	var businessErr *businessError
	if errors.As(err, &businessErr) {
		return businessErr
	}
	return nil
}

// errorCodeToGRPCCode maps business error codes to gRPC codes
func errorCodeToGRPCCode(code ErrorCode) codes.Code {
	switch code {
	case NotFoundErrCode:
		return codes.NotFound
	case InvalidArgumentErrCode:
		return codes.InvalidArgument
	case AlreadyExistsErrCode:
		return codes.AlreadyExists
	default:
		return codes.Unknown
	}
}

// BusinessErrorToGRPCStatus converts businessError to gRPC status
func BusinessErrorToGRPCStatus(err *businessError) *status.Status {
	grpcCode := errorCodeToGRPCCode(err.Code())
	return status.New(grpcCode, err.Error())
}
