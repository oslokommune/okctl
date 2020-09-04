// Package logger implements a logging middleware
package logger

import (
	"context"
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
func Logging(logger *logrus.Logger, serviceTag, endpointTag string) endpoint.Middleware {
	logCtx := logger.WithFields(logrus.Fields{
		"service":  serviceTag,
		"endpoint": endpointTag,
	})

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logCtx.Info("handling request")

			if anonReq, ok := request.(AnonymizeRequestLogger); ok {
				anon := anonReq.AnonymizeRequest(request)
				logCtx.Debug(litter.Sdump(anon))
			} else {
				logCtx.Debug(litter.Sdump(request))
			}

			defer func(begin time.Time) {
				if err != nil {
					logCtx.Warn("failed to process request, because: ", err.Error())
					logCtx.Debug(litter.Sdump(err))
				}

				if err == nil {
					logCtx.Info("done with request, sending response")

					if anonResp, ok := request.(AnonymizeResponseLogger); ok {
						anon := anonResp.AnonymizeResponse(response)
						logCtx.Debug(litter.Sdump(anon))
					} else {
						logCtx.Debug(litter.Sdump(response))
					}

				}

				logCtx.Info("request completed in: ", time.Since(begin).String())
			}(time.Now())

			return next(ctx, request)
		}
	}
}
