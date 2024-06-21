package convention_test

import (
	"testing"

	"github.com/mumoshu/gitimpart/convention"
)

func TestRepoURL(t *testing.T) {
	testcases := []struct {
		name string
		r    string
		want string
	}{
		{
			name: "https dot git",
			r:    "https://github.com/mumoshu/example.git",
			want: "https://github.com/mumoshu/example.git",
		},
		{
			name: "https without dot git",
			r:    "https://github.com/mumoshu/example",
			want: "https://github.com/mumoshu/example",
		},
		{
			name: "owner/repo",
			r:    "mumoshu/example",
			want: "https://github.com/mumoshu/example.git",
		},
		{
			name: "github.com/owner/repo",
			r:    "mumoshu/example",
			want: "https://github.com/mumoshu/example.git",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := convention.RepoURL(tc.r)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
