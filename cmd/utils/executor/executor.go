// Copyright 2021 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package executor

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui/screenbuf"
)

type ManifestExecutor interface {
	Execute(command string, env []string) error
	CleanUp()
}

type Executor struct {
	outputMode string
	displayer  executorDisplayer
}

type commandInfo struct {
	command string
	sb      *screenbuf.ScreenBuf
}

type executorDisplayer interface {
	display(scanner *bufio.Scanner)
	startCommand(cmd *exec.Cmd) (io.Reader, error)
	addCommandInfo(cmdInfo *commandInfo)
	cleanUp()
}

// NewExecutor returns a new executor
func NewExecutor(output string) *Executor {

	var displayer executorDisplayer
	switch output {
	case "tty":
		displayer = newTTYExecutorDisplayer()
	case "plain":
		displayer = newPlainExecutorDisplayer()
	case "json":
		displayer = newJsonExecutorDisplayer()
	default:
		displayer = newTTYExecutorDisplayer()
	}
	return &Executor{
		outputMode: output,
		displayer:  displayer,
	}
}

// Execute executes the specified command adding `env` to the execution environment
func (e *Executor) Execute(command string, env []string) error {

	cmd := exec.Command("bash", "-c", command)
	cmd.Env = append(os.Environ(), env...)

	reader, err := e.displayer.startCommand(cmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)

	sb := screenbuf.New(os.Stdout)
	e.displayer.addCommandInfo(&commandInfo{
		command: command,
		sb:      sb,
	})
	go e.displayer.display(scanner)

	err = cmd.Wait()
	if e.outputMode == "tty" {
		collapseTTY(command, err, sb)
	}
	return err
}

func startCommand(cmd *exec.Cmd) (io.Reader, error) {
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return reader, nil
}

// Execute executes the specified command adding `env` to the execution environment
func (e *Executor) CleanUp() {
	e.displayer.cleanUp()
}
