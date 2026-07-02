package delivery

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewLogging(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		var statusCode codes.Code
		if err != nil {
			s, _ := status.FromError(err)
			statusCode = s.Code()
		} else {
			statusCode = codes.OK
		}

		logger.Info("request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.Int("status_code", int(statusCode)),
			zap.String("error", statusCode.String()),
		)

		return resp, err
	}
}
