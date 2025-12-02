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

		// Skip if message is not a command
		if !update.Message.IsCommand() {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		// Get command without bot username (handles both /command and /command@botname)
		command := update.Message.Command()

		switch command {
		case "start":
			handleStart(bot, chatID)
		case "cock_size":
			handleCockSize(bot, chatID, userID)
		default:
			// Only respond to unknown commands in private chats
			if update.Message.Chat.Type == "private" {
				msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Используй /start или /cock_size")
				bot.Send(msg)
			}
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
			"Брат… это не кок, это USB Type-C разъём 😭",
			"Это не член, это насмешка природы 🤣",
			"Такой только муравьёв пугать, не людей 💀",
			"Печенье «топлёное молоко» и то длиннее 🤣",
			"Не переживай… главное — харизма 😭",
			"Размер как у детской сосиски из Магнума 😭",
		}

	case size <= 10:
		return []string{
			"Бюджетный вариант, но рабочий 🤣",
			"Средний класс — эконом, но уверенный 😎",
			"Нормас, по СНГ-стандарту проходишь 💪",
			"С таким хоть не стыдно в душ заходить 😭",
			"Нормальный кок, рабочая лошадка 😂",
		}

	case size <= 15:
		return []string{
			"Вот это уже техника! Девки хлопают стоя 🔥",
			"Уверенный среднячок, даже гордиться не стыдно 😎",
			"С таким можно говорить «у меня нормальный» без смеха 😭",
			"Рабочий кабанчик, уважаю 🚀",
			"Солидно. Можно хвастаться в чате 😏",
		}

	case size <= 20:
		return []string{
			"Вау. Тут уже тяжело жить с джинсами 🤣",
			"Это уже оружие массового развлечения 🔥🔥",
			"С таким тебе надо паспорт на член оформлять 😭",
			"У тебя там не кок — у тебя DLC к телу 😎",
			"Импозантно. Модно. Молодёжно. Опасно. 💀",
		}

	default: // 21–25
		return []string{
			"ЭТО НЕ ЧЛЕН. ЭТО ЛЕГЕНДА. 💀🔥",
			"Гигант. Монстр. Финальный босс Pornhub'а 😈",
			"Такой только в музее хранить… или в документах Marvel 🦣",
			"Абсолютный чемпион. Остальным стыдно рядом стоять 🏆",
			"С таким даже дверь открывать можно — ручка не нужна 😂",
			"Эпично. Бог дал, чтобы ты страдал в джинсах 😭",
		}
	}
}
