package stack_test

import (
	"path/filepath"
	"testing"

	"github.com/neermitt/opsos/pkg/stack"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackProcessorNoDependency(t *testing.T) {
	stacksPath, err := filepath.Abs("../../examples/complete/stacks")
	require.NoError(t, err)
	afs := afero.NewBasePathFs(afero.NewOsFs(), stacksPath)

	s, err := stack.ProcessYAMLFile(afs, "orgs/cp/_defaults.yaml")

	require.NoError(t, err)

	assert.NotNil(t, s)
}

func TestStackProcessorWithoutFileExt(t *testing.T) {
	stacksPath, err := filepath.Abs("../../examples/complete/stacks")
	require.NoError(t, err)
	afs := afero.NewBasePathFs(afero.NewOsFs(), stacksPath)

	s, err := stack.ProcessYAMLFile(afs, "orgs/cp/_defaults")

	require.NoError(t, err)

	assert.NotNil(t, s)
}

func TestStackProcessorSingleDependency(t *testing.T) {
	stacksPath, err := filepath.Abs("../../examples/complete/stacks")
	require.NoError(t, err)
	afs := afero.NewBasePathFs(afero.NewOsFs(), stacksPath)

	s, err := stack.ProcessYAMLFile(afs, "orgs/cp/tenant1/_defaults.yaml")

	require.NoError(t, err)

	assert.NotNil(t, s)
}
