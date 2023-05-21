package service

import (
	"context"
	"fmt"
	telBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"tel-bot/model"
	"tel-bot/mongodb"
)

type Alfred struct {
	bot     *telBotAPI.BotAPI
	updates telBotAPI.UpdatesChannel
	repo    *mongodb.Repo
}

func InitBot(ctx context.Context, token string, repo *mongodb.Repo) error {
	bot, err := telBotAPI.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = true

	updateConfig := telBotAPI.NewUpdate(0)
	updateConfig.Timeout = 60

	alf := Alfred{
		bot:     bot,
		updates: bot.GetUpdatesChan(updateConfig),
		repo:    repo,
	}

	alf.listen(ctx)

	return nil
}

func (a *Alfred) listen(ctx context.Context) {
	for update := range a.updates {
		if update.Message == nil {
			continue
		}

		reply := telBotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//reply.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				register := false
				user, err := a.repo.GetUser(ctx, update.Message.Chat.ID)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						register = true
					} else {
						log.Fatalln(err)
					}
				}
				reply.Text = welcome(user.FirstName)
				reply.ReplyMarkup = createMainKeyboard(register)
			case "close":
				reply.ReplyMarkup = telBotAPI.NewRemoveKeyboard(true)
			default:
				reply.Text = ""
			}
		} else {
			switch update.Message.Text {
			case registerUser:
				u := model.User{
					ID:        update.Message.Chat.ID,
					Username:  update.Message.Chat.UserName,
					FirstName: update.Message.Chat.FirstName,
					Account:   model.Account{Balance: 0},
				}
				user, err := a.repo.CreateUser(ctx, &u)
				if err != nil {
					log.Println(err)
					reply.Text = fmt.Sprintf("I'm afraid there was an error with creating your account: %s", err.Error())
					break
				}

				reply.Text = fmt.Sprintf("Welcome %s!", user.FirstName)
				reply.ReplyMarkup = createMainKeyboard(false)
			case listClients:
				clients, err := a.repo.ListClients(ctx, update.Message.Chat.ID)
				if err != nil {
					log.Println(err)
					reply.Text = fmt.Sprintf("I'm afraid there was an error with listing your clients: %s", err.Error())
					break
				}
				var res string
				for i, client := range *clients {
					res += fmt.Sprintf("%d: \nID: %d\nUsername: %s\n--------------", i, client.ID, client.Username)
				}
				reply.Text = res
				reply.ReplyMarkup = createMainKeyboard(false)
			case createClient:
				cID := int64(uuid.New().ID())
				c := model.Client{
					ID:         cID,
					ProviderID: update.Message.Chat.ID,
					Username:   fmt.Sprintf("bot_%d", cID),
				}
				client, err := a.repo.CreateClient(ctx, &c)
				if err != nil {
					log.Println(err)
					reply.Text = fmt.Sprintf("I'm afraid there was an error with creating this client: %s", err.Error())
					break
				}
				reply.Text = fmt.Sprintf("Your client was created successfuly\nUsername: %s", client.Username)
				reply.ReplyMarkup = createMainKeyboard(false)
			default:
				if update.Message.Text == "" {
					update.Message.Text = "your command"
				}
				reply.Text = fmt.Sprintf("I'm afraid I don't understand '%s', sir.", update.Message.Text)
			}
		}

		if _, err := a.bot.Send(reply); err != nil {
			log.Fatalln(err)
		}
	}
}
