package api

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	utils "tgbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
	Config  utils.Config
}

func (bot *Bot) initBot(apiKeyEnvVariant string) {
	//получаем токен бота из переменной окружения
	token, exists := os.LookupEnv(apiKeyEnvVariant)
	if !exists {
		log.Fatal("BOT_API_KEY не найден")
	}

	//создаем новый экземляр бота
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Api = api
	log.Printf("Аторизован как %s", bot.Api.Self.UserName)

	//устанавливаем режим отладки
	bot.Api.Debug = true
}

// канал для получения обновлений
func (bot *Bot) initUpdatesChannel() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	bot.Updates = bot.Api.GetUpdatesChan(u)
}

// обработка входящих сообщений
func (bot *Bot) loopThroughUpdates() {
	for update := range bot.Updates {
		if update.Message != nil {
			bot.processMessage(update)
		} else {
			if update.CallbackQuery != nil {
				bot.processCallbackQuery(update)
			}
		}
	}
}

func (bot *Bot) processMessage(update tgbotapi.Update) {
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
			msg.ReplyMarkup = buildInlineKeyboardWeatherCity(bot.Config.BotSettings[0].WeatherCities)
		}
	case "Курсы валют":
		{
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Заглушка")
		}
	}
	// отправка сообщения
	if _, err := bot.Api.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (bot *Bot) processCallbackQuery(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	arr := strings.Split(callbackData, ":")
	switch arr[0] {
	case "weatherCity":
		{
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите дату")
			msg.ReplyMarkup = buildInlineKeyboardWeatherDate(callbackData)
			// отправка сообщения
			deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			bot.Api.Send(deleteMsg)
			if _, err := bot.Api.Send(msg); err != nil {
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
			bot.Api.Send(deleteMsg)
			if _, err := bot.Api.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	default:
		{
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неизвестная команда")
			msg.ReplyMarkup = buildReplyKeyboard()
			if _, err := bot.Api.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func RunBot(config utils.Config) {

	bot := new(Bot)
	bot.Config = config
	bot.initBot("BOT_API_KEY")
	bot.initUpdatesChannel()
	bot.loopThroughUpdates()
}

func buildReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Погода"),
			tgbotapi.NewKeyboardButton("Курсы валют"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Кнопка 3"),
		),
	)
	return keyboard
}

func buildInlineKeyboardWeatherCity(weatherCities []utils.WeatherCity) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, data := range weatherCities {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(data.Button, "weatherCity:"+data.Value)))
	}
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
