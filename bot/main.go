package main

import (
	"bot/database"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "modernc.org/sqlite"
)

// Send any text message to the bot after the bot has been started

func main() {
	err := os.MkdirAll("/db", os.ModeDir)
	if err != nil {
		log.Fatalf("can't create folder for db: %s", err.Error())
	}

	db, err := sql.Open("sqlite", "/db/botdb.db")
	defer db.Close()
	if err != nil {
		log.Fatalf("can't create db connection: %s", err.Error())
	}

	schema, err := os.OpenFile("schema.sql", os.O_RDONLY, os.ModeExclusive)
	if err != nil {
		log.Fatalf("can't open schema file: %s", err.Error())
	}
	defer schema.Close()

	query, err := io.ReadAll(schema)
	if err != nil {
		log.Fatalf("can't read schema file: %s", err.Error())
	}
	schema.Close()

	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatalf("can't create tables: %s", err.Error())
	}

	stmt, err := database.Prepare(context.Background(), db)
	defer stmt.Close()
	if err != nil {
		log.Fatalf("can't create prepare statement: %s", err.Error())
	}

	err = stmt.CreateReport(context.Background(), database.CreateReportParams{
		Url:             "test" + strconv.Itoa(rand.Int()),
		Title:           "12345",
		StartingAt:      time.Now(),
		DurationMinutes: 45,
		Reporters:       "123 123 123",
		ConferenceID:    1,
		Status:          "active",
	})
	if err != nil {
		log.Fatalf("can't insert report in db: %s", err.Error())
	}

	reports, err := stmt.GetAllReports(context.Background(), 1)
	if err != nil {
		log.Fatalf("can't fetch reports: %s", err.Error())
	}
	for i, v := range reports {
		fmt.Printf("%d. url = %s, title = %s, starting at = %s, duration = %d, reporters = %s\n", i, v.Url, v.Title, v.StartingAt.String(), v.DurationMinutes, v.Reporters)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("EXAMPLE_TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: os.Getenv("VIRTUAL_HOST") + "/webhook",
	})

	go func() {
		http.ListenAndServe(":2000", b.WebhookHandler())
	}()

	// Use StartWebhook instead of Start
	b.StartWebhook(ctx)

	// call methods.DeleteWebhook if needed
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
