package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"notifier/database"
	"notifier/notification"
	"notifier/ui"
	"time"
)

func main() {
	application := app.NewWithID("com.github.loseagle.notifier")
	window := application.NewWindow("Notifier")
	notifier := notification.NewNotifier()

	db := database.Init()

	notifications := database.GetNotifications(db)
	notificationInterval := 30 * time.Second

	isRunningLabel := widget.NewLabel(fmt.Sprintf("Notification send is %v", false))

	appContainer := container.NewVBox(
		isRunningLabel,
		ui.CreateViewNotificationList(notifications),
		ui.CreateAddNotificationVBox(db),
		ui.CreateSendNotificationToggle(
			func() {
				notifier.Toggle(application, notifications, notificationInterval, isRunningLabel)
			},
		),
	)

	window.SetContent(appContainer)

	window.ShowAndRun()
}
