package english

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// slog examples
	slog.Info("hello world")
	slog.Info("ошибка подключения") // want "log message must contain only English letters"
	slog.Info("mixed ошибка")       // want "log message must contain only English letters"
	slog.Info("русский текст")      // want "log message must contain only English letters"

	// zap logger examples
	logger, _ := zap.NewProduction()
	logger.Info("hello world")
	logger.Info("ошибка подключения") // want "log message must contain only English letters"
	logger.Error("mixed ошибка")      // want "log message must contain only English letters"
	logger.Warn("русский текст")      // want "log message must contain only English letters"

	// zap sugared logger examples
	sugar := logger.Sugar()
	sugar.Infof("hello world")
	sugar.Errorf("ошибка подключения") // want "log message must contain only English letters"
	sugar.Warnf("mixed ошибка")        // want "log message must contain only English letters"
	sugar.Debugf("ошибка подключения") // want "log message must contain only English letters"
}
