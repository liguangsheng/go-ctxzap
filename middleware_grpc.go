package ctxzap

import (
	"context"
	"github.com/liguangsheng/go-randstr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func UnaryServerInterceptor(logger *zap.Logger, hooks ...HookFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error) {

		traceID := randstr.String(32)
		logger = logger.With(zap.String(traceIDKey, traceID),
			zap.String("path", info.FullMethod))

		newCtx := ToContext(ctx, logger)
		for _, hook := range hooks {
			newCtx = hook(newCtx)
		}

		logger.Info("grpc payload",
			zap.Any("request", req))

		start := time.Now()
		resp, err = handler(newCtx, req)
		elapsed := time.Since(start)

		if err != nil {
			logger.Error("grpc error", zap.Error(err))
		}
		logger.Info("grpc payload",
			zap.Any("response", resp),
			zap.Int64("elapsed_ms", elapsed.Microseconds()),
			zap.Error(err))

		return resp, err
	}
}
