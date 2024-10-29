package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const name = "dajare-api"

const version = "0.0.5"

var revision = "HEAD"

//go:embed static
var assets embed.FS

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
	e.GET("/api", func(c echo.Context) error {
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

	sub, _ := fs.Sub(assets, "static")
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(sub))))
	e.Logger.Fatal(e.Start(":8989"))
}
