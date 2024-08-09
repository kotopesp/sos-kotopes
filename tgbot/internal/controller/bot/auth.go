package bot

import (
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (r *Router) start(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID
	err := menu(bot, chatID)
	if err != nil {
		log.Println(err)
		return
	}
}

func (r *Router) back(bot *telego.Bot, update telego.Update) {
	chatID := update.CallbackQuery.Message.GetChat().ID
	err := menu(bot, chatID)
	if err != nil {
		log.Println(err)
		return
	}

	bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func (r *Router) authorize(bot *telego.Bot, update telego.Update) {
	userID := int(update.CallbackQuery.From.ID)
	url, err := r.authService.GetAuthURL(userID)
	if err != nil {
		log.Println(err)
		return
	}

	chatID := update.CallbackQuery.Message.GetChat().ID
	_, err = bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(chatID),
		Text:   "При нажатии на кнопку будет автоматически создан аккаунт и выполнен вход.",
		ReplyMarkup: tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Авторизоваться").WithURL(*url),
				tu.InlineKeyboardButton("Назад").WithCallbackData("back"),
			),
		),
	})
	if err != nil {
		log.Println(err)
		return
	}

	bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}
