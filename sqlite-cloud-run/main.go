package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *sql.DB // global, for now .. todo: OK to share this? Probably not FIXME

func main() {

	// Echo instance
	e := echo.New()
	e.Debug = true

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", form)
	e.POST("/", formResult)
	e.GET("/print", printRows)

	// Init db
	db = mustInitDB(e.Logger)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handlers
func printRows(c echo.Context) error {

	err := fprintAll(c.Logger(), db, c.Response())
	if err != nil {
		c.Logger().Errorf("printRows: db fprintAll: %s", err)
		return err
	}

	return nil
}

func form(c echo.Context) error {

	return c.HTML(http.StatusOK, `<html>
		<body>
			<form action="" method="post">
				<label for="msg">Message (min 4, max 500 characters):</label>
				<input type="text" id="msg" name="msg" required minlength="4" maxlength="500" size="100">
				<input type="submit" value="Store">
			</form>
			<p>
			</p>
		</body>
		</html>`)
	// fmt.Errorf("form error1!")
}

func formResult(c echo.Context) error {

	vals, err := c.FormParams()
	if err != nil {
		c.Logger().Errorf("formResult: %s", err)
		return err
	}

	msg := vals.Get("msg")
	if len(msg) < 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "formResult Input msg too short")
	}

	c.Logger().Infof("formResult: Msg: %q", msg)

	err = storeMsg(c.Logger(), db, msg)
	if err != nil {
		c.Logger().Errorf("formResult: storeMsg: %s", err)
		return err
	}

	return c.Redirect(http.StatusFound, "print")
}
