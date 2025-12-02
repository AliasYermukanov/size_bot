package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserData struct {
	LastSize    int
	LastDate    string
	LastMessage string
}

var userDataMap = make(map[int64]*UserData)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		switch update.Message.Command() {
		case "start":
			handleStart(bot, chatID)
		case "cock_size":
			handleCockSize(bot, chatID, userID)
		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Используй /start или /cock_size")
			bot.Send(msg)
		}
	}
}

func handleStart(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Привет! Используй /cock_size чтобы узнать свой размер на сегодня."
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func handleCockSize(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	today := time.Now().Format("2006-01-02")

	userData, exists := userDataMap[userID]
	if !exists {
		userData = &UserData{}
		userDataMap[userID] = userData
	}

	// Check if user already got their size today
	if userData.LastDate == today {
		message := userData.LastMessage
		msg := tgbotapi.NewMessage(chatID, message)
		bot.Send(msg)
		return
	}

	// Generate new size for today
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(25) + 1 // Random number from 1 to 25

	userData.LastSize = size
	userData.LastDate = today
	userData.LastMessage = formatSizeMessage(size)

	msg := tgbotapi.NewMessage(chatID, userData.LastMessage)
	bot.Send(msg)
}

func formatSizeMessage(size int) string {
	messages := getSizeMessages(size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMessage := messages[r.Intn(len(messages))]
	return fmt.Sprintf("Твой размер сегодня: %d см\n\n%s", size, randomMessage)
}

func getSizeMessages(size int) []string {
	switch {
	case size <= 5:
		return []string{
			"Не пазорься! Размер не главное 😏",
			"Компактно, но со вкусом! 🎯",
			"Маленький, да удаленький! 💪",
			"Не переживай, главное - техника! 🎪",
			"Скромно, но достойно! 😎",
			"Качество важнее количества! ✨",
		}
	case size <= 10:
		return []string{
			"Неплохо, неплохо! 👍",
			"Средний класс, ничего так! 😊",
			"Стабильно и надежно! 📊",
			"Нормальный размер для нормального парня! 👌",
			"Работает как часы! ⏰",
		}
	case size <= 15:
		return []string{
			"Вот это да! Неплохо так! 🔥",
			"Солидный размер, уважаю! 💯",
			"Выше среднего, молодец! 🚀",
			"Хороший результат! Продолжай в том же духе! 💪",
			"Достойно! Можно гордиться! 😎",
		}
	case size <= 20:
		return []string{
			"Вау! Это уже серьезно! 🔥🔥",
			"Вот это мощь! Респект! 💪💪",
			"Импозантно! Девчонки будут в восторге! 😏",
			"Солидный размер! Можно хвастаться! 🎉",
			"Выдающийся результат! Поздравляю! 🏆",
		}
	default: // 21-25
		return []string{
			"БОЖЕ МОЙ! Это легендарно! 🔥🔥🔥",
			"МОНСТР! Просто невероятно! 💀",
			"АБСОЛЮТНЫЙ ЧЕМПИОН! Все в шоке! 🏆🏆🏆",
			"МИФИЧЕСКИЙ РАЗМЕР! Таких единицы! ⚡",
			"НЕВЕРОЯТНО! Ты настоящий гигант! 🦣",
			"ЭПИЧНО! Это войдет в историю! 📜",
		}
	}
}
