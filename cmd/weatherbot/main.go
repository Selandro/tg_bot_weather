package main

import (
	"fmt"
	"os"

	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"main.go/cmd/internal/botlogic"
	"main.go/cmd/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Инициализация логгера
	log := setupLogger(cfg.Env)
	log.Info(fmt.Sprintf("Starting weather and recommendation bot, env: %s", cfg.Env))

	// Получение токена бота из переменных окружения
	botToken := cfg.TelegramBot.TelegramBotToken
	if botToken == "" {
		log.Error("TELEGRAM_BOT_TOKEN is not set")
		os.Exit(1)
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Error(fmt.Sprintf("Error initializing bot: %v", err))
		os.Exit(1)
	}

	// Включение режима отладки для бота
	bot.Debug = true

	// Логирование успешной авторизации
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получение канала для получения обновлений от бота
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting updates: %v", err))
		os.Exit(1)
	}

	// Обработка обновлений бота
	botlogic.ProcessUpdates(bot, updates, cfg, log)
}

// setupLogger настраивает логгер в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	// Настройка логгера в зависимости от окружения
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
