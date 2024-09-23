package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", indexHandler)
	e.GET("/a", aHandler)
	e.GET("/b", bHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

func indexHandler(c echo.Context) error {
	return c.File("index.html")
}

func aHandler(c echo.Context) error {
	return c.File("a.html")
}

func bHandler(c echo.Context) error {
	return c.File("b.html")
}
