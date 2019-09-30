package main

import (
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	initConfig()
}

const createSchema = `
CREATE TABLE entries (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	date TEXT NOT NULL,
	calories INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES budgets(id)
);

CREATE TABLE budgets (
	id INTEGER PRIMARY KEY,
	user_id TEXT NOT NULL UNIQUE,
	daily_budget INTEGER NOT NULL
);`

type Budget struct {
	ID          int64  `db:"id"`
	UserID      string `db:"user_id"`
	DailyBudget int64  `db:"daily_budget"`
}

type Entry struct {
	UserID   int64  `db:"user_id"`
	Date     string `db:"date"`
	Name     string `db:"name"`
	Calories int64  `db:"calories"`
}

func main() {
	db, err := sqlx.Connect("sqlite3", cfg.Db.SqlitePath)
	if err != nil {
		log.Fatalln(err)
	}

	db.Exec(createSchema)

	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Telegram.ApiKey,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	h := NewHub(db, b)
	gocron.Every(1).Day().At("07:00").Do(h.SendStats)

	go func() {
		<-gocron.Start()
	}()

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello")
	})

	b.Handle("/setbudget", h.SetBudget)
	b.Handle("/budget", h.GetBudget)
	b.Handle("/add", h.SetEntry)
	b.Handle("/clear", h.ClearEntries)

	b.Start()
}
