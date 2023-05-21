package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tel-bot/mongodb"
	"tel-bot/service"
	"tel-bot/utils"
	"time"
)

func main() {
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("Starting Server...")

	/*err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}*/

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

	indexes, err := repo.CreateIndexes(ctx, bson.D{
		{"provider_id", 1},
		{"username", 1},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("indexes added:", indexes)

	err = service.InitBot(ctx, apiToken, repo)
	if err != nil {
		log.Fatalln(err)
	}

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
