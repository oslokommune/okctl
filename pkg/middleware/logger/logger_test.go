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

// nolint: funlen
func TestLogging(t *testing.T) {
	testCases := []struct {
		name           string
		request        interface{}
		response       interface{}
		expectResponse interface{}
		truncateBytes  int
		err            error
		expect         []*regexp.Regexp
	}{
		{
			name:          "Should work",
			request:       Fake{Name: "hi"},
			response:      Fake{Name: "goodbye"},
			truncateBytes: 5000,
			expect: []*regexp.Regexp{
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request received\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=trace msg=\"request: logger_test.Fake(.*?)hi(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=trace msg=\"response: logger_test.Fake(.*?)goodbye(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request completed in: (.*?)\" endpoint=create service=something/else"),
			},
		},
		{
			name:          "With error",
			request:       Fake{Name: "hi"},
			truncateBytes: 5000,
			response:      Fake{Name: "goodbye"},
			err:           fmt.Errorf("oh no"),
			expect: []*regexp.Regexp{
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request received\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=trace msg=\"request: logger_test.Fake(.*?)hi(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=error msg=\"processing request: oh no\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request completed in: (.*?)\" endpoint=create service=something/else"),
			},
		},
		{
			name:           "Should keep response intact when truncating",
			request:        Fake{Name: "hi"},
			truncateBytes:  30,
			response:       Fake{Name: "goodbye and hello, my dear old friend"},
			expectResponse: Fake{Name: "goodbye and hello, my dear old friend"},
			expect: []*regexp.Regexp{
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request received\" endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=trace msg=\"request: logger_test.Fake(.*?)hi(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=trace msg=\"response: logger_test.Fake(.*?)gooXXXtruncated38bytesXXX(.*?)endpoint=create service=something/else"),
				regexp.MustCompile("time=\"(.*?)\" level=debug msg=\"request completed in: (.*?)\" endpoint=create service=something/else"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			logger.TruncateResponseAtLength = tc.truncateBytes

			l, hook := test.NewNullLogger()
			l.SetLevel(logrus.TraceLevel)
			got, _ := logger.Logging(l, "create", "something", "else")(Endpoint(tc.response, tc.err))(context.Background(), tc.request)

			if tc.expectResponse != nil {
				assert.Equal(t, tc.expectResponse, got)
			}

			for i, entry := range hook.AllEntries() {
				msg, err := entry.String()
				assert.NoError(t, err)
				assert.Regexp(t, tc.expect[i], msg)
			}
		})
	}
}
