package ctxzap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Context = *_context

type _context struct {
	logger *zap.Logger
	fields []zapcore.Field
}

func (c *_context) Logger() *zap.Logger {
	return c.logger.With(c.fields...)
}

func (c *_context) AddFields(fields ...zapcore.Field) {
	c.logger = c.logger.With(fields...)
	c.fields = append(c.fields, fields...)
}

type _marker struct{}

var (
	_key      = _marker{}
	nopLogger = zap.NewNop()
)

// ToContext place *zap.Logger to context
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, _key, &_context{logger: logger})
}

// Extract
func Extract(ctx context.Context) *_context {
	if ctx == nil {
		return nil
	}
	if _ctx, ok := ctx.Value(_key).(*_context); ok {
		return _ctx
	}
	return nil
}

// L return logger from context or zap.L()
func L(ctx context.Context) *zap.Logger {
	_ctx := Extract(ctx)
	if _ctx != nil {
		return _ctx.Logger()
	}
	return zap.L()
}

// S return logger from context or zap.S()
func S(ctx context.Context) *zap.SugaredLogger {
	return L(ctx).Sugar()
}

// N return logger from context or nopLogger
func N(ctx context.Context) *zap.Logger {
	_ctx := Extract(ctx)
	if _ctx != nil {
		return _ctx.Logger()
	}
	return nopLogger
}
