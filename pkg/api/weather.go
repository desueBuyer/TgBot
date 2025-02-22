package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetWeatherforecast(city string, date string) string {

	variant, existance := os.LookupEnv("WEATHER_API_KEY")

	if !existance {
		log.Fatal("Не найден API ключ для получения прогноза погоды")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?q=%s&days=2&dt=%s&key=%s", city, date, variant)

	//выполнение запроса
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса: ", err)
	}
	defer resp.Body.Close()

	//чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа: ", err)
	}

	return string(body)
}
