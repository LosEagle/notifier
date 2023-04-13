package notification

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"time"
)

type Notifier struct {
	IsRunning bool
	StopCh    chan struct{}
}

func NewNotifier() *Notifier {
	return &Notifier{
		IsRunning: false,
		StopCh:    make(chan struct{}),
	}
}

func (n *Notifier) sendNotificationInIntervals(
	sendNotification func(*fyne.Notification),
	notifications []string,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			sendNotification(&fyne.Notification{Title: "Notifier", Content: notifications[rand.Intn(len(notifications))]})
		case <-n.StopCh:
			ticker.Stop()
			return
		}
	}
}

func (n *Notifier) Toggle(
	application fyne.App,
	notifications []string,
	notificationInterval time.Duration,
	isRunningLabel *widget.Label,
) {
	notifierRunningLabelText := "Application is running"
	notifierNotRunningLabelText := "Application is not running"

	if !n.IsRunning {
		go n.sendNotificationInIntervals(
			application.SendNotification,
			notifications,
			notificationInterval,
		)
	} else {
		n.StopCh <- struct{}{}
	}

	n.IsRunning = !n.IsRunning

	if n.IsRunning {
		isRunningLabel.SetText(notifierRunningLabelText)
	} else {
		isRunningLabel.SetText(notifierNotRunningLabelText)
	}
}
