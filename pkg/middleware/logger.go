// Package middleware implements som common middlewares
package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"
)

// Logging returns a logging middleware
func Logging(logger *logrus.Logger, serviceTag, endpointTag string) endpoint.Middleware {
	logCtx := logger.WithFields(logrus.Fields{
		"service":  serviceTag,
		"endpoint": endpointTag,
	})

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logCtx.Info("handling request: ", litter.Sdump(request))

			defer func(begin time.Time) {
				if err != nil {
					logCtx.Warn("failed to process request, because: ", err.Error())
				}

				if err == nil {
					logCtx.Info("done with request, sending response: ", litter.Sdump(request))
				}

				logCtx.Info("request completed in: ", time.Since(begin).String())
			}(time.Now())

			return next(ctx, request)
		}
	}
}
