package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	Env         string            `yaml:"env" env:"ENV" env-default:"local"`
	WeatherAPI  WeatherAPIConfig  `yaml:"weather_api"`
	TelegramBot TelegramBotConfig `yaml:"telegram_bot"`
	GigaChat    GigachatConfig    `yaml:"gigachat"`
}

// WeatherAPIConfig содержит конфигурацию для доступа к WeatherAPI
type WeatherAPIConfig struct {
	Key     string `yaml:"key" env:"WEATHER_API_KEY"`
	BaseURL string `yaml:"base_url" env-default:"http://api.weatherapi.com/v1/current.json"`
}

// TelegramBotConfig содержит конфигурацию для доступа к Telegram Bot API
type TelegramBotConfig struct {
	TelegramBotToken string `yaml:"telegram_bot_token"`
}

// GigachatConfig содержит конфигурацию для доступа к GigaChat API
type GigachatConfig struct {
	ApiKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
}

// MustLoad загружает конфигурацию из файла и переменных окружения
func MustLoad() *Config {
	// Установка пути к конфигурационному файлу через переменную окружения
	configPath := os.Getenv("CONFIG_PATH_WEATHER")
	if configPath == "" {
		log.Fatal("CONFIG_PATH_WEATHER is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Чтение конфигурации из файла
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
