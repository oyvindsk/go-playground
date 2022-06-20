package main

import (
	"net/http/httputil"

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
