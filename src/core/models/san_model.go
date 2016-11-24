package models

import (
	"errors"
	"time"

	"github.com/fgrehm/go-san/model"

	"github.com/pborman/uuid"
)

type SanModelID string

func NextSanModelID() SanModelID {
	return SanModelID(uuid.New())
}

type SanModel struct {
	ID              SanModelID
	Name            string
	Title           string
	Description     string
	Source          string
	Hash            string
	TransitionsHash string
	ParsedModel     *sanmodel.Model
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Repository interface {
	Store(model *SanModel) error
	Find(id SanModelID) (*SanModel, error)
	All() ([]*SanModel, error)
}

// ErrUnknown is used when a model could not be found.
var ErrUnknown = errors.New("unknown model")
