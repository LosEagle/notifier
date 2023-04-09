package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"math/rand"
	"os"
	"time"
)

func sendNotificationsInIntervals(
	app fyne.App,
	text []string,
	interval time.Duration,
	ch chan struct{},
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.SendNotification(&fyne.Notification{Title: "Notifier", Content: text[rand.Intn(len(text))]})
		case <-ch:
			app.SendNotification(&fyne.Notification{Title: "Notifier", Content: "Notification send stopped"})
			return
		}
		time.Sleep(interval * time.Second)
	}
}

func getNotificationsFromJson(path string) []string {
	var notifications []string

	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonData, &notifications)
	if err != nil {
		log.Fatal(err)
	}

	return notifications
}

func main() {
	notifications := getNotificationsFromJson("./config/notifications.json")
	notificationChannel := make(chan struct{})
	notificationInterval := 5 * time.Second
	isNotifierRunning := false

	application := app.NewWithID("com.github.loseagle.notifier")
	window := application.NewWindow("Notifier")

	isRunningLabel := widget.NewLabel(fmt.Sprintf("Notification send is %v", isNotifierRunning))

	container := container.NewVBox(
		isRunningLabel,
		widget.NewButton("Toggle notification send", func() {
			if !isNotifierRunning {
				go sendNotificationsInIntervals(application, notifications, notificationInterval, notificationChannel)
			} else if isNotifierRunning {
				notificationChannel <- struct{}{}
			}

			isNotifierRunning = !isNotifierRunning

			labelText := fmt.Sprintf("Notification send is %v", isNotifierRunning)
			isRunningLabel.SetText(labelText)
		}))

	window.SetContent(container)

	window.ShowAndRun()
}
