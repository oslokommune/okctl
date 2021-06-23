// Package run knows how to execute binaries and collect their output and result
package run

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// Runner defines the interface required by a concrete
// runner implementation
type Runner interface {
	Run(progress io.Writer, args []string) ([]byte, error)
}

// Run stores state related to running a given command
type Run struct {
	WorkingDirectory string
	BinaryPath       string
	Env              []string
	CmdFn            CmdFn
	Logger           *logrus.Logger
}

// AnonymizeEnv will conceal any secret portions of the
// provided environment variables, given that it knows
// of them
func AnonymizeEnv(entries []string) []string {
	out := make([]string, len(entries))

	hide := []string{
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	}

	for i, entry := range entries {
		for _, h := range hide {
			if strings.Contains(entry, h) {
				entry = fmt.Sprintf("%s=XXXXXXX", h)
				break
			}
		}

		out[i] = entry
	}

	return out
}

// CmdFn defines a type for creating a compatible *exec.Cmd
type CmdFn func(workingDir, path string, env, args []string) *exec.Cmd

// Cmd returns a standard *exec.Cmd creator.
func Cmd() CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		return &exec.Cmd{
			Path: path,
			Args: append([]string{path}, args...),
			Env:  env,
			Dir:  workingDir,
		}
	}
}

// Run a command and record the progress to the provided writer
func (r *Run) Run(progress io.Writer, args []string) ([]byte, error) {
	cmd := r.CmdFn(r.WorkingDirectory, r.BinaryPath, r.Env, args)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	multiReader := io.MultiReader(stdoutIn, stderrIn)
	multiWriter := io.MultiWriter(progress, &buff)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(multiWriter, multiReader)

	if r.Logger != nil {
		scanner := bufio.NewScanner(multiReader)

		for scanner.Scan() {
			s := scanner.Text()
			r.Logger.Info(s)
		}
	}

	err = cmd.Wait()
	if err != nil {
		return buff.Bytes(), fmt.Errorf("executing command: %s, got: %w", buff.Bytes(), err)
	}

	return buff.Bytes(), nil
}

// New returns a runner capable of executing
// commands using the provided binary
func New(logger *logrus.Logger, workingDirectory, binaryPath string, env []string, cmdFn CmdFn) *Run {
	return &Run{
		WorkingDirectory: workingDirectory,
		BinaryPath:       binaryPath,
		Env:              env,
		CmdFn:            cmdFn,
		Logger:           logger,
	}
}
