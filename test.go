package test

import (
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	slog.Info("  Expired password 123 🚀 ой")

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Token и password expired!")

}
