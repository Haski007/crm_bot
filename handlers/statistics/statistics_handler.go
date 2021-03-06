package statistics

import (
	"fmt"
	"time"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"
	"../../utils"
	"../users"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


type m bson.M

type purchase struct {
	prodName string
	prodType string
	amount float64
	unit string
	cash float64
	profit float64
	seller string
	saleDate time.Time
	ID string
}

func GetStatisticsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
	update.CallbackQuery.Message.MessageID,
	emoji.GraphicIncrease + " *Стастистика* " + emoji.GraphicIncrease,
	keyboards.StatsKeyboard)

	answer.ParseMode = "MarkDown"

	bot.Send(answer)
}

func GetCurrentDayHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var products []betypes.Product
	fromDate := utils.GetTodayStartTime()

	query := m{
		"purchases": m{
			"$elemMatch": m{
				"sale_date": m{
					"$gt": fromDate,
				},
			},
		},
	}

	err := database.ProductsCollection.Find(query).All(&products)
	if err != nil {
		bot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			emoji.Warning + emoji.Warning + emoji.Warning + "ERROR: "+err.Error(),
			keyboards.MainMenu))
	}

	var message string

	var purchases []purchase
	for _, prod := range products {
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			purchases = append(purchases, purchase{
				prod.Name, prod.Type, prod.Purchases[i].Amount, prod.Unit,
				prod.Purchases[i].Amount*prod.Price,
				prod.Purchases[i].Amount*prod.Price - prod.Purchases[i].Amount*prod.PrimeCost,
				prod.Purchases[i].Seller,
				prod.Purchases[i].SaleDate.In(utils.Location),
				prod.Purchases[i].ID.String(),
			})
			i--
		}
	}

	// ---> Sort purchases by date.
	for i := len(purchases); i > 0; i-- {
		for j := 1; j < i; j++ {
			if purchases[j - 1].saleDate.After(purchases[j].saleDate) {
				purchases[j - 1], purchases[j] = purchases[j], purchases[j - 1]
			}
		}		
	}
		
	// ---> Build list of sorted purchases
	var index int = 1
	for _, pur := range purchases {
		message += fmt.Sprintf("%sПокупка #%d\nПродукт: %s\nТип: %s\nПродано: %v %s\nДенег получено: %.2f\nЧистыми: %.2f\nПродавец: %s\nДата продажи: %v\n%s\n",
		emoji.GreenDelimiter, index, pur.prodName, pur.prodType, pur.amount, pur.unit,
		pur.cash,
		pur.profit,
		pur.seller,
		pur.saleDate.Format("02.01.2006 15:04:05"),
		pur.ID)
		index++
	}

	var answer tgbotapi.EditMessageTextConfig
	if message != "" {
		answer = tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			message, keyboards.HistoryKeyboard)
		// answer.ParseMode = "MarkDown"
	} else {
		answer = tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			emoji.Warning + " Сегодня пока небыло продаж! " + emoji.Warning,
			keyboards.MainMenu)
	}

	bot.Send(answer)
}

func GetCurrentDayStatsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsAdmin(update.CallbackQuery.From) {
		answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			emoji.NoEntry + "У вас недостаточно прав!" + emoji.NoEntry,
			keyboards.MainMenu) 
		bot.Send(answer)
		return
	}
	message := getDailyStatistics()

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message, keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}