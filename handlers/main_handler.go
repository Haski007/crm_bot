package handlers

import (
	"../keyboards"
	"../emoji"
	"./settings"
	"./purchases"
	"./statistics"
	"./store"
	"./cashbox"
	"./users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func MainMenuHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	deleteAllQueues(update.CallbackQuery.From.ID)

	message := "........."+emoji.House+"......."+emoji.Tree+"..*Главное меню*........"+
		emoji.HouseWithGarden+"..."+emoji.Car+"...."
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message, keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}

func deleteAllQueues(id int) {
	delete(settings.AddProductQueue, id)
	delete(settings.AddTypeQueue, id)
	delete(settings.RemoveTypeQueue, id)
	delete(settings.EditProductQueue, id)
	delete(settings.EditProductNameQueue, id)
	delete(settings.EditProductPrimeQueue, id)
	delete(settings.EditProductPriceQueue, id)
	delete(settings.EditProductUnitQueue, id)
	delete(statistics.MonthStatsQueue, id)
	delete(purchases.MakePurchaseQueue, id)
	delete(purchases.RemovePurchaseQueue, id)
	delete(cashbox.PlusCashQueue, id)
	delete(cashbox.MinusCashQueue, id)
	delete(cashbox.TransactionsHostoryQueue, id)
	delete(cashbox.SetStartDailyMoneyQueue, id)
	delete(cashbox.EndDayQueue, id)
	delete(store.SupplyQueue, id)
	delete(users.RegisterUserQueue, id)
}