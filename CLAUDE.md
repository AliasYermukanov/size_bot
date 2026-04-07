# Size Bot

Telegram bot written in Go. Provides fun commands for group chats.

## Build & Run

```bash
# Build
go build -o size_bot main.go

# Run (requires BOT_TOKEN env var)
BOT_TOKEN=<token> ./size_bot
```

## Deployment

Deployed on **Railway** via Nixpacks. See `railway.json` for config.

## Project Structure

Single-file Go app (`main.go`) using `github.com/go-telegram-bot-api/telegram-bot-api/v5`.

## Bot Commands

- `/start` — greeting (private chats only)
- `/cock_size` — random daily size (cached per user per day, in-memory)
- `/door` — pings group admins

## Notes

- All user state is in-memory (`userDataMap`) — resets on restart
- Messages are in Russian
- Bot token comes from `BOT_TOKEN` environment variable
