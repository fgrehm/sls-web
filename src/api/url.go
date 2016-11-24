package api

import (
	"fmt"

	"github.com/fgrehm/sls-web/src/core/models"
	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/labstack/echo"
)

// REFACTOR: Should be using Echo's built in URL generation

func sanModelUrl(c echo.Context, id models.SanModelID) string {
	scheme := c.Request().Scheme()
	host := c.Request().Host()
	return fmt.Sprintf("%s://%s/models/%s", scheme, host, id)
}

func sanModelGraphUrl(c echo.Context, id models.SanModelID, transitionsHash string) string {
	return fmt.Sprintf("%s/graph/%s", sanModelUrl(c, id), transitionsHash)
}

func solutionUrl(c echo.Context, id solver.SolutionID) string {
	scheme := c.Request().Scheme()
	host := c.Request().Host()
	return fmt.Sprintf("%s://%s/solutions/%s", scheme, host, id)
}

func experimentUrl(c echo.Context, id solver.ExperimentID) string {
	scheme := c.Request().Scheme()
	host := c.Request().Host()
	return fmt.Sprintf("%s://%s/experiments/%s", scheme, host, id)
}
