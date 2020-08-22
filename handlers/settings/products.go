package settings

import (
	"encoding/json"
	"strconv"
	"strings"

	"../../betypes"
	"../../database"
	"../../keyboards"
	"../../emoji"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

// Queue of users who are trying to add new product
var AddProductQueue = make(map[int]*betypes.Product)

// GetAllProductsHandler prints all produtcs from "products" collection.
func GetAllProductsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {

	var products []betypes.Product

	database.ProductsCollection.Find(bson.M{}).Select(m{"purchases": 0}).All(&products)

	for i, prod := range products {
		prod.Name = "*" + prod.Name + "*"

		j, _ := json.Marshal(prod)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, strconv.Itoa(i)+") "+string(j))
		msg.ParseMode = "Markdown"
		prodKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Remove"+emoji.Minus, "remove_product "+prod.ID.Hex()),
			),
		)
		msg.ReplyMarkup = prodKeyboard
		bot.Send(msg)
	}
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Here you come!")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}

// RemoveProductHandler removes product from "products" collection
func RemoveProductHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	prodID := strings.Split(update.CallbackQuery.Data, " ")[1]

	err := database.ProductsCollection.Remove(bson.M{"_id": bson.ObjectIdHex(prodID)})
	if err != nil {
		return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Removing has been FAILED! {"+err.Error()+"}")
	}

	return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "The product has been removed succesfully!")
}

// AddProductHandler adds product to database collection "products"
func AddProductHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	AddProductQueue[update.CallbackQuery.From.ID] = new(betypes.Product)

	return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Enter product name:")
}

// addProduct prompt user to get name, type and prise of product. Save it in DB
func AddProduct(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID

	prod := AddProductQueue[update.Message.From.ID]

	var err error
	if prod.Name == "" {
		prod.Name = update.Message.Text
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter product type:")
	} else if prod.Type == "" {
		prod.Type = update.Message.Text
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter product price:")
	} else {
		prod.Price, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again")
		}

		var answer tgbotapi.MessageConfig

		err = database.ProductsCollection.Insert(prod)
		if err != nil {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Product has not beed added {"+err.Error()+"}")
		} else {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Product has been added succesfully!")
		}

		delete(AddProductQueue, userID)
		answer.ReplyMarkup = keyboards.MainMenu
		return answer
	}
}