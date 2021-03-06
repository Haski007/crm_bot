package statistics

import (
	"fmt"
	"time"

	"../../betypes"
	"../../database"
	"../../utils"
	"../../emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


func InitEveryDayStatistics(bot *tgbotapi.BotAPI) {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 23, 0, 0, 0, t.Location())
	d := n.Sub(t)

	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}

	for {
		time.Sleep(d)
		d = (24 * time.Hour)
		go utils.SendInfoToAdmins(bot, getDailyStatistics())
	}

}

func getDailyStatistics() string {

	var products []betypes.Product

	fromDate := utils.GetTodayStartTime()

	database.ProductsCollection.Find(nil).Sort("type").All(&products)

	var totalSum float64
	var totalMoney float64

	var message string = "  "

	var tmp string
	for index, prod := range products {

		if prod.Type != tmp {
			tmp = prod.Type
			message += "*" + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint +
			tmp + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + "*\n"
		}

		amount := 0.0
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			amount += prod.Purchases[i].Amount
			totalSum += prod.Purchases[i].Amount * prod.Price
			totalMoney += prod.Purchases[i].Amount * prod.PrimeCost
			i--
		}
		message += fmt.Sprintf("%02d) *%s*   продано: %.2f %s *(%.2f UAH)*\n", index + 1, prod.Name, amount, prod.Unit, amount * prod.Price)
	}

	message += fmt.Sprintf("Total cash: *%.2f UAH*\nTotal profit: *%.2f UAH*", totalSum, totalSum - totalMoney)

	return message
}