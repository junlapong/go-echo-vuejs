package main

import (
	"database/sql"
	"time"

	"todo/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tylerb/graceful"
)

func main() {

	db := initDB("storage.db")
	migrate(db)

	e := echo.New()

	// logging
	if logger, ok := e.Logger.(*log.Logger); ok {
		// logger.SetHeader("${time_rfc3339} ${level}")
		logger.SetHeader("${level}\t")
	}
	e.Logger.SetLevel(log.DEBUG)

	e.File("/", "public/index.html")
	e.GET("/tasks", handlers.GetTasks(db))
	e.PUT("/tasks", handlers.PutTask(db))
	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// server start !!!
	e.Server.Addr = ":8080"
	e.Logger.Infof("listen, url: http://localhost%s/", e.Server.Addr)

	err := graceful.ListenAndServe(e.Server, 30*time.Second)
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR NOT NULL
	);
	`

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}
