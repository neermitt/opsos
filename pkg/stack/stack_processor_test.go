package stack_test

import (
	"path/filepath"
	"testing"

	"github.com/neermitt/opsos/pkg/stack"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fs afero.Fs

func TestMain(m *testing.M) {
	stacksPath, _ := filepath.Abs("../../examples/complete/stacks")
	fs = afero.NewBasePathFs(afero.NewOsFs(), stacksPath)
	m.Run()
}

func TestStackProcessorNoDependency(t *testing.T) {
	proc := stack.NewStackProcessor(fs, []string{"orgs/**/*"}, []string{"**/_defaults.yaml"}, "test")
	s, err := proc.GetStack("orgs/cp/_defaults.yaml", stack.ProcessStackOptions{})
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorWithoutFileExt(t *testing.T) {
	proc := stack.NewStackProcessor(fs, []string{"orgs/**/*"}, []string{"**/_defaults.yaml"}, "test")
	s, err := proc.GetStack("orgs/cp/_defaults", stack.ProcessStackOptions{})
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorSingleDependency(t *testing.T) {
	proc := stack.NewStackProcessor(fs, []string{"orgs/**/*"}, []string{"**/_defaults.yaml"}, "test")
	s, err := proc.GetStack("orgs/cp/tenant1/_defaults", stack.ProcessStackOptions{})
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorMultipleFiles(t *testing.T) {
	proc := stack.NewStackProcessor(fs, []string{"orgs/**/*"}, []string{"**/_defaults.yaml"}, "test")
	names, err := proc.GetStackNames()
	require.NoError(t, err)
	s, err := proc.GetStacks(names)
	require.NoError(t, err)
	assert.Len(t, s, 15)
}

func TestStackProcessorLoadStackWithMixin(t *testing.T) {
	proc := stack.NewStackProcessor(fs, []string{"orgs/**/*"}, []string{"**/_defaults.yaml"}, "test")
	s, err := proc.GetStack("orgs/cp/tenant1/dev/us-east-2", stack.ProcessStackOptions{
		ComponentType: "terraform",
		ComponentName: "test/test-component-override-3",
	})
	require.NoError(t, err)
	assert.Len(t, s.Components, 1)
	assert.Len(t, s.Components["terraform"], 1)
}
