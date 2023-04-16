package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	isRunningLabel := ui.CreateIsRunningLabel()

	list := ui.CreateViewNotificationList(notifications)

	otherWidgets := container.NewVBox(
		ui.CreateAddNotificationVBox(db),
		ui.CreateSendNotificationToggle(
			func() {
				notifier.Toggle(application, notifications, notificationInterval, isRunningLabel)
			},
		),
		container.NewCenter(
			isRunningLabel,
		),
	)

	appContainer := container.NewBorder(
		nil,
		otherWidgets,
		nil,
		nil,
		list,
	)

	window.SetContent(appContainer)

	window.ShowAndRun()
}
