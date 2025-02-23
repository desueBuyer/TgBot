package utils

import (
	"encoding/json"
	"log"
	"os"
)

// Структуры для парсинга JSON
type SystemVariable struct {
	BotApiKey     string `json:"botApiKey"`
	WeatherApiKey string `json:"weatherApiKey"`
}

type WeatherCity struct {
	Button string `json:"button"`
	Value  string `json:"value"`
}

type BotSettings struct {
	WeatherCities []WeatherCity `json:"weatherCities"`
}

type Config struct {
	AppName         string           `json:"appName"`
	Version         string           `json:"version"`
	Environment     string           `json:"environment"`
	Debug           bool             `json:"debug"`
	SystemVariables []SystemVariable `json:"systemVariables"`
	BotSettings     []BotSettings    `json:"botSettings"`
}

func InitConfig() Config {
	// Парсинг JSON
	var config Config
	rawConfig, err := os.ReadFile("internal\\config\\config.json")
	if err != nil {
		log.Fatal("Не могу прочитать файл конфигурации!")
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatal("Ошибка при парсинге JSON:", err)
	}
	return config
}
