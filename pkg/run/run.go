package run

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Runner interface {
	Run(progress io.Writer, args []string) ([]byte, error)
}

type Run struct {
	WorkingDirectory string
	BinaryPath       string
	Env              []string
	Logger           *logrus.Logger
}

func (r *Run) Run(progress io.Writer, args []string) ([]byte, error) {
	var errOut, errErr error

	ctxLogger := logrus.WithFields(
		logrus.Fields{
			"component": "generic_runner",
			"binary":    r.BinaryPath,
			"args":      args,
			"env":       r.Env,
			"dir":       r.WorkingDirectory,
		},
	)

	cmd := &exec.Cmd{
		Path: r.BinaryPath,
		Args: append([]string{r.BinaryPath}, args...),
		Env:  r.Env,
		Dir:  r.WorkingDirectory,
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

	ctxLogger.Info("Starting execution of provided command")

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

func New(logger *logrus.Logger, workingDirectory, binaryPath string, env []string) *Run {
	return &Run{
		WorkingDirectory: workingDirectory,
		BinaryPath:       binaryPath,
		Env:              env,
		Logger:           logger,
	}
}
