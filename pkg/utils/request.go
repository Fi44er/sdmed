package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RequestOptions struct {
	Method   string            // Тип запроса (GET, POST, PUT, DELETE и т.д.)
	URL      string            // URL запроса
	Query    map[string]string // Query-параметры
	Headers  map[string]string // Заголовки запроса
	Body     interface{}       // Тело запроса (может быть nil)
	FormData map[string]string // Данные формы (для form-data)
	Files    map[string]string // Файлы для загрузки (ключ - имя поля формы, значение - путь к файлу)
}

// MakeRequest выполняет HTTP-запрос и возвращает ответ
func MakeRequest(options RequestOptions) ([]byte, error) {
	// Добавляем query-параметры к URL, если они есть
	if len(options.Query) > 0 {
		queryParams := url.Values{}
		for key, value := range options.Query {
			queryParams.Add(key, value)
		}
		options.URL += "?" + queryParams.Encode()
	}

	var bodyReader io.Reader
	var contentType string

	// Если есть файлы, создаем multipart/form-data запрос
	if len(options.Files) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Добавляем form-data
		for key, value := range options.FormData {
			err := writer.WriteField(key, value)
			if err != nil {
				return nil, fmt.Errorf("ошибка при добавлении form-data: %v", err)
			}
		}

		// Добавляем файлы
		for fieldName, filePath := range options.Files {
			file, err := os.Open(filePath)
			if err != nil {
				return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
			}
			defer file.Close()

			part, err := writer.CreateFormFile(fieldName, filePath)
			if err != nil {
				return nil, fmt.Errorf("ошибка при создании части файла: %v", err)
			}

			_, err = io.Copy(part, file)
			if err != nil {
				return nil, fmt.Errorf("ошибка при копировании файла: %v", err)
			}
		}

		// Закрываем writer для завершения формирования multipart запроса
		err := writer.Close()
		if err != nil {
			return nil, fmt.Errorf("ошибка при закрытии writer: %v", err)
		}

		bodyReader = body
		contentType = writer.FormDataContentType()
	} else if len(options.FormData) > 0 {
		// Если есть form-data, но нет файлов, отправляем как application/x-www-form-urlencoded
		formData := url.Values{}
		for key, value := range options.FormData {
			formData.Add(key, value)
		}
		bodyReader = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if options.Body != nil {
		// Преобразуем тело запроса в JSON, если оно есть
		bodyBytes, err := json.Marshal(options.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации тела запроса: %v", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
		contentType = "application/json"
	}

	// Создаем новый запрос
	req, err := http.NewRequest(options.Method, options.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	// Добавляем заголовки
	for key, value := range options.Headers {
		req.Header.Add(key, value)
	}

	// Устанавливаем Content-Type, если он не был установлен вручную
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ошибка сервера: %s, тело ответа: %s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}
