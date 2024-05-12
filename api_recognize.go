package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Done       bool   `json:"done"`
	ID         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	CreatedBy  string `json:"createdBy"`
	ModifiedAt string `json:"modifiedAt"`
}

func recognize() {
	upload_to_storage()

	// Укажите ссылку на ваш бакет и путь к файлу
	fileURL := "https://storage.yandexcloud.net/cipher/mic.ogg"

	// Укажите ваш API-ключ Object Storage
	apiKey := "AQVNzx4lwjTKXXxzNcWOOY9JnqC9TF4OrEtUAn3L"

	// Создание запроса на распознавание текста
	recognizeURL := "https://transcribe.api.cloud.yandex.net/speech/stt/v2/longRunningRecognize"

	// Установка параметров запроса
	bodyData := map[string]interface{}{
		"config": map[string]interface{}{
			"specification": map[string]interface{}{
				"languageCode": "ru-RU",
				"model":        "general",
			},
		},
		"audio": map[string]interface{}{
			"uri": fileURL,
		},
	}

	// Преобразование данных запроса в JSON
	requestBody, err := json.Marshal(bodyData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Отправка запроса на распознавание текста
	req, err := http.NewRequest("POST", recognizeURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Установка заголовков запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+apiKey)

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Десериализация ответа в структуру Response
	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := response.ID
	tt := 0
	for {

		GetURL := "https://operation.api.cloud.yandex.net/operations/" + id

		req, err = http.NewRequest("GET", GetURL, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Authorization", "Api-Key "+apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа
		re, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(re), &data); err != nil {
			panic(err)
		}

		done := data["done"].(bool)

		if done {
			text := []uint8(re)
			write_file("text.json", text)
			parser_json()
			fmt.Println("Done recognize")
			break
		}

		// Ожидаем
		time.Sleep(10 * time.Second)

		tt += 10
		fmt.Println("Time: ", tt, " seconds")
	}
}
