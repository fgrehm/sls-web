package compiler

import (
	"fmt"

	"github.com/fgrehm/go-san"
	"github.com/fgrehm/go-san/model"
	"github.com/mitchellh/hashstructure"
)

type ParseResult struct {
	ParsedModel     *sanmodel.Model `json:"parsedModel"`
	ModelHash       string          `json:"modelHash"`
	TransitionsHash string          `json:"transitionsHash"`
}

type Service interface {
	Parse(source []byte) (*ParseResult, error)
	Compile(m *sanmodel.Model) ([]byte, error)
	Hash(modelOrNetwork interface{}) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Parse(source []byte) (*ParseResult, error) {
	model, err := san.Parse(source)
	if err != nil {
		return nil, err
	}

	modelHash, err := s.Hash(model)
	if err != nil {
		return nil, err
	}

	transitionsHash, err := s.Hash(model.Network)
	if err != nil {
		return nil, err
	}

	return &ParseResult{
		ParsedModel:     model,
		ModelHash:       modelHash,
		TransitionsHash: transitionsHash,
	}, nil
}

func (s *service) Compile(m *sanmodel.Model) ([]byte, error) {
	return san.Compile(m)
}

func (s *service) Hash(modelOrNetwork interface{}) (string, error) {
	hash, err := hashstructure.Hash(modelOrNetwork, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
}
