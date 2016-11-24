package solver

import (
	"time"

	"github.com/fgrehm/go-san/model"

	"github.com/fgrehm/sls-web/src/core/compiler"
)

type SolutionInput struct {
	ExternalKey       *string
	MaxIterations     *uint
	Tolerance         *float64
	Source            string
	CustomIdentifiers CustomIdentifiers
}

type ExperimentInput struct {
	ExternalKey   *string
	MaxIterations *uint
	Tolerance     *float64
	Source        string
	Identifier    ExperimentIdentifier
}

type Scheduler interface {
	ScheduleSolution(*Solution) error
	ScheduleExperiment(*Experiment) error
}

type Service interface {
	ScheduleSolution(si *SolutionInput) (*Solution, error)
	Solve(id SolutionID) error
	FindSolution(id SolutionID, externalKey *string) (*Solution, error)
	ScheduleExperiment(ei *ExperimentInput) (*Experiment, error)
	ProcessExperiment(id ExperimentID) error
	FindExperiment(id ExperimentID, externalKey *string) (*Experiment, error)
}

const (
	DEFAULT_TOLERANCE      = 1e-10
	DEFAULT_MAX_ITERATIONS = uint(1000000)
)

type service struct {
	solutions   SolutionRepository
	experiments ExperimentRepository
	scheduler   Scheduler
	compiler    compiler.Service
}

func (s *service) ScheduleSolution(si *SolutionInput) (*Solution, error) {
	tolerance := DEFAULT_TOLERANCE
	if si.Tolerance != nil {
		tolerance = *si.Tolerance
	}
	maxIterations := DEFAULT_MAX_ITERATIONS
	if si.MaxIterations != nil {
		maxIterations = *si.MaxIterations
	}

	parseResult, err := s.compiler.Parse([]byte(si.Source))
	if err != nil {
		return nil, err
	}

	modelHash, err := s.applySolutionInput(parseResult.ParsedModel, si)
	if err != nil {
		return nil, err
	}
	source, err := s.generateSource(parseResult.ParsedModel, si)
	if err != nil {
		return nil, err
	}
	if modelHash == "" {
		modelHash = parseResult.ModelHash
	}

	if si.ExternalKey != nil {
		solution, err := s.solutions.FindForReuse(*si.ExternalKey, modelHash, maxIterations, tolerance)
		if err == nil {
			return solution, nil
		}
	}

	now := time.Now()
	solution := &Solution{
		ID:                NextSolutionID(),
		ExternalKey:       si.ExternalKey,
		ModelHash:         modelHash,
		Source:            source,
		MaxIterations:     maxIterations,
		Tolerance:         tolerance,
		CustomIdentifiers: si.CustomIdentifiers,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.solutions.Store(solution); err != nil {
		return nil, err
	}

	if err := s.scheduler.ScheduleSolution(solution); err != nil {
		return nil, err
	}

	return solution, nil
}

func (s *service) Solve(id SolutionID) error {
	solution, err := s.FindSolution(id, nil)
	if err != nil {
		return err
	}

	res, solutionError := sanLiteSolver([]byte(solution.Source), solution.MaxIterations, solution.Tolerance)
	now := time.Now()

	errored := (solutionError != nil)
	solution.Errored = &errored
	solution.UpdatedAt = now
	solution.FinishedAt = &now
	if !errored {
		solution.Found = &res.found
		solution.Steps = &res.steps
		solution.ExecutionTime = &res.executionTime
		solution.Results = &res.results
	}
	if err := s.solutions.Store(solution); err != nil {
		return err
	}

	// time.Sleep(5 * time.Second)
	return solutionError
}

func (s *service) FindSolution(id SolutionID, externalKey *string) (*Solution, error) {
	return s.solutions.Find(id, externalKey)
}

func (s *service) ScheduleExperiment(ei *ExperimentInput) (*Experiment, error) {
	tolerance := DEFAULT_TOLERANCE
	if ei.Tolerance != nil {
		tolerance = *ei.Tolerance
	}
	maxIterations := DEFAULT_MAX_ITERATIONS
	if ei.MaxIterations != nil {
		maxIterations = *ei.MaxIterations
	}

	parseResult, err := s.compiler.Parse([]byte(ei.Source))
	if err != nil {
		return nil, err
	}
	integrationFunctions := ExperimentIntegrationFunctions{}
	for _, res := range parseResult.ParsedModel.Results {
		integrationFunctions = append(integrationFunctions, ExperimentIntegrationFunction{
			Label:      res.Label,
			Expression: res.Expression,
		})
	}

	now := time.Now()
	experiment := &Experiment{
		ID:                   NextExperimentID(),
		ExternalKey:          ei.ExternalKey,
		Source:               ei.Source,
		MaxIterations:        maxIterations,
		Tolerance:            tolerance,
		Identifier:           ei.Identifier,
		IntegrationFunctions: integrationFunctions,
		CreatedAt:            now,
	}

	if err := s.experiments.Store(experiment); err != nil {
		return nil, err
	}

	if err := s.scheduler.ScheduleExperiment(experiment); err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *service) ProcessExperiment(id ExperimentID) error {
	experiment, err := s.FindExperiment(id, nil)
	if err != nil {
		return err
	}

	if experiment.Solutions == nil {
		experiment.Solutions = &ExperimentSolutions{}
	}
	solutions := *experiment.Solutions

	i := experiment.Identifier
	for val := i.From; val <= i.To; val += i.Increment {
		solution, err := s.ScheduleSolution(&SolutionInput{
			ExternalKey:   experiment.ExternalKey,
			MaxIterations: &experiment.MaxIterations,
			Tolerance:     &experiment.Tolerance,
			Source:        experiment.Source,
			CustomIdentifiers: map[string]interface{}{
				i.Name: val,
			},
		})
		if err != nil {
			return err
		}
		experimentSolution := &ExperimentSolution{
			ID:              solution.ID,
			IdentifierValue: val,
			Solution:        solution,
		}
		solutions = append(solutions, experimentSolution)
	}
	now := time.Now()
	experiment.Solutions = &solutions
	experiment.SolutionsScheduledAt = &now
	return s.experiments.Store(experiment)
}

func (s *service) FindExperiment(id ExperimentID, externalKey *string) (*Experiment, error) {
	return s.experiments.Find(id, externalKey)
}

func (s *service) generateSource(m *sanmodel.Model, si *SolutionInput) (string, error) {
	hasCustomIdentifiers := false
	for _, val := range si.CustomIdentifiers {
		if val != nil {
			hasCustomIdentifiers = true
			break
		}
	}

	if !hasCustomIdentifiers {
		return si.Source, nil
	}

	src, err := s.compiler.Compile(m)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func (s *service) applySolutionInput(m *sanmodel.Model, si *SolutionInput) (string, error) {
	hasCustomIdentifiers := false
	for _, val := range si.CustomIdentifiers {
		if val != nil {
			hasCustomIdentifiers = true
			break
		}
	}

	if !hasCustomIdentifiers {
		return "", nil
	}

	for _, id := range m.Identifiers {
		if si.CustomIdentifiers[id.Name] != nil {
			id.Value = si.CustomIdentifiers[id.Name]
		}
	}

	return s.compiler.Hash(m)
}

func NewService(solutions SolutionRepository, experiments ExperimentRepository, scheduler Scheduler, compiler compiler.Service) Service {
	return &service{solutions, experiments, scheduler, compiler}
}
