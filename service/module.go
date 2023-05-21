package service

import (
	"fmt"
	telBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	registerUser = "Register"
	listClients  = "List Clients"
	getClient    = "Get Client"
	createClient = "Create Client"
	deleteClient = "Delete Client"
)

func welcome(name string) string {
	prefix := "Welcome back "
	if name == "" {
		prefix = "Good day"
		name = "sir"
	}
	return fmt.Sprintf("%s %s! I'm Alfred, your personal user management butler. What can I do for you today?", prefix, name)
}

func createMainKeyboard(register bool) telBotAPI.ReplyKeyboardMarkup {
	if register {
		return telBotAPI.NewReplyKeyboard(
			telBotAPI.NewKeyboardButtonRow(
				telBotAPI.NewKeyboardButton(registerUser),
			),
		)
	} else {
		return telBotAPI.NewReplyKeyboard(
			telBotAPI.NewKeyboardButtonRow(
				telBotAPI.NewKeyboardButton(listClients),
				telBotAPI.NewKeyboardButton(createClient),
			),
			telBotAPI.NewKeyboardButtonRow(
				telBotAPI.NewKeyboardButton(getClient),
				telBotAPI.NewKeyboardButton(deleteClient),
			),
		)
	}
}

/*var inlineKeyboard = telBotAPI.NewInlineKeyboardMarkup(
	telBotAPI.NewInlineKeyboardRow(
		telBotAPI.NewInlineKeyboardButtonData("Add New User", AddUser),
		telBotAPI.NewInlineKeyboardButtonData("Test 1", Test1),
		telBotAPI.NewInlineKeyboardButtonData("Test 2", Test2),
	),
	telBotAPI.NewInlineKeyboardRow(
		telBotAPI.NewInlineKeyboardButtonData("Test 3", Test3),
		telBotAPI.NewInlineKeyboardButtonData("Test 4", Test4),
		telBotAPI.NewInlineKeyboardButtonData("Test 5", Test5),
	),
)*/
