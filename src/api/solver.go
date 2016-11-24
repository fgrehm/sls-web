package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/labstack/echo"
)

func registerSolverEndpoints(e *echo.Echo, ss solver.Service) {
	e.POST("/solve", scheduleSolution(ss))
	e.GET("/solutions/:solutionID", findSolution(ss))
	e.POST("/experiments", scheduleExperiment(ss))
	e.GET("/experiments/:experimentID", findExperiment(ss))
}

func scheduleSolution(ss solver.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		ct := contentType(c)
		input := &solver.SolutionInput{}

		switch {
		case strings.HasPrefix(ct, "text/plain"):
			src, err := requestBodyStr(c)
			if err != nil {
				return err
			}
			input.Source = src
			// TODO: Fetch maxIterations and tolerance from request params
		case strings.HasPrefix(ct, "application/json"):
			if err := c.Bind(input); err != nil {
				return serverError(c, "Error parsing request", err)
			}
		default:
			return notAcceptable(c)
		}

		solution, err := ss.ScheduleSolution(input)
		if err != nil {
			return serverError(c, "Error scheduling solution", err)
		}

		return c.JSON(http.StatusOK, buildSolutionResource(c, solution))
	}
}

func findSolution(ss solver.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := solver.SolutionID(c.Param("solutionID"))
		s, err := ss.FindSolution(id, nil)
		if err != nil {
			if err == solver.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Solution does not exist"})
			} else {
				return serverError(c, "Error loading solution", err)
			}
		}

		return c.JSON(http.StatusOK, buildSolutionResource(c, s))
	}
}

func scheduleExperiment(ss solver.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		ct := contentType(c)
		input := &solver.ExperimentInput{}

		switch {
		case strings.HasPrefix(ct, "text/plain"):
			src, err := requestBodyStr(c)
			if err != nil {
				return err
			}
			input.Source = src
			// REFACTOR
			if c.QueryParam("name") == "" || c.QueryParam("from") == "" || c.QueryParam("to") == "" || c.QueryParam("increment") == "" {
				// TODO: Render a better error
				return c.JSON(http.StatusBadRequest, "Missing parameters for experiment")
			}

			from, err := strconv.ParseFloat(c.QueryParam("from"), 64)
			if err != nil {
				return err
			}
			to, err := strconv.ParseFloat(c.QueryParam("to"), 64)
			if err != nil {
				return err
			}
			increment, err := strconv.ParseFloat(c.QueryParam("increment"), 64)
			if err != nil {
				return err
			}

			input.Identifier = solver.ExperimentIdentifier{
				Name:      c.QueryParam("name"),
				From:      from,
				To:        to,
				Increment: increment,
			}
		case strings.HasPrefix(ct, "application/json"):
			if err := c.Bind(input); err != nil {
				return serverError(c, "Error parsing request", err)
			}
		default:
			return notAcceptable(c)
		}

		experiment, err := ss.ScheduleExperiment(input)
		if err != nil {
			return serverError(c, "Error scheduling experiment", err)
		}

		return c.JSON(http.StatusOK, buildExperimentResource(c, experiment))
	}
}

func findExperiment(ss solver.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := solver.ExperimentID(c.Param("experimentID"))
		e, err := ss.FindExperiment(id, nil)
		if err != nil {
			if err == solver.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Experiment does not exist"})
			} else {
				return serverError(c, "Error loading experiment", err)
			}
		}

		return c.JSON(http.StatusOK, buildExperimentResource(c, e))
	}
}
