package auth

import (
	"fmt"
	"testing"
)

func TestAuthentication(t *testing.T) {
	var tests = []struct {
		username, pw string
		want         bool
	}{
		{"default_user", "123456", true},
		{"default_user", "abcdef", false},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s, %s", test.username, test.pw)
		t.Run(testname, func(t *testing.T) {
			authenticated := validate(test.username, test.pw)
			if authenticated != test.want {
				t.Errorf("got %t, want %t", authenticated, test.want)
			}
		})
	}
}
