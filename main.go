package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
	"booking/user"
	"time"
)

const URL = "https://inline.app/booking/-M3J_RCCRpHisvRHR5B7:inline-live-1/-M3J_RGIzVcfZGHU_9rj?language=zh-tw"
const TARGET = "線上訂位皆已滿"
const MSG = "Book now!!!!!!\n" + URL
var DB_URI, SAKIMOTO_TOKEN string

var repo user.Repository
var bot *tgbotapi.BotAPI

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB_URI = os.Getenv("DB_URI")
	SAKIMOTO_TOKEN = os.Getenv("SAKIMOTO_TOKEN")

	bot, err = tgbotapi.NewBotAPI(SAKIMOTO_TOKEN)
	if err != nil {
		log.Panic(err)
	}
}

func botListen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		fmt.Println("Command in")
		if update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				chatID := update.Message.Chat.ID
				username := update.Message.From.UserName
				_, err := repo.FindByChatID(chatID)
				if err != mongo.ErrNoDocuments {
					continue
				}
				err = repo.InsertUser(&user.User{
					Username: username,
					ChatID: chatID,
				})
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Panic(err)
	}
	repo = user.NewMongoRepository(client)

	go botListen()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var bookContent string
	for true {
		if err := chromedp.Run(ctx,
			chromedp.Navigate(URL),
			chromedp.InnerHTML("#book-now-content", &bookContent),
		); err != nil {
			fmt.Println(err)
		}

		if strings.Count(bookContent, TARGET) == 0 {
			users, err := repo.FindAll()
			if err != nil {
				log.Panic(err)
			}
			for _, v := range users {
				msg := tgbotapi.NewMessage(v.ChatID, MSG)
				bot.Send(msg)
			}
		}

		time.Sleep(1*time.Minute)
	}
}