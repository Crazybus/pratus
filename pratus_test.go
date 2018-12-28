package main

import "testing"

func TestParseGitHubURL(t *testing.T) {
	owner, repo, number := parseGitHubURL("https://github.com/", "https://github.com/Crazybus/pratus/pulls/123")
	if owner != "Crazybus" {
		t.Errorf("got: %s, want: %s", owner, "Crazybus")
	}
	if repo != "pratus" {
		t.Errorf("got: %s, want: %s", repo, "pratus")
	}
	if number != 123 {
		t.Errorf("got: %d, want: %d", number, 123)
	}
}
