package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"net/http"
	"os"
)

const (
	vapidPublicKey  = "<тут public ключ>"
	vapidPrivateKey = "<тут private ключ>"
	mail            = "supermail@mail.rru"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/register", registerHandler)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// обрабатываем входящие параметры
	decoder := json.NewDecoder(r.Body)
	var subscription SubscriptionJSON
	err := decoder.Decode(&subscription)
	if err != nil {
		fmt.Println(err)
	}

	// SaveInDB()

	SendPush(subscription.Subscription, "Поздравляю! Вы подписались")
}

type SubscriptionJSON struct {
	Subscription webpush.Subscription `json:"subscription"`
}

func SendPush(subscription webpush.Subscription, text string) {
	// Decode subscription
	// Send Notification
	resp, err := webpush.SendNotification([]byte(text), &subscription, &webpush.Options{
		Subscriber:      mail, // Do not include "mailto:"
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
}
