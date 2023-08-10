package main

import (
	"backend/handler"
	"context"
	_ "net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend/db"
)

func main() {
    e := echo.New()

	sqldb, err := db.PrepareDB(context.Background()); if err != nil {
		e.Logger.Fatal(err)
	}
	defer sqldb.Close()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Set("db", sqldb)
            return next(c)
        }
    })

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.Init(e)
    e.Logger.Fatal(e.Start(":1323"))
}
