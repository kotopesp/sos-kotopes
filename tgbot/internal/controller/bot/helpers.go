package bot

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func menu(bot *telego.Bot, chatID int64) error {
	_, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   "Какой сервис вас интересует?",
		ReplyMarkup: tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Авторизация").WithCallbackData("authorize"),
			),
		),
	})
	return err
}
