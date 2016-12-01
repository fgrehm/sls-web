package renderer

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/fgrehm/go-san/model"
	"github.com/fgrehm/sls-web/src/core/compiler"

	"github.com/awalterschulze/gographviz"
)

type Service interface {
	Render(m *sanmodel.Model) ([]byte, error)
	RenderFromSource(source []byte) ([]byte, error)
}

type service struct {
	cs compiler.Service
}

func (s *service) Render(m *sanmodel.Model) ([]byte, error) {
	dot, err := s.graph(m)
	if err != nil {
		return nil, err
	}

	return s.png(dot)
}

func (s *service) RenderFromSource(source []byte) ([]byte, error) {
	parseResult, err := s.cs.Parse(source)
	if err != nil {
		return nil, err
	}

	return s.Render(parseResult.ParsedModel)
}

func (s *service) graph(m *sanmodel.Model) (string, error) {
	g := gographviz.NewGraph()
	g.SetName(m.Network.Name)
	g.SetDir(true)

	for i := 0; i < len(m.Network.Automata); i++ {
		aut := m.Network.Automata[i]
		g.AddSubGraph(m.Network.Name, "cluster_"+aut.Name, map[string]string{
			"label": aut.Name,
		})
		for _, state := range s.extractStates(aut) {
			g.AddNode("cluster_"+aut.Name, aut.Name+"_"+state, map[string]string{
				"label": state,
			})
		}
		for _, transition := range aut.Transitions {
			for _, event := range transition.Events {
				g.AddEdge(aut.Name+"_"+transition.From, aut.Name+"_"+transition.To, true, map[string]string{
					"label": event.EventName,
				})
			}
		}
	}
	return g.String(), nil
}

func (s *service) extractStates(aut *sanmodel.Automaton) []string {
	statesMap := make(map[string]*struct{})
	states := []string{}
	for _, transition := range aut.Transitions {
		if statesMap[transition.From] == nil {
			statesMap[transition.From] = &struct{}{}
			states = append(states, transition.From)
		}
		if statesMap[transition.To] == nil {
			statesMap[transition.To] = &struct{}{}
			states = append(states, transition.To)
		}
	}
	return states
}

func (s *service) png(dotGraph string) ([]byte, error) {
	cmd := exec.Command("dot", "-T", "png")
	cmd.Stdin = strings.NewReader(dotGraph)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		println(out.String())
		return nil, err
	}
	return out.Bytes(), nil
}

func NewService(cs compiler.Service) Service {
	return &service{cs}
}
