package dot

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

type FastZerologHandler struct {
	logger zerolog.Logger
}

// Enabled implements [slog.Handler].
func (p *FastZerologHandler) Enabled(_ context.Context, l slog.Level) bool {
	var zl zerolog.Level
	switch {
	case l >= slog.LevelError:
		zl = zerolog.ErrorLevel
	case l >= slog.LevelWarn:
		zl = zerolog.WarnLevel
	case l >= slog.LevelInfo:
		zl = zerolog.InfoLevel
	default:
		zl = zerolog.DebugLevel
	}
	return p.logger.GetLevel() <= zl
}

// WithAttrs implements [slog.Handler].
func (p *FastZerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	ctx := p.logger.With()
	for _, a := range attrs {
		val := a.Value.Resolve()
		switch val.Kind() {
		case slog.KindString:
			ctx = ctx.Str(a.Key, val.String())
		case slog.KindInt64:
			ctx = ctx.Int64(a.Key, val.Int64())
		default:
			ctx = ctx.Any(a.Key, val.Any())
		}
	}
	return &FastZerologHandler{logger: ctx.Logger()}
}

// WithGroup implements [slog.Handler].
func (p *FastZerologHandler) WithGroup(name string) slog.Handler {
	return p
}

// Handle implements [slog.Handler].
func (p *FastZerologHandler) Handle(ctx context.Context, r slog.Record) error {
	var e *zerolog.Event
	switch {
	case r.Level >= slog.LevelError:
		e = p.logger.Error()
	case r.Level >= slog.LevelWarn:
		e = p.logger.Warn()
	case r.Level >= slog.LevelInfo:
		e = p.logger.Info()
	default:
		e = p.logger.Debug()
	}
	if e == nil {
		return nil
	}
	if !r.Time.IsZero() {
		e.Time(zerolog.TimestampFieldName, r.Time)
	}
	r.Attrs(func(a slog.Attr) bool {
		val := a.Value.Resolve()
		if a.Equal(slog.Attr{}) {
			return true
		}
		switch val.Kind() {
		case slog.KindString:
			e.Str(a.Key, val.String())
		case slog.KindInt64:
			e.Int64(a.Key, val.Int64())
		case slog.KindBool:
			e.Bool(a.Key, val.Bool())
		case slog.KindFloat64:
			e.Float64(a.Key, val.Float64())
		case slog.KindDuration:
			e.Dur(a.Key, val.Duration())
		case slog.KindTime:
			e.Time(a.Key, val.Time())
		case slog.KindUint64:
			e.Uint64(a.Key, val.Uint64())
		default:
			e.Any(a.Key, val.Any())
		}
		return true
	})

	e.Msg(r.Message)
	return nil
}
