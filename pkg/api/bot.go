package bot

import (
	"log"
	"os"

	answer "tgbot/app/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {

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
	processUpdates(bot, updates)
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

func processUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // проверка новго сообщения
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			responce := answer.CreateAnswer(update.Message.Text)

			// создаем ответное сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responce.AnswerMessage)
			msg.ReplyMarkup = buildReplyKeyboard()
			// отправка сообщения
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func buildReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Кнопка 1"),
			tgbotapi.NewKeyboardButton("Кнопка 2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Кнопка 3"),
		),
	)
	return keyboard
}

func buildInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Кнопка 1 ", "Кнопка 1"),
		),
	)
	return keyboard
}
