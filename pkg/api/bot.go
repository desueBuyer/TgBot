package api

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	log.Printf("Аторизован как %s", bot.Self.UserName)

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

			var msg tgbotapi.MessageConfig

			switch update.Message.Text {
			default:
				{
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
					msg.ReplyMarkup = buildReplyKeyboard()
				}
			case "Погода":
				{
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город")
					msg.ReplyMarkup = buildInlineKeyboardWeatherCity()
				}
			}
			// отправка сообщения
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else {
			callbackData := update.CallbackQuery.Data
			arr := strings.Split(callbackData, ":")
			switch arr[0] {
			case "weatherCity":
				{
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите дату")
					msg.ReplyMarkup = buildInlineKeyboardWeatherDate(callbackData)
					// отправка сообщения
					deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					bot.Send(deleteMsg)
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			case "weatherDate":
				{
					arr = strings.Split(callbackData, ";")
					res := GetWeatherforecast(strings.Split(arr[1], ":")[1], strings.Split(arr[0], ":")[1])
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, res)
					msg.ReplyMarkup = nil
					deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					bot.Send(deleteMsg)
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			default:
				{
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неизвестная команда")
					msg.ReplyMarkup = buildReplyKeyboard()
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			}
		}
	}
}

func buildReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Погода"),
			tgbotapi.NewKeyboardButton("Кнопка 2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Кнопка 3"),
		),
	)
	return keyboard
}

func buildInlineKeyboardWeatherCity() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Москва", "weatherCity:Moscow"),
		),
	)
	return keyboard
}

func buildInlineKeyboardWeatherDate(city string) tgbotapi.InlineKeyboardMarkup {

	now := time.Now()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(now.Format("02.01.2006"), fmt.Sprintf("weatherDate:%s;%s", now.Format("2006-01-02"), city)),
			tgbotapi.NewInlineKeyboardButtonData(now.AddDate(0, 0, 1).Format("02.01.2006"), fmt.Sprintf("weatherDate:%s;%s", now.AddDate(0, 0, 1).Format("2006-01-02"), city)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(now.AddDate(0, 0, 2).Format("02.01.2006"), fmt.Sprintf("weatherDate:%s;%s", now.AddDate(0, 0, 2).Format("2006-01-02"), city)),
		),
	)
	return keyboard
}
