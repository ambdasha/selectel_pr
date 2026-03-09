package basic

import (
	"fmt"
	"log/slog"

	"go.uber.org/zap"
)

type customLogger struct{}

func (customLogger) Info(msg string) {}

func testSlogPackage() {
	slog.Info("Starting server")   // want `log message must start with a lowercase letter`
	slog.Info("запуск сервера")    // want `log message must contain only English letters`
	slog.Info("server started!")   // want `log message must not contain special symbols or emoji`
	slog.Info("password leaked")   // want `log message must not contain potentially sensitive data`
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

func testConst() {
	const msg = "Starting server"
	slog.Info(msg) // want `log message must start with a lowercase letter`
}

func testVar() {
	var msg = "password leaked"
	slog.Info(msg) // want `log message must not contain potentially sensitive data`
}

func testAssign() {
	msg := "server started!"
	slog.Info(msg) // want `log message must not contain special symbols or emoji`
}

func testConcat() {
	msg := "Starting " + "server"
	slog.Info(msg) // want `log message must start with a lowercase letter`
}

func testSprintf() {
	slog.Info(fmt.Sprintf("Starting %s", "server")) // want `log message must start with a lowercase letter`
	slog.Info(fmt.Sprintf("password %s", "leaked")) // want `log message must not contain potentially sensitive data`
}

func testCustom(c customLogger) {
	c.Info("Starting server")
}