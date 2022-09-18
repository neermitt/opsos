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
	proc := stack.NewStackProcessor(fs)
	s, err := proc.GetStack("orgs/cp/_defaults.yaml")
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorWithoutFileExt(t *testing.T) {
	proc := stack.NewStackProcessor(fs)
	s, err := proc.GetStack("orgs/cp/_defaults")
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorSingleDependency(t *testing.T) {
	proc := stack.NewStackProcessor(fs)
	s, err := proc.GetStack("orgs/cp/tenant1/_defaults")
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestStackProcessorMultipleFiles(t *testing.T) {
	proc := stack.NewStackProcessor(fs)
	s, err := proc.GetStacks([]string{
		"orgs/cp/tenant1/dev/us-east-2",
		"orgs/cp/tenant1/dev/us-east-2-extras",
		"orgs/cp/tenant1/prod/us-east-2",
		"orgs/cp/tenant1/test1/us-east-2",
	})
	require.NoError(t, err)
	assert.Len(t, s, 4)
}
