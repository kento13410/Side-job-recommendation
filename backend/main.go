package main

import (
	"backend/handler"
	_"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.Init(e)
    e.Logger.Fatal(e.Start(":1323"))
}
