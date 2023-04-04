package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"time"
)

func sendNotificationsInIntervals(app fyne.App, text [2]string, interval time.Duration, ch chan struct{}) {
	fmt.Sprintf("sendNotificationsInIntervals: %v", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.SendNotification(&fyne.Notification{Title: "Notifier", Content: text[rand.Intn(len(text))]})
		case <-ch:
			return
		}
		time.Sleep(interval * time.Second)
	}
}

func main() {
	notificationLabels := [...]string{"First notification", "Second notification"}
	notificationInterval := 20 * time.Second
	notificationChannel := make(chan struct{})
	isNotificationSendRunning := false

	app := app.New()
	window := app.NewWindow("Notifier")

	isRunningLabel := widget.NewLabel(fmt.Sprintf("Notification send is %v", isNotificationSendRunning))

	container := container.NewVBox(
		isRunningLabel,
		widget.NewButton("Toggle notification send", func() {
			if !isNotificationSendRunning {
				go sendNotificationsInIntervals(app, notificationLabels, notificationInterval, notificationChannel)
			} else if isNotificationSendRunning {
				notificationChannel <- struct{}{}
			}

			isNotificationSendRunning = !isNotificationSendRunning
			isRunningLabel.Text = fmt.Sprintf("Notification send is %v", isNotificationSendRunning)
		}))

	window.SetContent(container)

	window.ShowAndRun()
}
