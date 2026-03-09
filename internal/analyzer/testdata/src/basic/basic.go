package basic

import (
	"log/slog"

	"go.uber.org/zap"
)

type customLogger struct{}

func (customLogger) Info(msg string) {}

func testSlog() {
	slog.Info("Starting server")  // want `log message must start with a lowercase letter`
	slog.Info("server started!")  // want `log message must not contain special symbols or emoji`
	slog.Info("password leaked")  // want `log message must not contain potentially sensitive data`
	slog.Info("starting server")
}

func testSlogLogger(logger *slog.Logger) {
	logger.Info("Starting server") // want `log message must start with a lowercase letter`
	logger.Info("starting server")
}

func testZap(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("Starting server")         // want `log message must start with a lowercase letter`
	sugar.Infow("password leaked", "k", 1) // want `log message must not contain potentially sensitive data`
	zap.L().Info("starting server")
}

func testCustom(c customLogger) {
	c.Info("Starting server")
}