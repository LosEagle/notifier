package main

import (
  "database/sql"
  "fmt"
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "math/rand"
  "path/filepath"
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
      return
    }
  }
}

func initDb() *sql.DB {
  appStorageURI := fyne.CurrentApp().Storage().RootURI()
  appStoragePath := filepath.Join(appStorageURI.Path(), "notifier.db")
  initialNotifications := []string{
    "Be in the present",
    "Live in the moment",
    "Be in the now",
  }

  db, err := sql.Open("sqlite3", appStoragePath)
  if err != nil {
    log.Fatal(err)
  }

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS notifications (id INTEGER PRIMARY KEY, content TEXT)")
  if err != nil {
    log.Fatal(err)
  }

  insertNotifications(db, initialNotifications)

  return db
}

func insertNotification(db *sql.DB, notification string) {
  if notification == "" {
  }

  stmt, err := db.Prepare("INSERT INTO notifications(content) VALUES(?)")
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec(notification)
  if err != nil {
    log.Println(err)
    return
  }
}

func insertNotifications(db *sql.DB, notifications []string) {
  if len(notifications) == 0 {
  }

  stmt, err := db.Prepare("INSERT INTO notifications(content) VALUES(?)")
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()

  for _, content := range notifications {
    _, err = stmt.Exec(content)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func getNotificationsFromDb(db *sql.DB) []string {
  var notifications []string

  rows, err := db.Query("SELECT * FROM notifications")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  for rows.Next() {
    var id int
    var notification string
    err = rows.Scan(&id, &notification)
    if err != nil {
      log.Fatal(err)
    }

    notifications = append(notifications, notification)
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

  return notifications
}

func createViewNotificationList(notifications []string) *widget.List {
  list := widget.NewList(
    func() int {
      return len(notifications)
    },
    func() fyne.CanvasObject {
      return widget.NewLabel("")
    },
    func(index widget.ListItemID, item fyne.CanvasObject) {
      item.(*widget.Label).SetText(notifications[index])
    },
  )

  return list
}

func createAddNotificationVBox(db *sql.DB) *fyne.Container {
  notificationAddEntry := widget.NewEntry()
  notificationAddEntry.SetPlaceHolder("Add a notification")

  notificationAddSubmitButton := widget.NewButton(
    "Add", func() {
      insertNotification(db, notificationAddEntry.Text)
    },
  )

  return container.NewVBox(
    notificationAddEntry,
    notificationAddSubmitButton,
  )
}

func main() {
  application := app.NewWithID("com.github.loseagle.notifier")
  window := application.NewWindow("Notifier")

  db := initDb()

  notifications := getNotificationsFromDb(db)
  stopNotificationCh := make(chan struct{})
  notificationInterval := 30 * time.Minute
  isNotifierRunning := false

  isRunningLabel := widget.NewLabel(fmt.Sprintf("Notification send is %v", false))

  appContainer := container.NewVBox(
    isRunningLabel,
    createViewNotificationList(notifications),
    createAddNotificationVBox(db),
    widget.NewButton(
      "Toggle sending of notifications", func() {
        labelText := ""
        isRunningLabelText := "Notifications are being sent"
        isNotRunningLabelText := "Notifications are not being sent"

        if !isNotifierRunning {
          go sendNotificationInIntervals(
            application.SendNotification,
            notifications,
            notificationInterval,
            stopNotificationCh,
          )

          labelText = fmt.Sprintf(isRunningLabelText)
        } else if isNotifierRunning {
          stopNotificationCh <- struct{}{}

          labelText = fmt.Sprintf(isNotRunningLabelText)
        }

        isNotifierRunning = !isNotifierRunning

        isRunningLabel.SetText(labelText)
      },
    ),
  )

  window.SetContent(appContainer)

  window.ShowAndRun()
}
