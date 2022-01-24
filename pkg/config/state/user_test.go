package state_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/stretchr/testify/assert"
)

func TestUserInfoUsernameValidation(t *testing.T) {
	testCases := []struct {
		name        string
		username    string
		expectError bool
		expect      interface{}
	}{
		{
			name:     "Correct username, minimum length",
			username: "ooo1234",
		},
		{
			name:     "Correct username, maximum length",
			username: "ooo1234567",
		},
		{
			name:        "Too short username",
			username:    "ooo123",
			expectError: true,
			expect:      "Username: username must be in the form: yyyXXXXXX (y = letter, x = digit).",
		},
		{
			name:        "Too long username",
			username:    "ooo12345678",
			expectError: true,
			expect:      "Username: username must be in the form: yyyXXXXXX (y = letter, x = digit).",
		},
	}
	uuID := "98e3ddf8-808f-47df-a65d-2dd173f93f8b"

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			userInfo := state.UserInfo{ID: uuID, Username: tc.username}
			err := userInfo.Validate()
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
