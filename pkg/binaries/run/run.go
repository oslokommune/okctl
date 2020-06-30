// Package run knows how to execute binaries and collect their output and result
package run

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
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

// Cmd returns a standard *exec.Cmd creator
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
	var errOut, errErr error

	cmd := r.CmdFn(r.WorkingDirectory, r.BinaryPath, r.Env, args)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	var errBuff, outBuff bytes.Buffer
	stdout := io.MultiWriter(progress, &outBuff)
	stderr := io.MultiWriter(progress, &errBuff)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		_, errOut = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errErr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		return errBuff.Bytes(), err
	}

	if errOut != nil || errErr != nil {
		return errBuff.Bytes(), err
	}

	return outBuff.Bytes(), nil
}

// New returns a runner capable of executing
// commands using the provided binary
func New(workingDirectory, binaryPath string, env []string) *Run {
	return &Run{
		WorkingDirectory: workingDirectory,
		BinaryPath:       binaryPath,
		Env:              env,
		CmdFn:            Cmd(),
	}
}
