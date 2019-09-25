package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	bud := Budget{}

	err := h.DB.Get(&bud, "SELECT user_id, daily_budget FROM budgets WHERE user_id=$1", fmt.Sprint(m.Sender.ID))
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't get daily budget for you. Please try again")
		return
	}

	h.Bot.Send(m.Sender, fmt.Sprintf("Current budget set is %d. Good luck!!", bud.DailyBudget))
}

// SetEntry sets an entry
func (h *Hub) SetEntry(m *tb.Message) {
	name := strings.Split(m.Text, " ")[1]

	calories, err := strconv.ParseInt(strings.Split(m.Text, " ")[2], 0, 64)
	if err != nil {
		h.Bot.Send(m.Sender, "Sorry!! Couldn't write entry")
		return
	}

	date := time.Now().Format("2-Jan-2006")

	// Get daily budget for the user
	bud := Budget{}

	err = h.DB.Get(&bud, "SELECT id, user_id, daily_budget FROM budgets WHERE user_id=$1", fmt.Sprint(m.Sender.ID))
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't write entry. Please try again")
		return
	}

	entry := Entry{
		UserID:   bud.ID,
		Name:     name,
		Date:     date,
		Calories: calories,
	}

	const q = `
		INSERT INTO entries (user_id, name, date, calories) VALUES (:user_id, :name, :date, :calories)
	`

	_, err = h.DB.NamedExec(q, entry)
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't write entry. Please try again")
		return
	}

	allCal := Entry{}
	// Get sum of all the calories for today
	const q2 = `
		SELECT sum(calories) as calories FROM entries WHERE date=$1
	`

	err = h.DB.Get(&allCal, q2, date)
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't write entry. Please try again")
		return
	}

	h.Bot.Send(m.Sender, fmt.Sprintf("Entry added. You can still consume: %d", bud.DailyBudget-allCal.Calories))
}

// ClearEntries clears out all entries for today
func (h *Hub) ClearEntries(m *tb.Message) {
	const q = `DELETE FROM entries WHERE DATE=$1`
	date := time.Now().Format("2-Jan-2006")

	_, err := h.DB.Exec(q, date)
	if err != nil {
		log.Printf("Err: %v", err)
		h.Bot.Send(m.Sender, "Sorry!! Couldn't clear entries. Please try again")
		return
	}

	h.Bot.Send(m.Sender, "All entries for today cleared")
}
