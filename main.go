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

func sendNotificationInIntervals(
	sendNotification func(*fyne.Notification),
	notifications []string,
	interval time.Duration,
	stop chan struct{},
) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			sendNotification(&fyne.Notification{Title: "Notifier", Content: notifications[rand.Intn(len(notifications))]})
		case <-stop:
			ticker.Stop()
			sendNotification(&fyne.Notification{Title: "Notifier", Content: "Notification send stopped"})
			return
		}
	}
}

func getNotificationsFromJson(path string) []string {
	var notifications []string

	jsonData, err := os.ReadFile(path)
	log.Println(jsonData)
	log.Println(err)
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
	stopNotificationCh := make(chan struct{})
	notificationInterval := 30 * time.Minute
	isNotifierRunning := false

	application := app.NewWithID("com.github.loseagle.notifier")
	window := application.NewWindow("Notifier")

	isRunningLabel := widget.NewLabel(fmt.Sprintf("Notification send is %v", false))

	appContainer := container.NewVBox(
		isRunningLabel,
		widget.NewButton(
			"Toggle notification send", func() {
				if !isNotifierRunning {
					go sendNotificationInIntervals(
						application.SendNotification,
						notifications,
						notificationInterval,
						stopNotificationCh,
					)
				} else if isNotifierRunning {
					stopNotificationCh <- struct{}{}
				}

				isNotifierRunning = !isNotifierRunning

				labelText := fmt.Sprintf("Notification send is %v", isNotifierRunning)
				isRunningLabel.SetText(labelText)
			},
		),
	)

	window.SetContent(appContainer)

	window.ShowAndRun()
}
