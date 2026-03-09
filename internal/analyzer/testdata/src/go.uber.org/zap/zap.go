package zap

type Logger struct{}
type SugaredLogger struct{}

func (l *Logger) Info(msg string, fields ...any) {}
func (l *SugaredLogger) Infow(msg string, keysAndValues ...any) {}
func L() *Logger { return &Logger{} }