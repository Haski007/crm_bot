package settings

import (
	"fmt"

	"../../emoji"
	"../../keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SettingsHandler handle "Configuration" callback (button)
func SettingsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// emojiRow := utils.MakeEmojiRow(emoji.Wrench, 12)
	message := fmt.Sprintf("%s\n*SETTINGS*\n%s", emoji.Gear, emoji.Gear)
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, message, keyboards.SettingsKeyboard)
	answer.ParseMode = "MarkDown"

	bot.Send(answer)
}
