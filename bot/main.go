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
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sethvargo/go-envconfig"
	_ "modernc.org/sqlite"
)

type appEnv struct {
	botToken           string   `env:"TELEGRAM_BOT_TOKEN, required"`
	admins             []string `env:"CONF_ADMINS, required"`
	confTitle          string   `env:"CONF_TITLE, required"`
	confLink           *url.URL `env:"CONF_LINK, required"`
	webhookAddr        *url.URL `env:"VIRTUAL_HOST, required"`
	confStart          time.Time
	confEnd            time.Time
	confratingDeadline time.Time
}

func (app *appEnv) fromArgs() error {
	confStart, err := time.Parse("02/01/2006 15:04:05", os.Getenv("CONF_START"))
	if err != nil {
		return fmt.Errorf("can't parse conference start: %w", err)
	}
	app.confStart = confStart

	confEnd, err := time.Parse("02/01/2006 15:04:05", os.Getenv("CONF_END"))
	if err != nil {
		return fmt.Errorf("can't parse conference end: %w", err)
	}
	app.confEnd = confEnd

	confRatingDeadline, err := time.Parse("02/01/2006 15:04:05", os.Getenv("CONF_RATING_DEADLINE"))
	if err != nil {
		return fmt.Errorf("can't parse conference rating deadline: %w", err)
	}
	app.confratingDeadline = confRatingDeadline

	if err := envconfig.Process(context.Background(), &app); err != nil {
		return err
	}

	return nil
}

func main() {
	var app appEnv
	err := app.fromArgs()
	if err != nil {
		log.Fatalf("can't parse env variables: %s", err.Error())
	}

	err = os.MkdirAll("/db", os.ModeDir)
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
		Status:          "active",
	})
	if err != nil {
		log.Fatalf("can't insert report in db: %s", err.Error())
	}

	reports, err := stmt.GetAllReports(context.Background(), app.confStart)
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

	b, err := bot.New(app.botToken, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: app.webhookAddr.String() + "/webhook",
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
