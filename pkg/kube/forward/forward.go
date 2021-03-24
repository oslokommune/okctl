// Package forward is based on, with some modifications here and there:
// - https://github.com/gianarb/kube-port-forward
//
// Copyright 2020 Gianluca Arbezzano
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
package forward

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// Forward contains the state required for forwarding
// traffic to a local port
type Forward struct {
	config  *restclient.Config
	stopCh  chan struct{}
	readyCh chan struct{}
	wg      sync.WaitGroup
	stream  genericclioptions.IOStreams
	out     io.Writer
}

// New returns an initialised forwarding client
func New(in io.Reader, out, err io.Writer, config *restclient.Config) *Forward {
	return &Forward{
		config:  config,
		stopCh:  make(chan struct{}, 1),
		readyCh: make(chan struct{}),
		stream: genericclioptions.IOStreams{
			In:     in,
			Out:    out,
			ErrOut: err,
		},
		out: err,
	}
}

// Start forwarding traffic to the provided port
func (f *Forward) Start(listenPort int32, p *v1.Pod) (err error) {
	f.wg.Add(1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		close(f.stopCh)
		f.wg.Done()
	}()

	go func() {
		err = f.portForwardToPod(listenPort, p)
	}()

	<-f.readyCh
	f.wg.Wait()

	return err
}

func (f *Forward) portForwardToPod(listenPort int32, pod *v1.Pod) error {
	path := fmt.Sprintf(
		"/api/v1/namespaces/%s/pods/%s/portforward",
		pod.Namespace,
		pod.Name,
	)

	hostIP := strings.TrimLeft(f.config.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(f.config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(
		upgrader,
		&http.Client{Transport: transport},
		http.MethodPost,
		&url.URL{Scheme: "https", Path: path, Host: hostIP},
	)

	fw, err := portforward.New(
		dialer,
		[]string{fmt.Sprintf("%d:%d", listenPort, listenPort)},
		f.stopCh,
		f.readyCh,
		f.stream.Out,
		f.stream.ErrOut,
	)
	if err != nil {
		return fmt.Errorf("starting port forwarding: %w", err)
	}

	return fw.ForwardPorts()
}
