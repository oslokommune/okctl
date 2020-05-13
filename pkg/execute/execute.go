package execute

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Executor interface {
	Execute(progress io.Writer, command string) (interface{}, error)
}

type executor struct {
}

func (e *executor) Execute(progress io.Writer, command string) (interface{}, error) {
	var errOut, errErr error

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = wd

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
		return nil, err
	}

	if errOut != nil || errErr != nil {
		return nil, err
	}

	return outBuff.Bytes(), nil
}

func New() Executor {
	return &executor{}
}
