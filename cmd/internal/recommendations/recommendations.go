package recommendations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"main.go/cmd/internal/model"
)

// GigaChatRequest представляет структуру запроса к GigaChat API
type GigaChatRequest struct {
	Model             string            `json:"model"`
	Messages          []GigaChatMessage `json:"messages"`
	Stream            bool              `json:"stream"`
	RepetitionPenalty int               `json:"repetition_penalty"`
}

// GigaChatMessage представляет структуру сообщения для GigaChat API
type GigaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GigaChatResponse представляет структуру ответа от GigaChat API
type GigaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GetRecommendations отправляет запрос к GigaChat API и получает рекомендации на основе погодных данных
// Параметры:
// - apiKey: API ключ для доступа к GigaChat
// - baseURL: базовый URL GigaChat API
// - weather: данные о погоде
// Возвращает:
// - string: рекомендация по выбору одежды
// - error: ошибка, если она произошла
func GetRecommendations(apiKey, baseURL string, weather *model.WeatherResponse) (string, error) {
	// Формирование тела запроса
	requestBody := GigaChatRequest{
		Model: "GigaChat",
		Messages: []GigaChatMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Погода за окном: Temperature: %.1f°C, Humidity: %d%%, Wind Speed: %.1f kph, что порекомендуешь надеть на улицу?", weather.Current.TempC, weather.Current.Humidity, weather.Current.WindKph),
			},
		},
		Stream:            false,
		RepetitionPenalty: 1,
	}

	// Сериализация тела запроса в JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	// Создание нового HTTP запроса
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Установка заголовков запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Отладочные выводы
	fmt.Printf("Request URL: %s\n", baseURL)
	fmt.Printf("Request Body: %s\n", string(body))

	// Инициализация HTTP клиента с таймаутом
	client := &http.Client{Timeout: 60 * time.Second}
	// Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Отладочный вывод
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Headers: %v\n", resp.Header)

	// Проверка статус кода ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодирование JSON ответа в структуру GigaChatResponse
	var response GigaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	// Проверка наличия рекомендаций в ответе и возвращение первой рекомендации
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	// Если рекомендаций нет, возвращаем пустую строку
	return "", nil
}
