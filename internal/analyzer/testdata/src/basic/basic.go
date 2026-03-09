package basic

import "log/slog"

func test() {
	
	slog.Info("Starting server")   // want `log message must start with a lowercase letter`
	slog.Info("запуск сервера")    // want `log message must contain only English letters`
	slog.Info("server started👋!")   // want `log message must not contain special symbols or emoji`
	slog.Info("password exposed")  // want `log message must not contain potentially sensitive data`
  

	slog.Info("starting server")
	slog.Info("server started")
}