package api

import (
	"time"

	"github.com/fgrehm/sls-web/src/core/models"
	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/labstack/echo"
)

type SanModelResource struct {
	Url             string      `json:"url"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Source          string      `json:"source"`
	Hash            string      `json:"hash"`
	TransitionsHash string      `json:"transitionsHash"`
	ParsedModel     interface{} `json:"parsedModel"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	GraphUrl        string      `json:"graphUrl"`
}

type SanModelSummaryResource struct {
	Url       string    `json:"url"`
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SanModelsCollectionResource []*SanModelSummaryResource

func buildSanModelResource(c echo.Context, m *models.SanModel) *SanModelResource {
	return &SanModelResource{
		Url:             sanModelUrl(c, m.ID),
		ID:              string(m.ID),
		Title:           m.Title,
		Description:     m.Description,
		Source:          m.Source,
		Hash:            m.Hash,
		TransitionsHash: m.TransitionsHash,
		ParsedModel:     m.ParsedModel,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		GraphUrl:        sanModelGraphUrl(c, m.ID, m.TransitionsHash),
	}
}

func buildSanModelsCollectionResource(models []*models.SanModel, c echo.Context) SanModelsCollectionResource {
	col := SanModelsCollectionResource{}
	for _, m := range models {
		col = append(col, &SanModelSummaryResource{
			Url:       sanModelUrl(c, m.ID),
			ID:        string(m.ID),
			Title:     m.Title,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return col
}

type SolutionResource struct {
	Url               string                     `json:"url"`
	ID                string                     `json:"id"`
	Source            string                     `json:"source"`
	MaxIterations     uint                       `json:"maxIterations"`
	Tolerance         float64                    `json:"tolerance"`
	Errored           *bool                      `json:"errored"`
	Found             *bool                      `json:"found"`
	Steps             *uint64                    `json:"steps"`
	ExecutionTime     *float64                   `json:"executionTime"`
	Results           *solver.IntegrationResults `json:"results"`
	CustomIdentifiers solver.CustomIdentifiers   `json:"customIdentifiers"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
	FinishedAt        *time.Time                 `json:"finishedAt"`
}

func buildSolutionResource(c echo.Context, s *solver.Solution) *SolutionResource {
	return &SolutionResource{
		Url:               solutionUrl(c, s.ID),
		ID:                string(s.ID),
		Source:            s.Source,
		MaxIterations:     s.MaxIterations,
		Tolerance:         s.Tolerance,
		Errored:           s.Errored,
		Found:             s.Found,
		Steps:             s.Steps,
		ExecutionTime:     s.ExecutionTime,
		Results:           s.Results,
		CustomIdentifiers: s.CustomIdentifiers,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		FinishedAt:        s.FinishedAt,
	}
}

type ExperimentResource struct {
	Url                  string                         `json:"url"`
	ID                   string                         `json:"id"`
	Source               string                         `json:"source"`
	MaxIterations        uint                           `json:"maxIterations"`
	Tolerance            float64                        `json:"tolerance"`
	Identifier           ExperimentIdentifier           `json:"identifier"`
	IntegrationFunctions ExperimentIntegrationFunctions `json:"integrationFunctions"`
	Solutions            ExperimentSolutions            `json:"solutions"`
	CreatedAt            time.Time                      `json:"createdAt"`
}

type ExperimentIdentifier struct {
	Name      string  `json:"name"`
	From      float64 `json:"from"`
	To        float64 `json:"to"`
	Increment float64 `json:"increment"`
}

type ExperimentIntegrationFunctions map[string]string

type ExperimentSolution struct {
	ID              string                     `json:"id"`
	IdentifierValue float64                    `json:"identifierValue"`
	FinishedAt      *time.Time                 `json:"finishedAt"`
	Found           *bool                      `json:"found"`
	Errored         *bool                      `json:"errored"`
	Steps           *uint64                    `json:"steps"`
	ExecutionTime   *float64                   `json:"executionTime"`
	Results         *solver.IntegrationResults `json:"results"`
}

type ExperimentSolutions []*ExperimentSolution

func buildExperimentResource(c echo.Context, e *solver.Experiment) *ExperimentResource {
	integrationFunctions := ExperimentIntegrationFunctions{}
	for _, f := range e.IntegrationFunctions {
		integrationFunctions[f.Label] = f.Expression
	}

	identifier := ExperimentIdentifier{
		Name:      e.Identifier.Name,
		From:      e.Identifier.From,
		To:        e.Identifier.To,
		Increment: e.Identifier.Increment,
	}

	return &ExperimentResource{
		Url:                  experimentUrl(c, e.ID),
		ID:                   string(e.ID),
		Source:               e.Source,
		MaxIterations:        e.MaxIterations,
		Tolerance:            e.Tolerance,
		Identifier:           identifier,
		IntegrationFunctions: integrationFunctions,
		Solutions:            buildExperimentSolutions(c, e),
		CreatedAt:            e.CreatedAt,
	}
}

func buildExperimentSolutions(c echo.Context, e *solver.Experiment) ExperimentSolutions {
	if e.Solutions == nil {
		return nil
	}
	solutions := ExperimentSolutions{}
	for _, s := range *e.Solutions {
		solutions = append(solutions, buildExperimentSolution(c, s))
	}
	return solutions
}

func buildExperimentSolution(c echo.Context, s *solver.ExperimentSolution) *ExperimentSolution {
	return &ExperimentSolution{
		ID:              string(s.Solution.ID),
		IdentifierValue: s.IdentifierValue,
		FinishedAt:      s.Solution.FinishedAt,
		Found:           s.Solution.Found,
		Errored:         s.Solution.Errored,
		Steps:           s.Solution.Steps,
		ExecutionTime:   s.Solution.ExecutionTime,
		Results:         s.Solution.Results,
	}
}
