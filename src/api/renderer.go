package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/fgrehm/go-san/model"
	"github.com/fgrehm/sls-web/src/core/renderer"

	"github.com/labstack/echo"
)

func registerRendererEndpoints(e *echo.Echo, rs renderer.Service) {
	e.Any("/graph", func(c echo.Context) error {
		ct := contentType(c)

		var (
			img []byte
			err error
		)

		switch {
		case strings.HasPrefix(ct, "text/plain"):
			src, err := requestBody(c)
			if err != nil {
				return err
			}
			img, err = rs.RenderFromSource(src)
		case strings.HasPrefix(ct, "application/json"):
			m := &sanmodel.Model{}
			if err = c.Bind(m); err != nil {
				return err
			}
			img, err = rs.Render(m)
		default:
			return c.JSON(http.StatusNotAcceptable, map[string]string{"error": "Invalid request type"})
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%s", err)})
		}

		// TODO: Check if the client wants a PNG instead of a base64 of the graph

		img64 := base64.StdEncoding.EncodeToString(img)
		return c.JSON(http.StatusOK, map[string]string{"url": img64})
	})
}
