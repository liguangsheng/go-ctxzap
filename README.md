# ctxzap

# install
```
go get github.com/liguangsheng/go-ctxzap
```

# example

## simple
```go
package main

func init() {
	ctxzap.BetterDefault() // this will replace global zap logger with a better default logger
}

func SomeFunction(ctx context.Context) {
    logger := ctxzap.L(ctx)
    logger.Info("some log") 
    ...
}

func main() {
	originContext := context.Background()
	newCtx := zapctx.ToContext(originContext, zap.L())
	SomeFunction(newCtx)
}
```

## gin middleware

```go
engine := gin.Default()
engine.Use(ctxzap.Gin(zap.L()))
```

## grpc middleware
```go
grpc.NewServer(
	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		zapctx.UnaryServerInterceptor(zap.L()),
	)))
```
