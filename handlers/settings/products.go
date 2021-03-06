package settings

import (
	"fmt"
	"strconv"
	"strings"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

// Queue of users who are trying to add new product
var (
	AddProductQueue = make(map[int]*betypes.Product)
)

// RemoveProductHandler removes product from "products" collection
func RemoveProductHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	prodID := strings.Split(update.CallbackQuery.Data, " ")[1]

	err := database.ProductsCollection.Remove(bson.M{"_id": bson.ObjectIdHex(prodID)})
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, emoji.Warning+" Ошибка при удалении продукта! {"+err.Error()+"}"))
		return
	}


	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Продукт был успешно удалён! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu

	bot.Send(answer)
}

// AddProductHandler adds product to database collection "products"
func AddProductHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	AddProductQueue[update.CallbackQuery.From.ID] = new(betypes.Product)

	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите название продукта:"))
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

// AddProduct prompt user to get name, type and prise of product. Save it in DB
func AddProduct(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID

	prod := AddProductQueue[update.Message.From.ID]

	var err error
	if prod.Name == "" {
		prod.Name = update.Message.Text

		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите тип продукта:")
		
		typesKeyboard := keyboards.GetTypesKeyboard("protyp")
		typesKeyboard.InlineKeyboard = append(typesKeyboard.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{keyboards.MainMenuButton})
		answer.ReplyMarkup = typesKeyboard

		bot.Send(answer)
	} else if prod.Unit == "" {
		prod.Unit = update.Message.Text

		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите *себестоимость* продукта")
		answer.ParseMode = "MarkDown"
		bot.Send(answer)
	} else if prod.PrimeCost == 0.0 {
		prod.PrimeCost, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный тип данных! Попробуйте ещё раз:"))
			return
		}

		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите *цену продажи*:")
		answer.ParseMode = "MarkDown"
		bot.Send(answer)
		
	} else {
		prod.Price, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный тип данных! Попробуйте ещё раз:"))
			return
		}

		var answer tgbotapi.MessageConfig

		prod.Margin = (prod.Price - prod.PrimeCost) / prod.Price * 100
		err = database.ProductsCollection.Insert(prod)
		if err != nil {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning+" Не удалось создать продукт! {"+err.Error()+"}")
		} else {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Продукт был успешно создан! " + emoji.Check)
		}

		delete(AddProductQueue, userID)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
	}
}

func AddTypeToProduct(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	t := strings.Join(strings.Fields(update.CallbackQuery.Data)[1:], " ")

	AddProductQueue[update.CallbackQuery.From.ID].Type = t
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Укажите единицу измерения:")
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}


func GetAllProductsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите тип продукта:")
		
	typesKeyboard := keyboards.GetTypesKeyboard("getprodstyp")
	typesKeyboard.InlineKeyboard = append(typesKeyboard.InlineKeyboard,
		[]tgbotapi.InlineKeyboardButton{keyboards.MainMenuButton})
	answer.ReplyMarkup = typesKeyboard

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	bot.Send(answer)
}

// GetAllProducts prints all produtcs from "products" collection.
func GetAllProducts(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	var products []betypes.Product
	var t = strings.Join(strings.Fields(update.CallbackQuery.Data)[1:], " ")

	database.ProductsCollection.Find(bson.M{"type":t}).Select(m{"purchases": 0}).Sort("name").All(&products)

	for i, prod := range products {
		message := fmt.Sprintf("======================\nПродукт #%d\nНазвание: *%s*\nТип: *%s*\n"+
		"Себестоимость: *%.2f*\nЦена продажи: *%.2f*\nМаржа: *%.2f %%*\nЕдиница измерения: *%s\n*",
			i, prod.Name, prod.Type, prod.PrimeCost, prod.Price, prod.Margin, prod.Unit)
		answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
		answer.ParseMode = "Markdown"
		prodKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Удалить "+emoji.Minus, "remove_product "+prod.ID.Hex()),
				tgbotapi.NewInlineKeyboardButtonData("Изменить "+emoji.Pencil, "edit_product "+prod.ID.Hex()),
			),
		)
		answer.ReplyMarkup = prodKeyboard
		bot.Send(answer)
	}

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вот и все продукты!")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}