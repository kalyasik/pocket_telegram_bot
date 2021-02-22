package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kalyasik/pocket_telegram_bot/pkg/telegram"
	"github.com/spf13/viper"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	/* Initialize config */
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to load config file: %s", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(viper.GetString("POCKET_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost/")
	if err = telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
