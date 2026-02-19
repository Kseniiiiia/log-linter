package symbols

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// slog examples
	slog.Info("hello world")
	slog.Info("error!")   // want "log message must not contain special characters or emoji"
	slog.Info("rocket 🚀") // want "log message must not contain special characters or emoji"
	slog.Info("path/to/file")
	slog.Info("status_ok")
	slog.Info("sure?") // want "log message must not contain special characters or emoji"

	// zap logger examples
	logger, _ := zap.NewProduction()
	logger.Info("hello world")
	logger.Info("error!")   // want "log message must not contain special characters or emojis"
	logger.Info("rocket 🚀") // want "log message must not contain special characters or emojis"
	logger.Error("path/to/file")
	logger.Warn("status_ok")
	logger.Debug("sure?") // want "log message must not contain special characters or emojis"

	// zap sugared logger examples
	sugar := logger.Sugar()
	sugar.Infof("hello world")
	sugar.Infof("error! %s", "details") // want "log message must not contain special characters or emojis"
	sugar.Errorf("rocket 🚀")            // want "log message must not contain special characters or emojis"
	sugar.Warnf("path/to/file")
	sugar.Debugf("status_ok")
	sugar.Warnf("sure?") // want "log message must not contain special characters or emojis"
}
