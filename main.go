package main

import (
	"fmt"
	"log"
	"time"

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
	date TEXT NOT NULL,
	calories INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES budgets(id)
);

CREATE TABLE budgets (
	id INTEGER PRIMARY KEY,
	user_id TEXT,
	daily_budget INTEGER NOT NULL
);`

type Budget struct {
	UserID      string `db:"user_id"`
	DailyBudget int64  `db:"daily_budget"`
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

	b.Handle("/hello", func(m *tb.Message) {
		fmt.Println(m.Text)
		b.Send(m.Sender, "hello world")
	})

	b.Handle("/setbudget", h.SetBudget)

	b.Handle("/budget", h.GetBudget)

	b.Start()
}
