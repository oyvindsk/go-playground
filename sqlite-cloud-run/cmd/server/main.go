package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/antage/eventsource"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// var db *sql.DB // global, for now .. todo: OK to share this? Probably not FIXME

type server struct {
	db          *sql.DB
	logger      echo.Logger
	eventsource eventsource.EventSource
}

func main() {

	dbpath := os.Getenv("DB_PATH")
	if dbpath == "" {
		log.Fatalln("DB_PATH environment variable must be sat")
	}

	var srv server

	// Echo instance
	e := echo.New()
	e.Debug = true
	srv.logger = e.Logger // also a ec.Logger, but this is used when we're not inside a http requets? hm..

	// Init db
	srv.db = mustInitDB(e.Logger, dbpath)

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// e.GET("/", form)
	e.GET(
		"/",
		func(ec echo.Context) error {

			ec.Logger().Debugf("GET / \nheaders: \n %+v", ec.Request().Header)

			// Get exisiting rows from the db
			messages, err := getAll(srv.db, ec.Logger())
			if err != nil {
				return err
			}

			// Render the page
			return pageHome(srv, ec, messages).Render(ec.Response())
		},
	)

	e.POST(
		"/",
		func(ec echo.Context) error {

			ec.Logger().Debugf("POST / \nheaders: \n %+v", ec.Request().Header)

			vals, err := ec.FormParams()
			if err != nil {
				ec.Logger().Errorf("formResult: %s", err)
				return err
			}

			author := vals.Get("author")

			msg := vals.Get("msg")
			if len(msg) < 5 {
				return echo.NewHTTPError(http.StatusBadRequest, "formResult Input msg too short")
			}

			ec.Logger().Infof("formResult: Msg: %q", msg)

			err = storeMsg(ec.Logger(), srv.db, author, msg)
			if err != nil {
				ec.Logger().Errorf("formResult: storeMsg: %s", err)
				return err
			}

			srv.eventsource.SendEventMessage("<tr><td>"+author+"</td> <td>"+msg+"</td></tr>", "tick-event", "1") // FIXME

			return nil // return ec.Redirect(http.StatusFound, "/")
		},
	)

	es := eventsource.New(nil, nil)
	defer es.Close()

	srv.eventsource = es

	e.Any("/sse-element", func(ec echo.Context) error {
		ec.Echo().Logger.Debugf("In /sse !!")
		ec.Logger().Debugf("headers: \n %+v", ec.Request().Header)
		es.ServeHTTP(ec.Response(), ec.Request())
		return nil
	})

	// go func() {
	// 	id := 1
	// 	for {
	// 		es.SendEventMessage(fmt.Sprintf("<li>tick - %d</li>", id), "tick-event", strconv.Itoa(id))
	// 		id++
	// 		e.Logger.Debugf("sse: on id %d, conected: %d", id, es.ConsumersCount())

	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
