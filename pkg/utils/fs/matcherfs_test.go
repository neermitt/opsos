package fs_test

import (
	"testing"

	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatcherFs(t *testing.T) {
	afs := afero.NewMemMapFs()
	require.NoError(t, afs.MkdirAll("src/a/b/c", 0755))
	require.NoError(t, afs.MkdirAll("src/d/e/f", 0755))
	require.NoError(t, afs.MkdirAll("org/a/b/c", 0755))
	require.NoError(t, afs.MkdirAll("org/d/e/f", 0755))
	require.NoError(t, afero.WriteFile(afs, "src/a/b/c/test1", []byte("file test1"), 0644))
	require.NoError(t, afero.WriteFile(afs, "src/d/e/f/test2", []byte("file test2"), 0644))
	require.NoError(t, afero.WriteFile(afs, "org/a/b/c/test3", []byte("file test3"), 0644))
	require.NoError(t, afero.WriteFile(afs, "org/d/e/f/test4", []byte("file test4"), 0644))
	require.NoError(t, afero.WriteFile(afs, "org/d/e/f/_defaults.yaml", []byte("file defaults"), 0644))

	gfs := fs.NewMatcherFs(afs, fs.And(fs.Glob("org/**/*"), fs.Not(fs.Glob("**/_defaults.yaml"))))
	assert.NotNil(t, gfs)

	match, err := fs.AllFiles(gfs)
	require.NoError(t, err)

	assert.ElementsMatch(t, match, []string{"org/a/b/c/test3", "org/d/e/f/test4"})
}
