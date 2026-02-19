package sensitive

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// slog examples
	slog.Info("all good")
	slog.Info("user password invalid")           // want "log message contains sensitive data"
	slog.Info("token expired")                   // want "log message contains sensitive data"
	slog.Info("api_key=123")                     // want "log message contains sensitive data"
	slog.Info("token : 12345, password : 12345") // want "log message contains sensitive data" "log message contains sensitive data"

	// zap logger examples
	logger, _ := zap.NewProduction()
	logger.Info("all good")
	logger.Info("user password invalid")            // want "log message contains sensitive data"
	logger.Error("api_key=123")                     // want "log message contains sensitive data"
	logger.Warn("token expired")                    // want "log message contains sensitive data"
	logger.Debug("token : 12345, password : 12345") // want "log message contains sensitive data" "log message contains sensitive data"

	// zap sugared logger examples
	sugar := logger.Sugar()
	sugar.Infof("all good")
	sugar.Errorf("token expired")                  // want "log message contains sensitive data"
	sugar.Warnf("user password invalid")           // want "log message contains sensitive data"
	sugar.Debugf("api_key=123")                    // want "log message contains sensitive data"
	sugar.Warnf("token : 12345, password : 12345") // want "log message contains sensitive data" "log message contains sensitive data"
}
