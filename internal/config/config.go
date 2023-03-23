package config

import "os"

var BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

const GroupChatID int64 = -865707097
