package solver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"github.com/kardianos/osext"
)

var runner *slsRunner
var WORKDIR = "/tmp/sls/solutions"

// TODO: This code should not be on an initializer
func init() {
	exe := os.Getenv("SLS_PATH")
	if exe == "" {
		exePath, err := osext.ExecutableFolder()
		if err == nil {
			exe = exePath + "/san-lite-solver"
		}
	}
	if exe == "" || !fileExists(exe) {
		panic("SAN Lite-Solver executable not found")
	}

	runner = &slsRunner{
		executablePath: exe,
		workdir:        WORKDIR,
	}
	if err := os.MkdirAll(WORKDIR, 0777); err != nil {
		panic(err)
	}
}

type slsRunner struct {
	source         []byte
	executablePath string
	workdir        string
}

type slsSolutionResult struct {
	found         bool
	steps         uint64
	executionTime float64
	results       IntegrationResults
}

func sanLiteSolver(source []byte, maxIterations uint, tolerance float64) (*slsSolutionResult, error) {
	return runner.solve(source, maxIterations, tolerance)
}

func (r *slsRunner) solve(source []byte, maxIterations uint, tolerance float64) (*slsSolutionResult, error) {
	modelPath := path.Join(r.workdir, fmt.Sprintf("%s-%s", time.Now().Format("20060102-150405.000"), uuid.New()))
	logPath := path.Join(modelPath, "san-lite-solver.log")

	if err := os.MkdirAll(modelPath, 0777); err != nil {
		return nil, err
	}

	started := time.Now()
	sanFilePath, err := r.createSanFile(modelPath, source)
	if err != nil {
		return nil, err
	}
	output, err := r.run(modelPath, sanFilePath, maxIterations, tolerance)
	if err != nil {
		fmt.Println(output.String())
		ioutil.WriteFile(logPath, output.Bytes(), 0666)
		return nil, err
	}
	if err := ioutil.WriteFile(logPath, output.Bytes(), 0666); err != nil {
		return nil, err
	}

	solution, err := r.parseResultsFromOutput(output)
	if err != nil {
		return nil, err
	}
	solution.executionTime = float64(time.Since(started)) / float64(time.Millisecond)
	return solution, nil
}

func (r *slsRunner) createSanFile(modelPath string, source []byte) (string, error) {
	outputPath := path.Join(modelPath, "model.san")
	if err := ioutil.WriteFile(outputPath, source, 0666); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (r *slsRunner) run(workDir, modelPath string, maxIterations uint, tolerance float64) (*bytes.Buffer, error) {
	toleranceStr := fmt.Sprintf("-tol=%e", tolerance)
	iterationsStr := fmt.Sprintf("-ite=%d", maxIterations)
	cmd := exec.Command(r.executablePath, modelPath, iterationsStr, toleranceStr)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Dir = workDir
	err := cmd.Run()
	return out, err
}

func (r *slsRunner) parseResultsFromOutput(output *bytes.Buffer) (*slsSolutionResult, error) {
	solution := &slsSolutionResult{results: IntegrationResults{}, found: true}
	var integrationLine = regexp.MustCompile(`^Integration of function ([a-zA-Z0-9_-]+) = ([0-9e\.\-\+]+)`)
	var stationarySolutionLine = regexp.MustCompile(`^Stationary solution (NOT )?found in ([a-zA-Z0-9_-]+) steps`)
	for {
		line, err := output.ReadString('\n')

		submatches := integrationLine.FindStringSubmatch(line)
		if len(submatches) == 3 {
			value, err := strconv.ParseFloat(submatches[2], 10)
			if err != nil {
				return nil, err
			}
			solution.results = append(solution.results, IntegrationResult{
				Label: submatches[1],
				Value: value,
			})
			continue
		}

		submatches = stationarySolutionLine.FindStringSubmatch(line)
		if len(submatches) > 0 {
			if submatches[1] == "NOT " {
				solution.found = false
			}
			value, err := strconv.ParseUint(submatches[len(submatches)-1], 10, 64)
			if err != nil {
				return nil, err
			}
			solution.steps = value
		}

		if err == io.EOF {
			break
		}
	}
	return solution, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}
