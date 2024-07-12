package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	var (
		stdout, stderr *os.File
	)

	stdout, err := os.CreateTemp(t.TempDir(), "stdout")
	require.NoError(t, err)

	stderr, err = os.CreateTemp(t.TempDir(), "stderr")
	require.NoError(t, err)

	pStdout := os.Stdout
	pStderr := os.Stderr

	defer func() {
		os.Stdout = pStdout
		os.Stderr = pStderr
	}()

	os.Stdout = stdout
	os.Stderr = stderr

	require.NoError(t, run([]string{
		"-file", "testdata/test.jsonnet",
		"-repo", "https://github.com/mumoshu/gitimpart_test.git",
		"-branch", "main",
		"-dry-run",
	}))

	assert.Equal(t, `diff --git a/a.txt b/a.txt
new file mode 100644
index 0000000000000000000000000000000000000000..78981922613b2afb6025042ff6bd878ac1994e85
--- /dev/null
+++ b/a.txt
@@ -0,0 +1 @@
+a
diff --git a/b.json b/b.json
new file mode 100644
index 0000000000000000000000000000000000000000..666608b8979d7b3932723a4fdbc9b176ecddeb13
--- /dev/null
+++ b/b.json
@@ -0,0 +1 @@
+{"b":"B"}
\ No newline at end of file
diff --git a/b.yaml b/b.yaml
new file mode 100644
index 0000000000000000000000000000000000000000..8b2007f3a16f9b5b8575bfc5ee33dab104523df0
--- /dev/null
+++ b/b.yaml
@@ -0,0 +1 @@
+b: B
diff --git a/c.json b/c.json
new file mode 100644
index 0000000000000000000000000000000000000000..2cff59a2f2d020e2e757b40483470c6df22295ae
--- /dev/null
+++ b/c.json
@@ -0,0 +1 @@
+{"c":"C"}
diff --git a/c.yaml b/c.yaml
new file mode 100644
index 0000000000000000000000000000000000000000..0e7435b72e46e3b1a6f7d8b5866c9a8832fd81b7
--- /dev/null
+++ b/c.yaml
@@ -0,0 +1 @@
+c: C
diff --git a/d/e/f.txt b/d/e/f.txt
new file mode 100644
index 0000000000000000000000000000000000000000..71bc41cb110a98401c3f27d0341e69459f434596
--- /dev/null
+++ b/d/e/f.txt
@@ -0,0 +1 @@
+d/e/f
\ No newline at end of file

`, readAll(t, stdout))
	assert.Empty(t, readAll(t, stderr))
}

func readAll(t *testing.T, f *os.File) string {
	t.Helper()

	b, err := os.ReadFile(f.Name())
	require.NoError(t, err)

	return string(b)
}
