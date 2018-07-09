package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/VOID404/GoWebChat/utils"
)

func main() {
	e := echo.New()

	e.Static("/", "static/")
	e.GET("/", fileHandler("html/index.html"))
	e.GET("/socket/:username", socket)

	// TODO: Check username availability endpoint

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${latency_human} ${uri}\n",
	}))

	port := "3000"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func fileHandler(filename string) echo.HandlerFunc {
	return func(e echo.Context) error {
		return e.File(filename)
	}
}

var chat = utils.NewChat()

func socket(c echo.Context) error {
	return chat.Chatify(c)
}
