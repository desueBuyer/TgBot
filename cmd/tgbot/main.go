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
	bot := initBot(botToken)

	//устанавливаем режим отладки
	bot.Debug = true

	log.Print("Аторизован как %s", bot.Self.UserName)

	//канал для получения обновлений
	u := initUpdatesConfig()

	updates := bot.GetUpdatesChan(u)

	//обработка входящих сообщений
	proseccUpdates(bot, updates)
}

func initBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	return bot
}

func initUpdatesConfig() tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return u
}

func proseccUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
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
