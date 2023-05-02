package slogx

import (
	"context"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

var _ slog.Handler = (*Handler)(nil)

type HandlerOptions struct {
	Level         slog.Leveler
	Prefix        string
	AddSource     bool
	AddCaller     bool
	AddStackTrace bool
	ReplaceAttr   func(groups []string, a slog.Attr) slog.Attr
}

type Handler struct {
	opts HandlerOptions // 控制写什么
	e    Encoder        // 控制怎么写
	w    io.Writer      // 控制往那写

	group string
	core  slog.Handler
}

func NewHandler(w io.Writer, opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return opts.NewHandler(w)
}

func (opts HandlerOptions) NewHandler(w io.Writer) *Handler {
	return &Handler{
		opts: opts,
		w:    w,
		core: slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: opts.ReplaceAttr,
		}.NewJSONHandler(w),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.core.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	return h.core.Handle(ctx, record)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		opts: h.opts,
		e:    h.e,
		w:    h.w,
		core: h.core.WithAttrs(attrs),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		opts:  h.opts,
		e:     h.e,
		w:     h.w,
		group: name,
		core:  h.core.WithGroup(name),
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
