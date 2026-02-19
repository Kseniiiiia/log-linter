package multi

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// slog examples
	slog.Info("good message")
	slog.Info("Ошибка подключения password! 🚀") // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emoji" "log message contains sensitive data"
	slog.Info("user password invalid")          // want "log message contains sensitive data"
	slog.Info("Bad сообщение!!")                // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emoji"

	// zap logger examples
	logger, _ := zap.NewProduction()
	logger.Info("good message")
	logger.Info("Ошибка подключения password! 🚀") // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emojis" "log message contains sensitive data"
	logger.Error("user password invalid")         // want "log message contains sensitive data"
	logger.Warn("Bad сообщение!!")                // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emojis"

	// zap sugared logger examples
	sugar := logger.Sugar()
	sugar.Infof("good message")
	sugar.Infof("Ошибка подключения password! 🚀") // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emojis" "log message contains sensitive data"
	sugar.Errorf("user password invalid")         // want "log message contains sensitive data"
	sugar.Warnf("Bad сообщение!!")                // want "log message must start with lowercase letter" "log message must contain only English letters" "log message must not contain special characters or emojis"
}
