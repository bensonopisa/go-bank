package main

import (
	"math/rand"
	"net/http"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type Route struct {
	Method      string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
	Path        string
	Description string
}

type BaseResponse struct {
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
}

func NewAccount(name string) *Account {
	return &Account{
		ID:        rand.Intn(100),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
}

type Accounts []Account

type Routes []Route
