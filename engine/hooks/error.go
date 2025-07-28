package hooks

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DefaultHTTPErrorResponse(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		fmt.Println(he.Message, he.Unwrap())
	}

	c.JSON(code, map[string]any{
		"error": err.Error(),
	})
}
