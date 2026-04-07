package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserData struct {
	LastSize    int
	LastDate    string
	LastMessage string
}

var userDataMap = make(map[int64]*UserData)

const elsMessage = "🍆 Твой размер сегодня: @els_15 см \n\nС таким только детей пугать"
const alikMessage = "🍆 Твой размер сегодня: 1 см\n\nу @AliasYermukanov разводов больше было чем у тебя см"
const aliMessage = "🍆 Твой размер сегодня: 2 см\n\nМеньше только у @StylebenderAli"
const akimMessage = "🍆 Твой размер сегодня: 3 см\n\nПрям как трезубец у @nicekz"

//todo add Akim message

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

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
				if _, err := bot.Send(msg); err != nil {
					log.Printf("failed to send message: %v", err)
				}
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
	if _, err := bot.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
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
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
		return
	}

	size := rand.Intn(25) + 1

	userData.LastSize = size
	userData.LastDate = today
	userData.LastMessage = formatSizeMessage(size)

	msg := tgbotapi.NewMessage(chatID, userData.LastMessage)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func handleDoor(bot *tgbotapi.BotAPI, chatID int64, message *tgbotapi.Message) {
	if message.Chat.Type == "private" {
		msg := tgbotapi.NewMessage(chatID, "🚫 Команда работает только в группах")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
		return
	}

	admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	})
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "🚪 откройте дверь")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
		return
	}

	text := "🚪 откройте дверь\n\n"
	offset := len(utf16.Encode([]rune(text)))
	var entities []tgbotapi.MessageEntity

	for _, admin := range admins {
		if admin.User.ID == bot.Self.ID {
			continue
		}

		if admin.User.UserName != "" {
			mention := "@" + admin.User.UserName + " "
			text += mention
			offset += len(utf16.Encode([]rune(mention)))
		} else {
			name := admin.User.FirstName
			if admin.User.LastName != "" {
				name += " " + admin.User.LastName
			}

			nameLen := len(utf16.Encode([]rune(name)))
			text += name + " "
			entities = append(entities, tgbotapi.MessageEntity{
				Type:   "text_mention",
				Offset: offset,
				Length: nameLen,
				User:   admin.User,
			})
			offset += nameLen + 1
		}
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.Entities = entities
	if _, err := bot.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

// ================= HELPERS =================

func formatSizeMessage(size int) string {
	switch size {
	case 15:
		return elsMessage
	case 1:
		return alikMessage
	case 2:
		return aliMessage
	case 3:
		return akimMessage
	}

	messages := getSizeMessages(size)
	return fmt.Sprintf("🍆 Твой размер сегодня: %d см\n\n%s", size, messages[rand.Intn(len(messages))])
}

func getSizeMessages(size int) []string {
	switch {
	case size <= 5:
		return []string{
			"@tynezloi у тебя просто член маленький",
			"Размер как у детской сосиски из Магнума 😭",
			"Можно считать инвалидностью, но у @Karama_magan еще меньше",
		}

	case size <= 10:
		return []string{
			"Бюджетный вариант, но рабочий 🤣 жаль не как @xRyden и @rchum",
			"С таким даже твое имя не запомнят прям как @adilnrglm",
			"Рыбный четверг у @Aibek09 обеспечен",
			"Как раз такие любит @Dekirr",
		}

	case size <= 15:
		return []string{
			"Чеее происходит @StylebenderAli уже возбудился",
			"@StylebenderAli будет доказывать что у него больше, жаль не уточнит что в жопе",
			"Рабочий кабанчик, уважаю 🚀",
		}

	case size <= 20:
		return []string{
			"Размерчик до сломанного колена @els_15",
			"Даже @azamat_doc за одну хапку такой не проглотит",
			"Вызывайте скорую у @aynkrmv сейчас сердце остановится",
			"Любой Ктлщик мечтает о таком во рту @SSMuchacho не даст соврать",
		}

	default: // 21–25
		return []string{
			"С таким точно не останешься без работы как @xRyden",
			"Сантиметров больше чем волос у @rchum",
			"@nice_kz cock bro",
		}
	}
}
