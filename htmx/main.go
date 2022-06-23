package main

import (
	"fmt"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/antage/eventsource"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	// *template.Template
}

func main() {

	srv := server{}

	e := echo.New()
	e.Debug = true

	e.Logger.SetHeader(`${time_rfc3339_nano} ${level} ${short_file}:${line}`)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: `${time_rfc3339_nano} ${id} ${remote_ip}, ${host} ${method} ${uri} ${user_agent},` +
		`status: ${status}, error: ${error},latency: ${latency_human}, b_in: ${bytes_in}, b_out: ${bytes_out}` + "\n"}))
	e.Any("/*", srv.handler)

	es := eventsource.New(nil, nil)
	defer es.Close()
	e.Any("/sse-element", func(ec echo.Context) error {
		ec.Echo().Logger.Debugf("In /sse !!")
		es.ServeHTTP(ec.Response(), ec.Request())
		return nil
	})

	go func() {
		id := 1
		for {
			es.SendEventMessage(fmt.Sprintf("<li>tick - %d</li>", id), "tick-event", strconv.Itoa(id))
			id++
			e.Logger.Debugf("sse: on id %d, conected: %d", id, es.ConsumersCount())

			time.Sleep(2 * time.Second)
		}
	}()

	// log.Fatal(http.ListenAndServe(":8080", nil))

	e.Logger.Fatal(e.Start(":8080"))
}

func (srv server) handler(ec echo.Context) error {

	templates := mustParseTemplates()

	path := ec.Request().URL.Path
	path = path[1:] + ".html"
	ec.Echo().Logger.Debugf("handler: serving path: %q", path)

	debugBody, err := httputil.DumpRequest(ec.Request(), true)
	if err != nil {
		ec.Echo().Logger.Errorf("handler: failed dumping Request: %w", err)
	}

	ec.Echo().Logger.Debugf("handler: requst:\t %s", debugBody)

	return templates.ExecuteTemplate(ec.Response(), path, nil)
}
