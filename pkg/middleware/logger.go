package middleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

func Logging(logger *logrus.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			logger.Info("request", request)
			logger.Info("response", response)

			return next(ctx, request)
		}
	}
}
