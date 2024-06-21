package gitimpart_test

import (
	"os"
	"testing"

	"github.com/mumoshu/gitimpart"
	"github.com/stretchr/testify/require"
)

func TestGitimpartRender(t *testing.T) {
	r, err := gitimpart.RenderFile("testdata/test.jsonnet")
	require.NoError(t, err)

	require.Equal(t, gitimpart.Contents{
		Files: map[string]interface{}{
			"a.txt": "aaa",
			"b.json": map[string]interface{}{
				"foo": "bar",
			},
		},
	}, *r)
}

func TestGitimpartPush(t *testing.T) {
	ghtoken := os.Getenv("GITHUB_TOKEN")
	if ghtoken == "" {
		t.Skip("GITHUB_TOKEN is not set")
	}

	r, err := gitimpart.RenderFile("testdata/test.jsonnet")
	require.NoError(t, err)

	err = gitimpart.Push(
		*r,
		"https://github.com/mumoshu/gitimpart_test.git",
		"main",
		gitimpart.WithGitHubToken(ghtoken),
	)
	require.NoError(t, err)
}

func TestGitimpartPush_PullRequest(t *testing.T) {
	ghtoken := os.Getenv("GITHUB_TOKEN")
	if ghtoken == "" {
		t.Skip("GITHUB_TOKEN is not set")
	}

	r, err := gitimpart.RenderFile("testdata/test.jsonnet")
	require.NoError(t, err)

	err = gitimpart.Push(
		*r,
		"https://github.com/mumoshu/gitimpart_test.git",
		"main",
		gitimpart.WithGitHubToken(ghtoken),
		gitimpart.WithPullRequest(),
	)
	require.NoError(t, err)
}
