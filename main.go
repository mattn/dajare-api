package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const name = "dajare-api"

const version = "0.0.3"

var revision = "HEAD"

func main() {
	var dsn string
	var ver bool

	flag.StringVar(&dsn, "dsn", os.Getenv("DATABASE_URL"), "Database source")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.Parse()

	if ver {
		fmt.Println(version)
		os.Exit(0)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		var text string
		err = db.QueryRow(`SELECT text FROM dajare ORDER BY RANDOM() LIMIT 1`).Scan(&text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(text)
		return c.JSON(http.StatusOK, struct {
			Text string `json:"text"`
		}{
			Text: text,
		})
	})
	e.Logger.Fatal(e.Start(":8989"))
}
