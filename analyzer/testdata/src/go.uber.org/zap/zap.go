package zap

type Logger struct{}

func (l *Logger) Info(msg string, fields ...any)   {}
func (l *Logger) Error(msg string, fields ...any)  {}
func (l *Logger) Warn(msg string, fields ...any)   {}
func (l *Logger) Debug(msg string, fields ...any)  {}
func (l *Logger) DPanic(msg string, fields ...any) {}
func (l *Logger) Panic(msg string, fields ...any)  {}
func (l *Logger) Fatal(msg string, fields ...any)  {}

func NewProduction() (*Logger, error) { return &Logger{}, nil }

func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

type SugaredLogger struct{}

func (s *SugaredLogger) Infof(template string, args ...any)     {}
func (s *SugaredLogger) Debugf(template string, args ...any)    {}
func (s *SugaredLogger) Warnf(template string, args ...any)     {}
func (s *SugaredLogger) Errorf(template string, args ...any)    {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...any) {}
