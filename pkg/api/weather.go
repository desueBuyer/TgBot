package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	utils "tgbot/internal/utils"
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
		log.Fatal("Ошибка при выполнении запроса: ", err)
	}
	defer resp.Body.Close()

	//чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при чтении ответа: ", err)
	}

	res := utils.ParceWeatherJson(string(body))

	return res
}
