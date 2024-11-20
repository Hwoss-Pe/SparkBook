package logger

type NoOpLogger struct {
}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n NoOpLogger) Debug(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (n NoOpLogger) Info(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (n NoOpLogger) Warn(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (n NoOpLogger) Error(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (n NoOpLogger) With(args ...Field) Logger {
	//TODO implement me
	panic("implement me")
}
