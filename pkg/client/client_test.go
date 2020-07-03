package client_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestClient(t *testing.T) {
	defaultClusterCreateOpts := mock.DefaultClusterCreateOpts()
	defaultClusterDeleteOpts := mock.DefaultClusterDeleteOpts()
	defaultVpcCreateOpts := mock.DefaultVpcCreateOpts()
	defaultVpcDeleteOpts := mock.DefaultVpcDeleteOpts()
	defaultClusterConfigCreateOpts := mock.DefaultCreateClusterConfigOpts()

	defaultResponse := []byte("hi there")
	defaultExpect := "hi there"

	testCases := []struct {
		name      string
		path      string
		method    string
		expect    interface{}
		response  []byte
		expectErr bool
		fn        func(c *client.Client) error
	}{
		{
			name:     "Create cluster works",
			path:     "/clusters/",
			method:   http.MethodPost,
			response: defaultResponse,
			expect:   defaultExpect,
			fn: func(c *client.Client) error {
				return c.CreateCluster(&defaultClusterCreateOpts)
			},
		},
		{
			name:     "Delete cluster works",
			path:     "/clusters/",
			method:   http.MethodDelete,
			response: defaultResponse,
			expect:   defaultExpect,
			fn: func(c *client.Client) error {
				return c.DeleteCluster(&defaultClusterDeleteOpts)
			},
		},
		{
			name:     "Create vpc works",
			path:     "/vpcs/",
			method:   http.MethodPost,
			expect:   defaultExpect,
			response: defaultResponse,
			fn: func(c *client.Client) error {
				return c.CreateVpc(&defaultVpcCreateOpts)
			},
		},
		{
			name:     "Delete vpc works",
			path:     "/vpcs/",
			method:   http.MethodDelete,
			expect:   defaultExpect,
			response: defaultResponse,
			fn: func(c *client.Client) error {
				return c.DeleteVpc(&defaultVpcDeleteOpts)
			},
		},
		{
			name:      "Create cluster config works",
			path:      "/clusterconfigs/",
			method:    http.MethodPost,
			expect:    defaultExpect,
			response:  defaultResponse,
			expectErr: false,
			fn: func(c *client.Client) error {
				return c.CreateClusterConfig(&defaultClusterConfigCreateOpts)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, request.Method, tc.method)
				assert.Equal(t, request.URL.Path, tc.path)

				_, err := writer.Write(tc.response)
				assert.NoError(t, err)
			}))

			got := new(bytes.Buffer)

			c := client.New(got, fmt.Sprintf("%s/", server.URL))

			err := tc.fn(c)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got.String())
			}

			defer server.Close()
		})
	}
}
