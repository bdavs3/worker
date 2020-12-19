package auth

import (
	"testing"
)

func TestAuthentication(t *testing.T) {
	var tests = []struct {
		comment      string
		username, pw string
		want         bool
	}{
		{
			comment:  "user with correct password",
			username: "default_user",
			pw:       "123456",
			want:     true,
		},
		{
			comment:  "user with incorrect password",
			username: "default_user",
			pw:       "abcdef",
			want:     false,
		},
		{
			comment:  "unknown username",
			username: "four_tet",
			pw:       "16oceans",
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.comment, func(t *testing.T) {
			authenticated := validate(test.username, test.pw)
			if authenticated != test.want {
				t.Errorf("got %t, want %t", authenticated, test.want)
			}
		})
	}
}
