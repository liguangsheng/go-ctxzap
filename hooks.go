package ctxzap

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

type HookFunc func(ctx context.Context) context.Context

func OpenTraceFields(ctx context.Context) context.Context {
	zapCtx := Extract(ctx)
	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		zapCtx.AddFields(
			zap.String("trace.traceid", spanCtx.TraceID.String()),
			zap.String("trace.spanid", spanCtx.SpanID.String()),
			zap.Bool("trace.sampled", spanCtx.IsSampled()),
		)
	}
	return ctx
}
