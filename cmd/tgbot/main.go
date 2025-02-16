package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	//получае токен бота из переменной окружения
	botToken, exists := os.LookupEnv("BOT_API_KEY")
	if !exists {
		log.Fatal("BOT_API_KEY не найден")
	}

	//создаем новый экземляр бота
	bot := botInitializer(botToken)

	//устанавливаем режим отладки
	bot.Debug = true

	log.Print("Аторизован как %s", bot.Self.UserName)

	//канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//обработка входящих сообщений
	for update := range updates {
		if update.Message != nil { // проверка новго сообщения
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// создаем ответное сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// отправка сообщения
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func botInitializer(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	return bot
}
