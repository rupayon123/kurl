package main

import "testing"

func TestJoinURLKeepsQueryAndFragmentOutOfPath(t *testing.T) {
	for _, tc := range []struct{ base, path, want string }{
		{"https://api.test/v1?tenant=one", "users", "https://api.test/v1/users?tenant=one"},
		{"https://api.test/v1?tenant=one", "users?limit=2", "https://api.test/v1/users?limit=2"},
		{"https://api.test/v1#old", "users#new", "https://api.test/v1/users#new"},
		{"https://api.test/v1", "?limit=2", "https://api.test/v1?limit=2"},
		{"https://api.test/v1?old=1", "users?", "https://api.test/v1/users?"},
		{"https://api.test/v1", "a%2Fb", "https://api.test/v1/a%2Fb"},
	} {
		if got := joinURL(tc.base, tc.path); got != tc.want {
			t.Errorf("joinURL(%q,%q)=%q want %q", tc.base, tc.path, got, tc.want)
		}
	}
}
