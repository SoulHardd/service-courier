package logger

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Info(msg string, fields ...Field)
}
