package botlogic

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"main.go/cmd/internal/config"
	"main.go/cmd/internal/recommendations"
	weather "main.go/cmd/internal/weatherApi"
)

// ProcessUpdates обрабатывает обновления Telegram бота
// Параметры:
// - bot: экземпляр бота
// - updates: канал обновлений
// - cfg: конфигурация приложения
// - log: логгер
func ProcessUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, cfg *config.Config, log *slog.Logger) {
	for update := range updates {
		// Обработка команды /start
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// Отправка приветственного сообщения при команде /start
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот погоды и рекомендаций. Отправьте мне название города, чтобы узнать погоду и получить рекомендации по выбору одежды для выхода на улицу.")
				bot.Send(msg)
				continue
			}
		}

		// Пропуск обработки, если сообщение пустое
		if update.Message == nil {
			continue
		}

		// Получение названия города из текста сообщения
		location := update.Message.Text
		// Получение данных о погоде
		weather, err := weather.GetWeather(cfg.WeatherAPI.Key, cfg.WeatherAPI.BaseURL, location)
		if err != nil {
			// Отправка сообщения об ошибке при получении данных о погоде
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка получения данных о погоде в %s, попробуйте снова.", location))
			bot.Send(msg)
			continue
		}

		// Формирование ответа с информацией о погоде
		response := fmt.Sprintf("Город: %s, %s, %s\nТемпература: %.1f°C\nВлажность: %d%%\nСкорость ветра: %.1f км/ч\nУсловия: %s",
			weather.Location.Name, weather.Location.Region, weather.Location.Country,
			weather.Current.TempC, weather.Current.Humidity, weather.Current.WindKph, weather.Current.Condition.Text)

		// Отправка сообщения с информацией о погоде
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)

		// Отладочный вывод URL GigaChat API
		fmt.Printf("GigaChat BaseURL: %s\n", cfg.GigaChat.BaseURL)

		// Получение рекомендаций на основе погодных данных
		recommendationText, err := recommendations.GetRecommendations(cfg.GigaChat.ApiKey, cfg.GigaChat.BaseURL, weather)
		if err != nil {
			// Отправка сообщения об ошибке при получении рекомендаций
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error getting recommendations: %v", err))
			bot.Send(msg)
			continue
		}

		// Отправка сообщения с рекомендацией
		recMsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Совет: %s", recommendationText))
		bot.Send(recMsg)
	}
}
