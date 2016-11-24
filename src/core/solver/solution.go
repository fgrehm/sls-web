package solver

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

type SolutionID string

func NextSolutionID() SolutionID {
	return SolutionID(uuid.New())
}

type Solution struct {
	ID                SolutionID
	ExternalKey       *string
	Source            string
	ModelHash         string
	CustomIdentifiers CustomIdentifiers
	MaxIterations     uint
	Tolerance         float64
	Errored           *bool
	Found             *bool
	Steps             *uint64
	ExecutionTime     *float64
	Results           *IntegrationResults
	CreatedAt         time.Time
	UpdatedAt         time.Time
	FinishedAt        *time.Time
}

type CustomIdentifiers map[string]interface{}

type IntegrationResult struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}
type IntegrationResults []IntegrationResult

type SolutionRepository interface {
	Store(solution *Solution) error
	Find(id SolutionID, externalKey *string) (*Solution, error)
	FindForReuse(externalKey, modelHash string, maxIterations uint, tolerance float64) (*Solution, error)
}

// ErrUnknown is used when a solution could not be found.
var ErrUnknown = errors.New("unknown solution")
