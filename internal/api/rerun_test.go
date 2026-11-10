package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReRunRun(t *testing.T) {
	for _, tc := range []struct {
		name       string
		failedOnly bool
		action     string
		status     int
		wantError  bool
	}{
		{"all jobs", false, "rerun", http.StatusCreated, false},
		{"failed jobs", true, "rerun-failed-jobs", http.StatusCreated, false},
		{"rejected", true, "rerun-failed-jobs", http.StatusForbidden, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}
				if want := "/repos/some/repo/actions/runs/123/" + tc.action; r.URL.Path != want {
					t.Errorf("path = %s, want %s", r.URL.Path, want)
				}
				w.WriteHeader(tc.status)
			}))
			defer server.Close()
			client := API{url: server.URL, httpClient: server.Client()}
			err := client.ReRunRun("some/repo", "123", tc.failedOnly)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
		})
	}
}
