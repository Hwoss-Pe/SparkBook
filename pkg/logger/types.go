package logger

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) Logger
}

type Field struct {
	Key   string
	Value any
}
