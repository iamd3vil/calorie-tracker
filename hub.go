package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Hub contains everything needed for the app
type Hub struct {
	DB  *sqlx.DB
	Bot *tb.Bot
}

// NewHub returns a new Hub instance
func NewHub(db *sqlx.DB, bot *tb.Bot) *Hub {
	return &Hub{
		DB:  db,
		Bot: bot,
	}
}

// SetBudget sets daily caloie budget for a user
func (h *Hub) SetBudget(m *tb.Message) {
	fmt.Println(m.Text)
	budget, err := strconv.ParseInt(strings.Split(m.Text, " ")[1], 0, 64)
	if err != nil {
		h.Bot.Send(m.Sender, "Sorry!! Couldn't set daily budget")
		return
	}

	bud := Budget{
		UserID:      fmt.Sprint(m.Sender.ID),
		DailyBudget: budget,
	}

	const q = `
			INSERT INTO budgets (user_id, daily_budget) VALUES (:user_id, :daily_budget)
		`

	_, err = h.DB.NamedExec(q, bud)
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't set daily budget")
		return
	}

	h.Bot.Send(m.Sender, fmt.Sprintf("Daily budget of %d set. Good luck!!", budget))
}

// GetBudget gets set daily calorie budget
func (h *Hub) GetBudget(m *tb.Message) {
	fmt.Println(m.Text)
	bud := Budget{}

	err := h.DB.Get(&bud, "SELECT user_id, daily_budget FROM budgets WHERE user_id=$1", fmt.Sprint(m.Sender.ID))
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't get daily budget for you. Please try again")
		return
	}

	h.Bot.Send(m.Sender, fmt.Sprintf("Current budget set is %d. Good luck!!", bud.DailyBudget))
}
