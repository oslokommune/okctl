// Package rotatefilehook implements a rotating file logger for logrus
package rotatefilehook

// MIT License
//
// Copyright (c) 2017 Zach
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This code was shamelessly stolen from: https://github.com/snowzach/rotatefilehook

import (
	"io"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateFileConfig contains the configuration
type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Levels     []logrus.Level
	Formatter  logrus.Formatter
}

// RotateFileHook saves the state of the logger
type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

// NewRotateFileHook returns a hook that writes to a log file and rotates it
func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	hook := RotateFileHook{
		Config: config,
	}

	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	return &hook, nil
}

// Levels returns the loglevels
func (hook *RotateFileHook) Levels() []logrus.Level {
	return hook.Config.Levels
}

// Fire writes to the file
func (hook *RotateFileHook) Fire(entry *logrus.Entry) error {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.logWriter.Write(b)
	if err != nil {
		return err
	}

	return nil
}
