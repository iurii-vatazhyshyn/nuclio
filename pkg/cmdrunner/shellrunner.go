/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmdrunner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/nuclio/nuclio/pkg/common"

	"github.com/nuclio/errors"
	"github.com/nuclio/logger"
)

type ShellRunner struct {
	logger logger.Logger
	shell  string
}

func NewShellRunner(parentLogger logger.Logger) (*ShellRunner, error) {
	return &ShellRunner{
		logger: parentLogger.GetChild("runner"),
		shell:  "/bin/sh",
	}, nil
}

func (sr *ShellRunner) Run(runOptions *RunOptions, format string, vars ...interface{}) (RunResult, error) {

	// support missing runOptions for tests that send nil
	if runOptions == nil {
		runOptions = &RunOptions{}
	}

	// format the command
	formattedCommand := fmt.Sprintf(format, vars...)
	redactedCommand := common.Redact(runOptions.LogRedactions, formattedCommand)
	sr.logger.DebugWith("Executing", "command", redactedCommand)

	// create a command
	cmd := exec.Command(sr.shell, "-c", formattedCommand)

	// if there are runOptions, set them
	if runOptions.WorkingDir != nil {
		cmd.Dir = *runOptions.WorkingDir
	}

	// get environment variables if any
	if runOptions.Env != nil {
		cmd.Env = sr.getEnvFromOptions(runOptions)
	}

	if runOptions.Stdin != nil {
		cmd.Stdin = strings.NewReader(*runOptions.Stdin)
	}

	runResult := RunResult{
		ExitCode: 0,
	}

	err := sr.runAndCaptureOutput(cmd, runOptions, &runResult)

	if err != nil {
		var exitCode int

		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		}

		runResult.ExitCode = exitCode

		sr.logger.DebugWith("Failed to execute command",
			"output", runResult.Output,
			"stderr", runResult.Stderr,
			"exitCode", runResult.ExitCode,
			"err", err)

		err = errors.Wrapf(err, "stdout:\n%s\nstderr:\n%s", runResult.Output, runResult.Stderr)

		return runResult, err
	}

	sr.logger.DebugWith("Command executed successfully",
		"output", runResult.Output,
		"stderr", runResult.Stderr,
		"exitCode", runResult.ExitCode)

	return runResult, nil
}

func (sr *ShellRunner) SetShell(shell string) {
	sr.shell = shell
}

func (sr *ShellRunner) getEnvFromOptions(runOptions *RunOptions) []string {
	envs := []string{}

	for name, value := range runOptions.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", name, value))
	}

	return envs
}

func (sr *ShellRunner) runAndCaptureOutput(cmd *exec.Cmd,
	runOptions *RunOptions,
	runResult *RunResult) error {

	switch runOptions.CaptureOutputMode {

	case CaptureOutputModeCombined:
		stdoutAndStderr, err := cmd.CombinedOutput()
		runResult.Output = common.Redact(runOptions.LogRedactions, string(stdoutAndStderr))
		return err

	case CaptureOutputModeStdout:
		var stdOut, stdErr bytes.Buffer
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr

		err := cmd.Run()

		runResult.Output = common.Redact(runOptions.LogRedactions, stdOut.String())
		runResult.Stderr = common.Redact(runOptions.LogRedactions, stdErr.String())

		return err
	}

	return fmt.Errorf("Invalid output capture mode: %d", runOptions.CaptureOutputMode)
}
