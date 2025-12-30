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
	log.Printf("Authorized as @%s", bot.Self.UserName)

	// 🔥 регистрируем команды
	registerCommands(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		command := update.Message.Command()

		switch command {

		case "start":
			if update.Message.Chat.Type == "private" {
				handleStart(bot, chatID)
			}

		case "cock_size":
			handleCockSize(bot, chatID, userID)

		case "door":
			handleDoor(bot, chatID, update.Message)

		default:
			if update.Message.Chat.Type == "private" {
				msg := tgbotapi.NewMessage(chatID, "Неизвестная команда 🤨")
				bot.Send(msg)
			}
		}
	}
}

// ================= COMMAND REGISTRATION =================

func registerCommands(bot *tgbotapi.BotAPI) {
	// 📌 Команды для групп
	groupCommands := []tgbotapi.BotCommand{
		{Command: "door", Description: "🚪 Открыть дверь"},
		{Command: "cock_size", Description: "🍆 Узнать размер на сегодня"},
	}

	groupCfg := tgbotapi.NewSetMyCommands(groupCommands...)
	groupCfg.Scope = &tgbotapi.BotCommandScope{
		Type: "all_group_chats",
	}
	bot.Request(groupCfg)

	// 📌 Команды для лички
	privateCommands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать"},
	}

	privateCfg := tgbotapi.NewSetMyCommands(privateCommands...)
	privateCfg.Scope = &tgbotapi.BotCommandScope{
		Type: "all_private_chats",
	}
	bot.Request(privateCfg)
}

// ================= HANDLERS =================

func handleStart(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		"Привет 👋\n\n"+
			"Добавь меня в группу и используй:\n"+
			"/cock_size — узнать размер 🍆\n"+
			"/door — открыть дверь 🚪",
	)
	bot.Send(msg)
}

func handleCockSize(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	today := time.Now().Format("2006-01-02")

	userData, exists := userDataMap[userID]
	if !exists {
		userData = &UserData{}
		userDataMap[userID] = userData
	}

	if userData.LastDate == today {
		msg := tgbotapi.NewMessage(chatID, userData.LastMessage)
		bot.Send(msg)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(25) + 1

	userData.LastSize = size
	userData.LastDate = today
	userData.LastMessage = formatSizeMessage(size)

	msg := tgbotapi.NewMessage(chatID, userData.LastMessage)
	bot.Send(msg)
}

func handleDoor(bot *tgbotapi.BotAPI, chatID int64, message *tgbotapi.Message) {
	if message.Chat.Type == "private" {
		msg := tgbotapi.NewMessage(chatID, "🚫 Команда работает только в группах")
		bot.Send(msg)
		return
	}

	admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	})
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "🚪 откройте дверь")
		bot.Send(msg)
		return
	}

	text := "🚪 откройте дверь\n\n"
	offset := len(text)
	var entities []tgbotapi.MessageEntity

	for _, admin := range admins {
		if admin.User.ID == bot.Self.ID {
			continue
		}

		if admin.User.UserName != "" {
			mention := "@" + admin.User.UserName + " "
			text += mention
			offset += len(mention)
		} else {
			name := admin.User.FirstName
			if admin.User.LastName != "" {
				name += " " + admin.User.LastName
			}

			text += name + " "
			entities = append(entities, tgbotapi.MessageEntity{
				Type:   "text_mention",
				Offset: offset,
				Length: len(name),
				User:   admin.User,
			})
			offset += len(name) + 1
		}
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.Entities = entities
	bot.Send(msg)
}

// ================= HELPERS =================

func formatSizeMessage(size int) string {
	messages := getSizeMessages(size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("🍆 Твой размер сегодня: %d см\n\n%s", size, messages[r.Intn(len(messages))])
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
