Weather and Recommendation Bot

Этот проект представляет собой Telegram-бота для получения информации о погоде
и рекомендаций по выбору одежды в зависимости от погодных условий. Бот использует
API погоды и GigaChat для генерации рекомендаций.

Функциональность:
Получение текущей погоды по запросу пользователя.
Генерация рекомендаций по выбору одежды на основе погодных условий.
Приветственное сообщение и инструкции по использованию бота при вводе команды /start.

Требования:
Go 1.20+
Go Modules
Telegram Bot API Token
Weather API Key
GigaChat API Key

Установка:

git clone https://github.com/selandro/tg_bot_weather.git
cd weather-recommendation-bot
Установите переменную окружения CONFIG_PATH_WEATHER, которая указывает на путь к вашему конфигурационному файлу.

Конфигурация
Создайте файл конфигурации config.yml в корне проекта. Пример конфигурационного файла:
env: local
weather_api:
  key: YOUR_WEATHER_API_KEY
  base_url: "http://api.weatherapi.com/v1/current.json"
telegram_bot:
  telegram_bot_token: YOUR_TELEGRAM_BOT_TOKEN
gigachat:
  api_key: YOUR_GIGACHAT_API_KEY
  base_url: "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

