package cli

import "testing"

func TestParsePositiveID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in    string
		want  string
		bad   bool
	}{
		{in: "1", want: "1"},
		{in: "42", want: "42"},
		{in: "18446744073709551615", want: "18446744073709551615"},
		{in: "0", bad: true},
		{in: "-1", bad: true},
		{in: "abc", bad: true},
		{in: "1.5", bad: true},
		{in: "../tags/7", bad: true},
		{in: "1/../../users", bad: true},
		{in: "..%2f..%2fusers", bad: true},
		{in: "1?x=..", bad: true},
		{in: " 7", bad: true},
		{in: "7 ", bad: true},
		{in: "", bad: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, err := parsePositiveID(tt.in)
			if tt.bad {
				if err == nil {
					t.Fatalf("expected error for %q", tt.in)
				}
				if ee, ok := err.(*exitError); !ok || ee.code != exitUsage {
					t.Fatalf("expected usage exit for %q: %v", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
