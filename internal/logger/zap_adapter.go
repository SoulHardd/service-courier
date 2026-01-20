package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: l}
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, toZapFields(fields)...)
}

func toZapFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zf = append(zf, zap.Any(f.Key, f.Value))
	}
	return zf
}
