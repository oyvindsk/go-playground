package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antage/eventsource"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaeljs1990/sqlitestore"
)

const (
	sessionName = "ost-session"
)

type server struct {
	db           *sql.DB
	logger       echo.Logger
	eventsource  eventsource.EventSource
	sessionStore *sqlitestore.SqliteStore
}

func main() {

	var srv server

	// Echo instance
	e := echo.New()
	e.Debug = true
	srv.logger = e.Logger // also a ec.Logger, but this is used when we're not inside a http requets? hm..

	// Init db
	dbpath := os.Getenv("DB_PATH")
	if dbpath == "" {
		log.Fatalln("DB_PATH environment variable must be sat")
	}
	srv.db = mustInitDB(e.Logger, dbpath)

	var err error
	// srv.sessionStore, err = sqlitestore.NewSqliteStore(dbpath, "sessions", "/", 3600, []byte("<SecretKey>"))
	srv.sessionStore, err = sqlitestore.NewSqliteStoreFromConnection(srv.db, "sessions", "/", 3600, []byte("<SecretKey>"))
	if err != nil {
		panic(err)
	}

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// e.GET("/", form)
	e.GET(
		"/",
		func(ec echo.Context) error {

			ec.Logger().Debugf("GET / \nheaders: \n %+v", ec.Request().Header)

			_, foo, err := srv.sessionGetOrCreate(ec.Request(), ec.Response())
			if err != nil {
				return err
			}

			if foo == "" {
				return ec.Redirect(http.StatusSeeOther, "/password")
			}

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

	e.GET("/password", func(ec echo.Context) error {

		_, foo, err := srv.sessionGetOrCreate(ec.Request(), ec.Response())
		if err != nil {
			return err
		}

		if foo != "" {
			return ec.Redirect(http.StatusSeeOther, "/")
		}

		return pagePassword(srv, ec).Render(ec.Response())
	})

	e.POST("/password", func(ec echo.Context) error {

		if ec.FormValue("password") != "røykepause" {
			ec.Echo().Logger.Infof("Wrong password from client: %q", ec.FormValue("password"))
			return echo.ErrUnauthorized
		}

		err := srv.sessionPut(ec.Request(), ec.Response())
		if err != nil {
			return err
		}
		return ec.Redirect(http.StatusSeeOther, "/")
		// )(srv, ec).Render(ec.Response())
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

// created, error
func (srv *server) sessionGetOrCreate(r *http.Request, w http.ResponseWriter) (bool, string, error) {

	session, err := srv.sessionStore.Get(r, sessionName)
	if err != nil {
		return false, "", err
	}

	srv.logger.Debugf("sessTest: Session:\n%+v", session)

	// if session.ID == "" {
	// 	return false, nil // fmt.Errorf("sessionGet: ID empty, no session")
	// }

	bar, ok := session.Values["bar"].(string)
	if ok {
		return session.IsNew, bar, nil
	}

	return session.IsNew, "", nil
}

func (srv *server) sessionPut(r *http.Request, w http.ResponseWriter) error {
	session, err := srv.sessionStore.Get(r, sessionName)
	session.Values["bar"] = "baz"
	session.Values["baz"] = time.Now().String()

	err = session.Save(r, w)
	// srv.logger.Debugf("sessTest: Session:\n%+v", session)
	return err
}
