package cli

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestIDArgsRejectPathTraversal(t *testing.T) {
	t.Parallel()
	badIDs := []string{
		"../tags/7",
		"1/../../users",
		"..%2f..%2fusers",
		"1?x=..",
		"0",
		"-3",
	}
	commands := []struct {
		name string
		args []string
	}{
		{"bookmarks get", []string{"bookmarks", "get"}},
		{"bookmarks update", []string{"bookmarks", "update", "--json", `{"url":"https://x"}`}},
		{"bookmarks archive", []string{"bookmarks", "archive"}},
		{"bookmarks unarchive", []string{"bookmarks", "unarchive"}},
		{"bookmarks delete", []string{"bookmarks", "delete"}},
		{"tags get", []string{"tags", "get"}},
	}
	for _, cmd := range commands {
		cmd := cmd
		for _, id := range badIDs {
			id := id
			t.Run(cmd.name+"/"+id, func(t *testing.T) {
				t.Parallel()
				var hits int32
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					atomic.AddInt32(&hits, 1)
					w.WriteHeader(http.StatusOK)
				}))
				defer srv.Close()
				args := append(append([]string(nil), cmd.args...), id)
				_, errOut, code := runCLI(t, srv.URL, args...)
				if code != exitUsage {
					t.Fatalf("code=%d want %d stderr=%s", code, exitUsage, errOut)
				}
				if atomic.LoadInt32(&hits) != 0 {
					t.Fatalf("server was contacted for malicious id %q", id)
				}
			})
		}
	}
}

func TestIDArgsValidRequestPath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"bookmarks", "get", "7"}, "/api/bookmarks/7/"},
		{[]string{"tags", "get", "3"}, "/api/tags/3/"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			var gotPath string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":7}`))
			}))
			defer srv.Close()
			_, errOut, code := runCLI(t, srv.URL, tc.args...)
			if code != 0 {
				t.Fatalf("code=%d stderr=%s", code, errOut)
			}
			if gotPath != tc.want {
				t.Fatalf("path %q want %q", gotPath, tc.want)
			}
		})
	}
}
