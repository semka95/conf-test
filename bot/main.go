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
	"github.com/jszwec/csvutil"
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
	if err := envconfig.Process(context.Background(), app); err != nil {
		return err
	}

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

	b.RegisterHandler(bot.HandlerTypeMessageText, "/sendReports", bot.MatchTypeExact, getReportsHandler)
}

func getReportsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// update.Message.Document.FileID
	fileParams := &bot.GetFileParams{
		FileID: update.Message.Document.FileID,
	}
	file, err := b.GetFile(ctx, fileParams)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "can't get file, upload csv file please",
		})
		return
	}

	resp, err := http.Get(b.FileDownloadLink(file))
	defer resp.Body.Close()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "can't download file",
		})
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "can't read file",
		})
		return
	}

	var talks []Talk
	if err := csvutil.Unmarshal(data, &talks); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("can't parse csv: %s", err.Error()),
		})
		return
	}

	for i, v := range talks {
		fmt.Printf("%d. %+v", i, v)
	}
}

type Talk struct {
	StartingAt      Time            `csv:"Start (MSK Time Zone)"`
	DurationMinutes DurationMinutes `csv:"Duration (min)"`
	Title           string          `csv:"Title"`
	Speakers        string          `csv:"Speakers"`
	URL             URL             `csv:"URL"`
}

type Time struct {
	time.Time
}

const format = "02/01/2006 15:04:05"

func (t *Time) UnmarshalCSV(data []byte) error {
	tt, err := time.Parse(format, string(data))
	if err != nil {
		return err
	}
	*t = Time{Time: tt}
	return nil
}

type DurationMinutes struct {
	time.Duration
}

func (d *DurationMinutes) UnmarshalCSV(data []byte) error {
	minutes, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	*d = DurationMinutes{Duration: time.Minute * time.Duration(minutes)}
	return nil
}

type URL struct {
	*url.URL
}

func (u *URL) UnmarshalCSV(data []byte) error {
	link, err := url.Parse(string(data))
	if err != nil {
		return err
	}

	*u = URL{URL: link}
	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
