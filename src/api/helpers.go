package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fgrehm/sls-web/src/core/models"

	"github.com/labstack/echo"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func contentType(c echo.Context) string {
	return c.Request().Header().Get("Content-Type")
}

func requestBody(c echo.Context) ([]byte, error) {
	return ioutil.ReadAll(c.Request().Body())
}

func requestBodyStr(c echo.Context) (string, error) {
	body, err := requestBody(c)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func serverError(c echo.Context, msg string, err error) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": msg,
		"error":   fmt.Sprintf("%s", err),
	})
}

func notAcceptable(c echo.Context) error {
	return c.JSON(http.StatusNotAcceptable, map[string]string{"error": "Content type"})
}

func parseSanModelUpsertInput(c echo.Context) (*models.UpsertInput, error) {
	ct := contentType(c)
	input := &models.UpsertInput{}

	switch {
	case strings.HasPrefix(ct, "text/plain"):
		src, err := requestBodyStr(c)
		if err != nil {
			return nil, serverError(c, "Unable to fetch request body", err)
		}
		input.Source = &src
	case strings.HasPrefix(ct, "application/json"):
		if err := c.Bind(input); err != nil {
			return nil, serverError(c, "Unable to parse request", err)
		}
	default:
		return nil, notAcceptable(c)
	}

	return input, nil
}
