package main

import (
	"bot/database"
	"context"
	"database/sql"
	"fmt"
	"io"
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
	os.MkdirAll("/data", os.ModeDir)
	db, _ := sql.Open("sqlite", "/data/botdb.db")
	defer db.Close()
	schema, _ := os.OpenFile("schema.sql", os.O_RDONLY, os.ModeExclusive)
	query, _ := io.ReadAll(schema)
	schema.Close()
	db.Exec(string(query))

	stmt, _ := database.Prepare(context.Background(), db)
	defer stmt.Close()
	stmt.CreateReport(context.Background(), database.CreateReportParams{
		Url:             "test" + strconv.Itoa(rand.Int()),
		Title:           "12345",
		StartingAt:      time.Now(),
		DurationMinutes: 45,
		Reporters:       "123 123 123",
		ConferenceID:    1,
		Status:          "active",
	})

	reports, _ := stmt.GetAllReports(context.Background())
	for i, v := range reports {
		fmt.Printf("%d. url = %s, title = %s, starting at = %s, duration = %d, reporters = %s, conference id = %d, status = %s\n", i, v.Url, v.Title, v.StartingAt.String(), v.DurationMinutes, v.Reporters, v.ConferenceID, v.Status)
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
