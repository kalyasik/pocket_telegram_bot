package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kalyasik/pocket_telegram_bot/pkg/repository"
	"github.com/kalyasik/pocket_telegram_bot/pkg/repository/boltdb"
	"github.com/kalyasik/pocket_telegram_bot/pkg/server"
	"github.com/kalyasik/pocket_telegram_bot/pkg/telegram"
	"github.com/spf13/viper"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	/* Initialize config */
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to load config file: %s", err.Error())
	}

	/* Initialize tgbotapi */
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	/* Initialize pocket */
	pocketClient, err := pocket.NewClient(viper.GetString("POCKET_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	/* Initialize database */
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	autorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https:/t.me/pocket_telegram_bot")

	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err = autorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
