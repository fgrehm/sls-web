package api

import (
	"bytes"
	"net/http"

	"github.com/fgrehm/sls-web/src/core/models"
	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/labstack/echo"
)

func registerSanModelsEndpoints(e *echo.Echo, ms models.Service) {
	e.GET("/models", listSanModels(ms))
	e.POST("/models", createSanModel(ms))
	e.PUT("/models/:modelID", updateSanModel(ms))
	e.GET("/models/:modelID", findSanModel(ms))
	e.GET("/models/:modelID/graph/*any", renderSanModelGraph(ms))
	e.POST("/models/:modelID/solutions", scheduleModelSolution(ms))
	e.GET("/models/:modelID/solutions/:solutionID", findModelSolution(ms))
	e.POST("/models/:modelID/experiments", scheduleModelExperiment(ms))
	e.GET("/models/:modelID/experiments/:experimentID", findModelExperiment(ms))
}

func listSanModels(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		models, err := ms.All()
		if err != nil {
			return serverError(c, "Error loading SAN Models", err)
		}
		return c.JSON(http.StatusOK, buildSanModelsCollectionResource(models, c))
	}
}

func createSanModel(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		input, err := parseSanModelUpsertInput(c)
		if err != nil {
			return err
		}

		id, err := ms.Create(input)
		if err != nil {
			return serverError(c, "Error creating SAN Model", err)
		}

		m, err := ms.Find(id)
		if err != nil {
			return serverError(c, "Error loading SAN Model", err)
		}

		return c.JSON(http.StatusOK, buildSanModelResource(c, m))
	}
}

func updateSanModel(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := models.SanModelID(c.Param("modelID"))

		input, err := parseSanModelUpsertInput(c)
		if err != nil {
			return err
		}

		m, err := ms.Update(id, input)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			} else {
				return serverError(c, "Error updating SAN Model", err)
			}
		}

		return c.JSON(http.StatusOK, buildSanModelResource(c, m))
	}
}

func findSanModel(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := models.SanModelID(c.Param("modelID"))
		m, err := ms.Find(id)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			} else {
				return serverError(c, "Error loading SAN Model", err)
			}
		}

		return c.JSON(http.StatusOK, buildSanModelResource(c, m))
	}
}

func renderSanModelGraph(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := models.SanModelID(c.Param("modelID"))

		outputFile := &bytes.Buffer{}
		err := ms.RenderGraph(id, outputFile)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			}

			return serverError(c, "Error rendering SAN Model graph", err)
		}

		return c.Blob(http.StatusOK, "image/png", outputFile.Bytes())
	}
}

func scheduleModelSolution(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := models.SanModelID(c.Param("modelID"))

		si := &models.ModelSolutionInput{}
		if err := c.Bind(si); err != nil {
			return err
		}

		solution, err := ms.ScheduleSolution(id, si)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			}

			return serverError(c, "Error scheduling SAN Model solution", err)
		}

		return c.JSON(http.StatusOK, buildSolutionResource(c, solution))
	}
}

func findModelSolution(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		modelID := models.SanModelID(c.Param("modelID"))
		solutionID := solver.SolutionID(c.Param("solutionID"))
		solution, err := ms.FindSolution(modelID, solutionID)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			} else if err == solver.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model solution does not exist"})
			} else {
				return serverError(c, "Error loading SAN Model", err)
			}
		}

		return c.JSON(http.StatusOK, buildSolutionResource(c, solution))
	}
}

func scheduleModelExperiment(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := models.SanModelID(c.Param("modelID"))

		ei := &models.ModelExperimentInput{}
		if err := c.Bind(ei); err != nil {
			return err
		}

		experiment, err := ms.ScheduleExperiment(id, ei)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			}

			return serverError(c, "Error scheduling SAN Model experiment", err)
		}

		return c.JSON(http.StatusOK, buildExperimentResource(c, experiment))
	}
}

func findModelExperiment(ms models.Service) func(c echo.Context) error {
	return func(c echo.Context) error {
		modelID := models.SanModelID(c.Param("modelID"))
		experimentID := solver.ExperimentID(c.Param("experimentID"))
		experiment, err := ms.FindExperiment(modelID, experimentID)
		if err != nil {
			if err == models.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model does not exist"})
			} else if err == solver.ErrUnknown {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "SAN Model experiment does not exist"})
			} else {
				return serverError(c, "Error loading SAN Model", err)
			}
		}

		return c.JSON(http.StatusOK, buildExperimentResource(c, experiment))
	}
}
