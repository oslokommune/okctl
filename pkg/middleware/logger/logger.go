// Package logger implements a logging middleware
package logger

import (
	"context"
	"github.com/oslokommune/okctl/pkg/truncate"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"
)

// AnonymizeRequestLogger can be implemented by the request types that want to conceal
// items from the logs
type AnonymizeRequestLogger interface {
	AnonymizeRequest(request interface{}) interface{}
}

// AnonymizeResponseLogger can be implemented by the response types that want to conceal
// things from the logs
type AnonymizeResponseLogger interface {
	AnonymizeResponse(response interface{}) interface{}
}

// Logging returns a logging middleware
func Logging(logger *logrus.Logger, endpointTag string, serviceTags ...string) endpoint.Middleware {
	l := &logging{
		log: logger.WithFields(logrus.Fields{
			"service":  strings.Join(serviceTags, "/"),
			"endpoint": endpointTag,
		}),
	}

	return l.ProcessRequest
}

type logging struct {
	log *logrus.Entry
}

// ProcessRequest handles logging of the request
func (l *logging) ProcessRequest(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		l.log.Debug("request received")

		var d string

		switch r := request.(type) {
		case AnonymizeRequestLogger:
			d = litter.Sdump(r.AnonymizeRequest(request))
		default:
			d = litter.Sdump(request)
		}

		l.log.Trace("request: ", d)

		defer func() {
			l.ProcessResponse(err, response, time.Now())
		}()

		return next(ctx, request)
	}
}

// ProcessResponse handles logging of the response
func (l *logging) ProcessResponse(err error, response interface{}, begin time.Time) {
	if err != nil {
		l.log.Errorf("processing request: %s", err.Error())
	}

	if err == nil {
		var d string

		switch r := response.(type) {
		case AnonymizeResponseLogger:
			d = litter.Sdump(r.AnonymizeResponse(response))
		default:
			d = litter.Sdump(response)
		}

		truncatedDump := truncate.String(&d, 5000)
		l.log.Trace("response: ", truncatedDump)
	}

	l.log.Debug("request completed in: ", time.Since(begin).String())
}
