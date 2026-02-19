package lowercase

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// slog examples
	slog.Info("Bad message")   // want "log message must start with lowercase letter"
	slog.Info("  Bad message") // want "log message must start with lowercase letter"
	slog.Info("404 Nоt found")
	slog.Info("good message")

	// zap logger examples
	logger, _ := zap.NewProduction()
	logger.Info("Bad message")    // want "log message must start with lowercase letter"
	logger.Error("  Bad message") // want "log message must start with lowercase letter"
	logger.Debug("404 Nоt found")
	logger.Info("good message")

	// zap sugared logger examples
	sugar := logger.Sugar()
	sugar.Infof("Bad message")    // want "log message must start with lowercase letter"
	sugar.Errorf("  Bad message") // want "log message must start with lowercase letter"
	sugar.Warnf("404 Nоt found")
	sugar.Infof("good message")
}
