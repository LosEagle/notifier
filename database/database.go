package database

import (
	"database/sql"
	"fyne.io/fyne/v2"
	"log"
	"path/filepath"
)

func Init() *sql.DB {
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

	InsertNotifications(db, initialNotifications)

	return db
}

func InsertNotification(db *sql.DB, notification string) {
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

func InsertNotifications(db *sql.DB, notifications []string) {
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

func GetNotifications(db *sql.DB) []string {
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
