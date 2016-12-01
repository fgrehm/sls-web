package api

import (
	"net/http"
	"strings"

	"github.com/fgrehm/go-san/model"
	"github.com/fgrehm/sls-web/src/core/compiler"

	"github.com/labstack/echo"
)

func registerCompilerEndpoints(e *echo.Echo, cs compiler.Service) {
	e.Any("/parse", func(c echo.Context) error {
		if !strings.HasPrefix(contentType(c), "text/plain") {
			return notAcceptable(c)
		}

		src, err := requestBody(c)
		if err != nil {
			return err
		}

		res, err := cs.Parse(src)
		if err != nil {
			return serverError(c, "Error parsing model", err)
		}

		return c.JSON(http.StatusOK, res)
	})

	e.Any("/compile", func(c echo.Context) error {
		if !strings.HasPrefix(contentType(c), "application/json") {
			return notAcceptable(c)
		}

		m := &sanmodel.Model{}
		if err := c.Bind(m); err != nil {
			return err
		}

		out, err := cs.Compile(m)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, string(out))
	})
}
