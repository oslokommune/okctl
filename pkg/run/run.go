package run

import (
	"bytes"
	"io"
	"os/exec"
)

type Runner interface {
	Run(progress io.Writer, args string) ([]byte, error)
}

type Run struct {
	WorkingDirectory string
	BinaryPath       string
	Env              []string
}

func (e *Run) Run(progress io.Writer, args []string) ([]byte, error) {
	var errOut, errErr error

	cmd := &exec.Cmd{
		Path: e.BinaryPath,
		Args: args,
		Env:  e.Env,
		Dir:  e.WorkingDirectory,
	}

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

func New(workingDirectory, binaryPath string) *Run {
	return &Run{
		WorkingDirectory: workingDirectory,
		BinaryPath:       binaryPath,
	}
}
