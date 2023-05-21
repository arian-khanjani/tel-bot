package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tel-bot/mongodb"
	"tel-bot/utils"
	"time"

	telBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var bot *telBotAPI.BotAPI

func main() {
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Server...")

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	apiToken := utils.GetEnv("TELEGRAM_API_TOKEN", "", true)
	uri := utils.GetEnv("MONGO_URI", "", true)
	db := utils.GetEnv("MONGO_DB", "", true)
	coll := utils.GetEnv("MONGO_COLLECTION", "", true)

	repo, err := mongodb.New(mongodb.ConnProps{
		URI:  uri,
		DB:   db,
		Coll: coll,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("MongoDB connection established")

	defer func(repo *mongodb.Repo, ctx context.Context) {
		err := repo.Disconnect(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("MongoDB client disconnected")
	}(repo, ctx)

	/*indexes, err := repo.CreateIndexes(ctx, bson.D{ // TODO: indexes
		{"name", 1},
		{"email", 1},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("indexes added:", indexes)*/

	// ************************************************************************************

	bot, err = telBotAPI.NewBotAPI(apiToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := telBotAPI.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	listen(&updateConfig)

	// ************************************************************************************

	serverCtx, serverStopCtx := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		log.Println("Closing MongoDB connection...")
		err := repo.Disconnect(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Shutting down...")
		serverStopCtx()
	}()

	<-serverCtx.Done()
}

const (
	AddUser = "Add User"
	Test1   = "Test1"
	Test2   = "Test2"
	Test3   = "Test3"
	Test4   = "Test4"
	Test5   = "Test5"
)

var keyboard1 = telBotAPI.NewReplyKeyboard(
	telBotAPI.NewKeyboardButtonRow(
		telBotAPI.NewKeyboardButton(AddUser),
		telBotAPI.NewKeyboardButton(Test1),
		telBotAPI.NewKeyboardButton(Test2),
	),
	telBotAPI.NewKeyboardButtonRow(
		telBotAPI.NewKeyboardButton(Test3),
		telBotAPI.NewKeyboardButton(Test4),
		telBotAPI.NewKeyboardButton(Test5),
	),
)

/*var keyboard2 = telBotAPI.NewInlineKeyboardMarkup(
	telBotAPI.NewInlineKeyboardRow(
		telBotAPI.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		telBotAPI.NewInlineKeyboardButtonData("2", "two"),
		telBotAPI.NewInlineKeyboardButtonData("3", "three"),
	),
	telBotAPI.NewInlineKeyboardRow(
		telBotAPI.NewInlineKeyboardButtonData("4", "four"),
		telBotAPI.NewInlineKeyboardButtonData("5", "five"),
		telBotAPI.NewInlineKeyboardButtonData("6", "six"),
	),
)*/

func listen(updateConfig *telBotAPI.UpdateConfig) {
	updates := bot.GetUpdatesChan(*updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		reply := telBotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//reply.ReplyToMessageID = update.Message.MessageID

		/*switch update.Message.Command() {
		case "help":
			reply.Text = "I understand /sayhi and /status."
		case "sayhi":
			reply.Text = "Hi :)"
		case "status":
			reply.Text = "I'm ok."
		default:
			reply.Text = "I don't know that command"
		}*/

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				reply.Text = `
شلام جیگر
/start
/open
/close
`
			case "open":
				reply.ReplyMarkup = keyboard1
			case "close":
				reply.ReplyMarkup = telBotAPI.NewRemoveKeyboard(true)
			default:
				reply.Text = "I don't know that command"
			}
		} else {
			switch update.Message.Text {
			case AddUser:
				reply.Text = "ماخوای آدم اد کنی؟"
			case Test1:
				reply.Text = "شی ماخوای د‌ه"
			case Test2:
				reply.Text = "برو پی کارت کره خر"
			case Test3:
				reply.Text = "برو با او قین خرابت"
			case Test4:
				reply.Text = "روته برم پسررررر"
			case Test5:
				reply.Text = "خاله چیه تخمم نیستی"
			default:

			}
		}

		if _, err := bot.Send(reply); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}
