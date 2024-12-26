package main

import (
	"GymBot/internal/application"
	"GymBot/internal/infrastructure/postgres"
	"GymBot/internal/interface/telegram"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func init() {

	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found, using system environment variables.")
	}
}

func main() {
	// Read environment variables
	dbConnStr := os.Getenv("DB_CONNECTION")
	botToken := os.Getenv("BOT_TOKEN")

	if dbConnStr == "" || botToken == "" {
		log.Fatal("Failed to load environment variables. Check BOT_TOKEN and DB_CONNECTION.")
	}

	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	repo := postgres.NewUserRepositoryDb(db)
	service := application.Initialize(repo) // Initialize service

	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	botHandler := telegram.NewBotHandler(service) // Pass initialized service
	bot.Handle(telebot.OnText, botHandler.MsgMainHandler)
	bot.Handle(telebot.OnCallback, botHandler.DataHandler)

	slog.Info("Bot started.")
	bot.Start()
}

//TODO Перенести все подключения и инициализации чтобы тут было только чтение енв файла и старт всего нужного
