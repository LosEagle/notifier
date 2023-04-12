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
      sendNotification(&fyne.Notification{Title: "Notifier", Content: "Notification send stopped"})
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

func insertNotifications(db *sql.DB, notifications []string) {
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

func main() {
  application := app.NewWithID("com.github.loseagle.notifier")
  window := application.NewWindow("Notifier")

  db := initDb()

  notifications := getNotificationsFromDb(db)
  stopNotificationCh := make(chan struct{})
  notificationInterval := 10 * time.Second
  isNotifierRunning := false

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
