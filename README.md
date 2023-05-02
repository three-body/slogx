# slogx

[![tag](https://img.shields.io/github/tag/three-body/slogx.svg)](https://github.com/three-body/slogx/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20.3-%23007d9c)
[![GoDoc](https://godoc.org/github.com/three-body/slogx?status.svg)](https://pkg.go.dev/github.com/three-body/slogx)
[![Go report](https://goreportcard.com/badge/github.com/three-body/slogx)](https://goreportcard.com/report/github.com/three-body/slogx)
[![Coverage](https://img.shields.io/codecov/c/github/three-body/slogx)](https://codecov.io/gh/three-body/slogx)
[![Contributors](https://img.shields.io/github/contributors/three-body/slogx)](https://github.com/three-body/slogx/graphs/contributors)
[![License](https://img.shields.io/github/license/three-body/slogx)](./LICENSE)


slogx provides enhanced extensions for slog, including handler, writetr, format, middleware, etc.

## 🚀 Install
```go
go get -u golang.org/x/exp/slog
```

```go
import "github.com/three-body/slogx"
```

**Compatibility**: go >= 1.20.3

This library is v0 and follows SemVer strictly. On slog final release (go 1.21), this library will go v1.

No breaking changes will be made to exported APIs before v1.0.0.

## 💡 Usage

### Handler

#### Handler
```go
writer, err := slogx.NewFileWriter()
if err != nil {
    panic(err)
}

h := slogx.HandlerOptions{
    Level: slog.LevelError,
}.NewHandler(writer)

logger := slog.New(h)
slog.SetDefault(logger)

slog.Info("hello world")
```

#### MutilHandler
```go
errHandler := slogx.HandlerOptions{
    Level: slog.LevelError,
}.NewHandler(os.Stderr)

infoHandler := slogx.HandlerOptions{
    Level: slog.LevelInfo,
}.NewHandler(os.Stdout)

h := slogx.NewMultiHandler(errHandler, infoHandler)
logger := slog.New(h)
slog.SetDefault(logger)

slog.Info("hello world")
```
### Writer
Writer is a io.Writer that can write log to a file or a stream.

#### FileWriter
```go
writer, err := slogx.FileWriterOptions{
    Path:             "./logs",
    FileName:         "app.log",
    MaxTime:          0,
    MaxCount:         0,
    RotateTimeLayout: slogx.RotateTimeLayoutEveryHour,
    RotateSize:       100,
    Compress:         false,
}.NewFileWriter()
if err != nil {
    panic(err)
}

logger := slog.New(slog.NewJSONHandler(writer))
slog.SetDefault(logger)

slog.Info("hello world")
```

#### KafkaWriter

#### RedisWriter

#### MysqlWriter

#### NsqWriter

### Formatter
Formatter is a function that can format log entry to a string.

### Middleware

## 📝 License

Copyright 2023 [three-body](https://github.com/three-body).

This project is [Apache-2.0](./LICENSE) licensed.git 