package models

import (
	"io"
	"time"

	"github.com/fgrehm/sls-web/src/core/compiler"
	"github.com/fgrehm/sls-web/src/core/renderer"
	"github.com/fgrehm/sls-web/src/core/solver"
)

type UpsertInput struct {
	Source      *string `json:"source"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type ModelSolutionInput struct {
	MaxIterations     *uint                  `json:"maxIterations,omitempty"`
	Tolerance         *float64               `json:"tolerance,omitempty"`
	CustomIdentifiers map[string]interface{} `json:"customIdentifiers,omitempty"`
}

type ModelExperimentInput struct {
	IdentifierParams struct {
		Name      string  `json:"name"`
		From      float64 `json:"from"`
		To        float64 `json:"to"`
		Increment float64 `json:"increment"`
	} `json:"identifier"`
	MaxIterations *uint    `json:"maxIterations"`
	Tolerance     *float64 `json:"tolerance"`
}

type Service interface {
	Create(input *UpsertInput) (SanModelID, error)
	Find(id SanModelID) (*SanModel, error)
	Update(id SanModelID, input *UpsertInput) (*SanModel, error)
	All() ([]*SanModel, error)
	RenderGraph(id SanModelID, out io.Writer) error
	ScheduleSolution(id SanModelID, req *ModelSolutionInput) (*solver.Solution, error)
	FindSolution(id SanModelID, solutionID solver.SolutionID) (*solver.Solution, error)
	ScheduleExperiment(id SanModelID, req *ModelExperimentInput) (*solver.Experiment, error)
	FindExperiment(id SanModelID, experimentID solver.ExperimentID) (*solver.Experiment, error)
}

type service struct {
	repo     Repository
	compiler compiler.Service
	renderer renderer.Service
	solver   solver.Service
}

func (s *service) Create(input *UpsertInput) (SanModelID, error) {
	now := time.Now()

	m := &SanModel{ID: NextSanModelID(), CreatedAt: now, UpdatedAt: now}

	if input.Title != nil {
		m.Title = *input.Title
	}
	if input.Description != nil {
		m.Description = *input.Description
	}

	if input.Source != nil {
		parseResult, err := s.compiler.Parse([]byte(*input.Source))
		if err != nil {
			return SanModelID(""), err
		}
		m.Source = *input.Source
		m.ParsedModel = parseResult.ParsedModel
		m.Hash = parseResult.ModelHash
		m.TransitionsHash = parseResult.TransitionsHash
	}

	err := s.repo.Store(m)
	if err != nil {
		return SanModelID(""), err
	}

	return m.ID, nil
}

func (s *service) Find(id SanModelID) (*SanModel, error) {
	return s.repo.Find(id)
}

func (s *service) Update(id SanModelID, input *UpsertInput) (*SanModel, error) {
	now := time.Now()

	m, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}

	m.UpdatedAt = now
	if input.Title != nil {
		m.Title = *input.Title
	}
	if input.Description != nil {
		m.Description = *input.Description
	}

	if input.Source != nil {
		parseResult, err := s.compiler.Parse([]byte(*input.Source))
		if err != nil {
			return nil, err
		}
		m.Source = *input.Source
		m.ParsedModel = parseResult.ParsedModel
		m.Hash = parseResult.ModelHash
		m.TransitionsHash = parseResult.TransitionsHash
	}

	if err = s.repo.Store(m); err != nil {
		return nil, err
	}
	return m, err
}

func (s *service) All() ([]*SanModel, error) {
	return s.repo.All()
}

func (s *service) RenderGraph(id SanModelID, out io.Writer) error {
	m, err := s.Find(id)
	if err != nil {
		return err
	}

	img, err := s.renderer.RenderFromSource([]byte(m.Source))
	if err != nil {
		return err
	}
	_, err = out.Write(img)
	return err
}

func (s *service) ScheduleSolution(id SanModelID, req *ModelSolutionInput) (*solver.Solution, error) {
	m, err := s.Find(id)
	if err != nil {
		return nil, err
	}

	externalKey := string(m.ID)
	solutionReq := &solver.SolutionInput{
		ExternalKey:       &externalKey,
		MaxIterations:     req.MaxIterations,
		Tolerance:         req.Tolerance,
		Source:            m.Source,
		CustomIdentifiers: req.CustomIdentifiers,
	}

	solution, err := s.solver.ScheduleSolution(solutionReq)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

func (s *service) FindSolution(modelID SanModelID, solutionID solver.SolutionID) (*solver.Solution, error) {
	m, err := s.Find(modelID)
	if err != nil {
		return nil, err
	}

	externalKey := string(m.ID)
	return s.solver.FindSolution(solutionID, &externalKey)
}

func (s *service) ScheduleExperiment(id SanModelID, req *ModelExperimentInput) (*solver.Experiment, error) {
	m, err := s.Find(id)
	if err != nil {
		return nil, err
	}

	externalKey := string(m.ID)
	experimentReq := &solver.ExperimentInput{
		ExternalKey:   &externalKey,
		MaxIterations: req.MaxIterations,
		Tolerance:     req.Tolerance,
		Source:        m.Source,
		Identifier: solver.ExperimentIdentifier{
			Name:      req.IdentifierParams.Name,
			From:      req.IdentifierParams.From,
			To:        req.IdentifierParams.To,
			Increment: req.IdentifierParams.Increment,
		},
	}

	experiment, err := s.solver.ScheduleExperiment(experimentReq)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *service) FindExperiment(modelID SanModelID, experimentID solver.ExperimentID) (*solver.Experiment, error) {
	m, err := s.Find(modelID)
	if err != nil {
		return nil, err
	}

	externalKey := string(m.ID)
	return s.solver.FindExperiment(experimentID, &externalKey)
}

func NewService(repo Repository, compiler compiler.Service, renderer renderer.Service, solver solver.Service) Service {
	return &service{repo, compiler, renderer, solver}
}
