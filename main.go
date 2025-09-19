package main

import (
	"log"
	"os"

	"go-telegram-expense-bot/handler"
	"go-telegram-expense-bot/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[!] .env file not found, using system env vars")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	credFile := os.Getenv("GOOGLE_CREDENTIALS")
	sheetID := os.Getenv("GOOGLE_SHEET_ID")

	if botToken == "" || credFile == "" || sheetID == "" {
		log.Fatal("[X] Missing required env vars (TELEGRAM_BOT_TOKEN, GOOGLE_CREDENTIALS, GOOGLE_SHEET_ID)")
	}

	gsService, err := service.NewGoogleSheets(credFile, sheetID)
	if err != nil {
		log.Fatalf("[X] Failed to init Google Sheets service: %v", err)
	}

	tgService, err := service.NewTelegram(botToken)
	if err != nil {
		log.Fatalf("[X] Failed to init Telegram service: %v", err)
	}

	log.Printf("[V] Authorized on account %s", tgService.Bot.Self.UserName)
	handler.NewTelegramHandler(tgService, gsService).Start()
}
