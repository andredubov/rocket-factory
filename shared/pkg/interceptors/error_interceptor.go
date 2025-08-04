package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/shared/pkg/errors"
)

// UnaryErrorInterceptor handles error conversion for unary RPC calls
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, convertError(err, info.FullMethod)
		}
		return resp, nil
	}
}

// convertError converts business errors to appropriate gRPC errors
func convertError(err error, method string) error {
	// Check if it's a businessError
	if businessErr := errors.GetBusinessError(err); businessErr != nil {
		grpcStatus := errors.BusinessErrorToGRPCStatus(businessErr)
		log.Printf(
			"BusinessError in method %s: code=%d, message=%s",
			method,
			businessErr.Code(),
			businessErr.Error(),
		)
		return grpcStatus.Err()
	}

	// Check if it's already a gRPC status error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// For unknown errors, return internal error
	log.Printf("Unknown error in method %s: %v", method, err)
	return status.Error(codes.Internal, "internal server error")
}
