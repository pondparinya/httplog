package httplog

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	*zap.SugaredLogger
}

type zapOtelLogger struct {
	*otelzap.SugaredLogger
}

type ZapConfig struct {
	Level               string // debug, info, warn, error, panic, fatal
	OutputPaths         []string
	Format              string // json, console
	DisableFunctionName bool
	DisableLog          bool
	OutputFile          string
	EnableOtel          bool
}

// NewZapLogger creates a new logger instance with underlying zap logger engine
func NewZapLogger(cfg ZapConfig) (Logger, error) {
	config := zap.NewProductionConfig()

	var ws zapcore.WriteSyncer

	switch strings.ToLower(cfg.Level) {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		config.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	switch strings.ToLower(cfg.Format) {
	case "json", "console":
		config.Encoding = strings.ToLower(cfg.Format)
	default:
		config.Encoding = "json"
	}

	if cfg.DisableLog {
		zap.RegisterSink(NoOpWriterKey, func(u *url.URL) (zap.Sink, error) {
			return NewNoOpWriter(), nil
		})
		config.OutputPaths = []string{fmt.Sprintf("%s:whatever", NoOpWriterKey)}
	} else {
		config.OutputPaths = cfg.OutputPaths
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = KeyTimestamp
	config.EncoderConfig.NameKey = KeyServiceName
	config.EncoderConfig.MessageKey = KeyMessage

	if !cfg.DisableFunctionName {
		config.EncoderConfig.FunctionKey = KeyFunctionName
	}

	if !cfg.DisableLog && cfg.OutputFile != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.OutputFile,
			MaxAge:     0,
			MaxBackups: 0,
			Compress:   true,
			LocalTime:  true,
		}

		ws = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(fileWriter),
		)

		config.OutputPaths = nil

	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	var encoder zapcore.Encoder
	if config.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(config.EncoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(config.EncoderConfig)
	}

	l, err := config.Build(
		zap.AddCallerSkip(1),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				encoder,
				ws,
				config.Level,
			)
		}),
	)
	if err != nil {
		return nil, err
	}

	if cfg.EnableOtel {
		// Wrap original zap with otelzap
		zapOtelLog := &zapOtelLogger{otelzap.New(l).Sugar()}
		// Set the global logger to the otelzap logger
		SetGlobalLogger(zapOtelLog)
		return zapOtelLog, nil
	} else {
		zapLog := &zapLogger{l.Sugar()}
		SetGlobalLogger(zapLog)
		return zapLog, nil
	}
}

func (l *zapLogger) Debugf(msg string, args ...any) {
	l.SugaredLogger.Debugf(fmt.Sprintf(msg, args...))
}

func (l *zapLogger) Warnf(msg string, args ...any) {
	l.SugaredLogger.Warnf(fmt.Sprintf(msg, args...))
}

func (l *zapLogger) Infof(msg string, args ...any) {
	l.SugaredLogger.Infof(fmt.Sprintf(msg, args...))
}

func (l *zapLogger) Errorf(msg string, args ...any) {
	l.SugaredLogger.Errorf(fmt.Sprintf(msg, args...))
}

func (l *zapLogger) Panicf(msg string, args ...any) {
	l.SugaredLogger.Panicf(msg, args...)
}

func (l *zapLogger) Debugw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Debugw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapLogger) Infow(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Infow(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapLogger) Warnw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Warnw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapLogger) Errorw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Errorw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapLogger) Panicw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Panicw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapLogger) Debugcf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Debugw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapLogger) Infocf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Infow(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapLogger) Warncf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Warnw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapLogger) Errorcf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Errorw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapLogger) Paniccf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Panicw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapLogger) Debugcw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Debugw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapLogger) Infocw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Infow(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapLogger) Warncw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Warnw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapLogger) Errorcw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Errorw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapLogger) Paniccw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Panicw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{l.SugaredLogger.Named(name)}
}

// ====================================================================================

func (l *zapOtelLogger) Debugf(msg string, args ...any) {
	l.SugaredLogger.Debugf(fmt.Sprintf(msg, args...))
}

func (l *zapOtelLogger) Warnf(msg string, args ...any) {
	l.SugaredLogger.Warnf(fmt.Sprintf(msg, args...))
}

func (l *zapOtelLogger) Infof(msg string, args ...any) {
	l.SugaredLogger.Infof(fmt.Sprintf(msg, args...))
}

func (l *zapOtelLogger) Errorf(msg string, args ...any) {
	l.SugaredLogger.Errorf(fmt.Sprintf(msg, args...))
}

func (l *zapOtelLogger) Panicf(msg string, args ...any) {
	l.SugaredLogger.Panicf(msg, args...)
}

func (l *zapOtelLogger) Debugw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Debugw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapOtelLogger) Infow(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Infow(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapOtelLogger) Warnw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Warnw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapOtelLogger) Errorw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Errorw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapOtelLogger) Panicw(fields []Field, msg string, args ...any) {
	l.SugaredLogger.Panicw(fmt.Sprintf(msg, args...), ParseFields(fields)...)
}

func (l *zapOtelLogger) Debugcf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Debugw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapOtelLogger) Infocf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Infow(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapOtelLogger) Warncf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Warnw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapOtelLogger) Errorcf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Errorw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapOtelLogger) Paniccf(ctx context.Context, msg string, args ...any) {
	l.SugaredLogger.Panicw(fmt.Sprintf(msg, args...), ParseRequestContext(ctx)...)
}

func (l *zapOtelLogger) Debugcw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Ctx(ctx).Debugw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapOtelLogger) Infocw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Ctx(ctx).Infow(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapOtelLogger) Warncw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Ctx(ctx).Warnw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapOtelLogger) Errorcw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Ctx(ctx).Errorw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapOtelLogger) Paniccw(ctx context.Context, fields []Field, msg string, args ...any) {
	l.SugaredLogger.Ctx(ctx).Panicw(fmt.Sprintf(msg, args...), append(ParseRequestContext(ctx), ParseFields(fields)...)...)
}

func (l *zapOtelLogger) Named(name string) Logger {
	return &zapOtelLogger{otelzap.New(l.SugaredLogger.Named(name).Desugar()).Sugar()}
}
