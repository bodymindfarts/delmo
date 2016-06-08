package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type Suite struct {
	config                *Config
	globalTaskEnvironment TaskEnvironment
}

func NewSuite(config *Config, globalTaskEnvironment TaskEnvironment) (*Suite, error) {
	suite := &Suite{
		config:                config,
		globalTaskEnvironment: globalTaskEnvironment,
	}
	return suite, nil
}

type BuildOutput struct {
	ui cli.Ui
}

func (o *BuildOutput) Write(p []byte) (int, error) {
	o.ui.Output(strings.TrimSpace(string(p)))
	return len(p), nil
}

func (s *Suite) Run(ui cli.Ui) int {
	s.initializeSystem(ui)
	ui.Info(fmt.Sprintf("\nRunning Test Suite for System %s", s.config.Suite.Name))

	failed := []*TestReport{}
	succeeded := []*TestReport{}
	for _, test := range s.config.Tests {
		runner := NewTestRunner(test, s.config.Tasks, s.globalTaskEnvironment)
		runtime, err := NewDockerCompose(s.config.Suite.System, test.Name)
		if err != nil {
			ui.Error(fmt.Sprintf("Error creating runtime! %s", err))
			return 1
		}

		ui.Info(fmt.Sprintf("Running test %s", test.Name))
		report := runner.RunTest(runtime, ui)
		if report.Success {
			succeeded = append(succeeded, report)
			runner.Cleanup()
			ui.Info(fmt.Sprintf("Test %s Succeeded!", test.Name))
		} else {
			failed = append(succeeded, report)
			ui.Info(fmt.Sprintf("Test %s Failed!\nRuntime Output:\n%s", test.Name, report.SystemOutput()))
		}
	}

	outputSummary(ui, failed, succeeded)
	if len(failed) != 0 {
		return 1
	}
	return 0
}

func outputSummary(ui cli.Ui, failed []*TestReport, succeeded []*TestReport) {
	ui.Output(
		fmt.Sprintf("\n\nSUMMARY:\n%d tests succeeded\n%d tests failed",
			len(succeeded),
			len(failed)))
}

func (s *Suite) initializeSystem(ui cli.Ui) error {
	dc, err := NewDockerCompose(s.config.Suite.System, "")
	output := &BuildOutput{ui: ui}
	if err != nil {
		return fmt.Errorf("Could not initialize docker-compose\n%s", err)
	}

	ui.Info(fmt.Sprintf("Pulling images for system %s", s.config.Suite.Name))
	err = dc.Pull(output)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error pulling images\n%s", err))
	}

	ui.Info(fmt.Sprintf("Builing images for system %s", s.config.Suite.Name))
	if s.config.Suite.OnlyBuildTask {
		err = dc.Build(output, s.config.Suite.TaskService)
	} else {
		err = dc.Build(output)
	}
	if err != nil {
		return fmt.Errorf("Could not build system\n%s", err)
	}

	return nil
}
