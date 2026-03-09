# logmsglint

Кастомный линтер для `golangci-lint`, который проверяет сообщения логов в `log/slog` и `go.uber.org/zap`.
---

Проверки:
- сообщение начинается со строчной буквы;
- сообщение содержит только английский текст;
- сообщение не содержит спецсимволы и emoji;
- сообщение не содержит потенциально чувствительные данные.

Требования задания соответствуют тестовому описанию: поддержка `slog`, `zap`, unit-тесты, проверка на другом проекте и интеграция с `golangci-lint`. 
Реализованы обязательные 4 этапа + добавлены CI и зеализация правил
---

## Что проверяется

Примеры нарушений:

```go
slog.Info("Starting server")
slog.Info("запуск сервера")
slog.Info("server started!")
slog.Info("password leaked")
```

Примеры корректных сообщений:

```go
slog.Info("starting server")
slog.Info("server started")
slog.Info("token validated")
```

## Поддерживаемые вызовы

- `slog.Info / Warn / Error / Debug`
- `(*slog.Logger).Info / Warn / Error / Debug`
- `(*zap.Logger).Info / Warn / Error / Debug`
- `(*zap.SugaredLogger).Infow / Warnw / Errorw / Debugw`

Также разбираются:
- строковые константы;
- переменные со строками;
- конкатенация строк;
- `fmt.Sprintf(...)`, если аргументы можно вычислить статически.

---

## Сборка

### 1. Обычный запуск анализатора

```bash
go build -o logmsglint ./cmd/logmsglint
```

### 2. Сборка кастомного golangci-lint

```bash
golangci-lint custom -v
```

После сборки появится `custom-gcl.exe`.

---

## Запуск на другом проекте
проверка проводилась на нескольких проектах один из:https://github.com/ambdasha/Price-Monitor

Пример запуска:

```powershell
PS C:\Users\user\Desktop\pet-proecti\price-monitor> C:\Users\user\Desktop\pet-proecti\selectel_pr\custom-gcl.exe run -c C:\Users\user\Desktop\pet-proecti\price-monitor\.golangci.yml
```

Когда все нормально:

```powershell
PS C:\Users\user\Desktop\pet-proecti\price-monitor> C:\Users\user\Desktop\pet-proecti\selectel_pr\custom-gcl.exe run -c C:\Users\user\Desktop\pet-proecti\price-monitor\.golangci.yml
0 issues.
```

Если добавить новый некорректный кейс:

```powershell
C:\Users\user\Desktop\pet-proecti\price-monitor\.golangci.yml
main.go:6:12: log message must start with a lowercase letter (logmsglint)
        slog.Info("Starting server")
                  ^
1 issues:
* logmsglint: 1
```
---  
