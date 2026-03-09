package basic

import "log/slog"

func test() {
	slog.Info("Starting server")   
	slog.Info("запуск сервера")    
	slog.Info("server started!👋")   
	slog.Info("password exposed")  

	slog.Info("starting server")
	slog.Info("server started")
}