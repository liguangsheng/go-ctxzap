package ctxzap

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/liguangsheng/go-randstr"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"time"
)

const traceIDKey = "X-Trace-Id"

type _writer struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r _writer) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Gin(logger *zap.Logger, hooks ...HookFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader(traceIDKey)
		if traceID == "" {
			traceID = randstr.String(32)
		}
		ctx.Header(traceIDKey, traceID)
		logger = logger.With(zap.String(traceIDKey, traceID),
			zap.String("path", ctx.Request.URL.Path))
		newCtx := ToContext(ctx.Request.Context(), logger)

		for _, hook := range hooks {
			newCtx = hook(newCtx)
		}
		ctx.Request = ctx.Request.WithContext(newCtx)

		var buf bytes.Buffer
		tee := io.TeeReader(ctx.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		ctx.Request.Body = ioutil.NopCloser(&buf)

		w := &_writer{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		logger.Info("gin payload",
			zap.Any("request", json.RawMessage(body)),
			zap.Any("headers", ctx.Request.Header))

		start := time.Now()
		ctx.Next()
		elapsed := time.Since(start)

		logger.Info("gin payload",
			zap.Any("response", json.RawMessage(w.body.String())),
			zap.Duration("elapsed", elapsed))
	}
}
