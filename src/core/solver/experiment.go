package solver

import (
	"time"

	"github.com/pborman/uuid"
)

func NextExperimentID() ExperimentID {
	return ExperimentID(uuid.New())
}

type ExperimentID string

type Experiment struct {
	ID                   ExperimentID
	ExternalKey          *string
	Source               string
	MaxIterations        uint
	Tolerance            float64
	Identifier           ExperimentIdentifier
	IntegrationFunctions ExperimentIntegrationFunctions
	SolutionsScheduledAt *time.Time
	Solutions            *ExperimentSolutions
	CreatedAt            time.Time
}

type ExperimentIdentifier struct {
	Name      string
	From      float64
	To        float64
	Increment float64
}

type ExperimentIntegrationFunction struct {
	Label      string
	Expression string
}

type ExperimentIntegrationFunctions []ExperimentIntegrationFunction

type ExperimentSolution struct {
	ID              SolutionID
	IdentifierValue float64
	Solution        *Solution `bson:"-"`
}

type ExperimentSolutions []*ExperimentSolution

type ExperimentRepository interface {
	Store(experiment *Experiment) error
	Find(id ExperimentID, externalKey *string) (*Experiment, error)
}
