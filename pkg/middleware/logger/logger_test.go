package logger_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/go-kit/kit/endpoint"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/oslokommune/okctl/pkg/middleware/logger"
)

type Fake struct {
	Name string
}

func Endpoint(response interface{}, err error) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return response, err
	}
}

func TestLogging(t *testing.T) {
	testCases := []struct {
		name     string
		request  interface{}
		response interface{}
		err      error
		expect   []*regexp.Regexp
	}{
		{
			name:     "Should work",
			request:  Fake{Name: "hi"},
			response: Fake{Name: "goodbye"},
			expect: []*regexp.Regexp{
				regexp.MustCompile("time=\"(.*?)\" level=info msg=\"request received\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request: logger_test.Fake(.*?)hi(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"response: logger_test.Fake(.*?)goodbye(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=info msg=\"request completed in: (.*?)\" endpoint=create service=something/else"),
			},
		},
		{
			name:     "With error",
			request:  Fake{Name: "hi"},
			response: Fake{Name: "goodbye"},
			err:      fmt.Errorf("oh no"),
			expect: []*regexp.Regexp{
				regexp.MustCompile("time=\"(.*?)\" level=info msg=\"request received\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request: logger_test.Fake(.*?)hi(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=error msg=\"processing request: oh no\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=info msg=\"request completed in: (.*?)\" endpoint=create service=something/else"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			l, hook := test.NewNullLogger()
			l.SetLevel(logrus.DebugLevel)
			_, _ = logger.Logging(l, "create", "something", "else")(Endpoint(tc.response, tc.err))(context.Background(), tc.request)

			for i, entry := range hook.AllEntries() {
				msg, err := entry.String()
				assert.NoError(t, err)
				assert.Regexp(t, tc.expect[i], msg)
			}
		})
	}
}
