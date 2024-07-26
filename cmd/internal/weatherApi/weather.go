package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"main.go/cmd/internal/model"
)

// GetWeather получает данные о погоде от WeatherAPI
// Параметры:
// - apiKey: API ключ для доступа к WeatherAPI
// - baseURL: базовый URL WeatherAPI
// - location: название местоположения для получения данных о погоде
// Возвращает:
// - *model.WeatherResponse: структура с данными о погоде
// - error: ошибка, если она произошла
func GetWeather(apiKey, baseURL, location string) (*model.WeatherResponse, error) {
	// Кодируем параметр location для использования в URL
	encodedLocation := url.QueryEscape(location)
	// Формируем URL запроса с параметрами apiKey и location
	url := fmt.Sprintf("%s?key=%s&q=%s", baseURL, apiKey, encodedLocation)

	// Создаем HTTP клиент с таймаутом 10 секунд
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Выполняем GET запрос к WeatherAPI
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодируем JSON ответ в структуру WeatherResponse
	var weatherResponse model.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	// Возвращаем полученные данные о погоде
	return &weatherResponse, nil
}
