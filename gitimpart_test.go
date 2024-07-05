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
			"a.txt": "a\n",
			"b.json": map[string]interface{}{
				"b": "B",
			},
			"b.yaml": map[string]interface{}{
				"b": "B",
			},
			"c.json":    "{\"c\":\"C\"}\n",
			"c.yaml":    "c: C\n",
			"d/e/f.txt": "d/e/f",
		},
	}, *r)
}

func TestGitimpartRender_Kustomize(t *testing.T) {
	r, err := gitimpart.RenderFile("testdata/test.kustomize.jsonnet", gitimpart.Vars(map[string]string{
		"project": "myproject",
	}))
	require.NoError(t, err)

	require.Equal(t, gitimpart.Contents{
		Files: map[string]interface{}{
			"a.txt": "a\n",
			"path/to/kustomization.yaml/dir/projects/myproject.yaml": `metadata:
  name: "myproject"
`,
		},
		Kustomize: map[string]map[string]interface{}{
			"path/to/kustomization.yaml/dir": {
				"projects/myproject.yaml": nil,
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

func TestGitimpartPush_Kustomize(t *testing.T) {
	ghtoken := os.Getenv("GITHUB_TOKEN")
	if ghtoken == "" {
		t.Skip("GITHUB_TOKEN is not set")
	}

	r, err := gitimpart.RenderFile("testdata/test.kustomize.jsonnet")
	require.NoError(t, err)

	err = gitimpart.Push(
		*r,
		"https://github.com/mumoshu/gitimpart_test.git",
		"main",
		gitimpart.WithGitHubToken(ghtoken),
		gitimpart.WithKustomizeBin("./kustomize"),
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
