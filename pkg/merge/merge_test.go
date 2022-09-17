package merge_test

import (
	"testing"

	"github.com/neermitt/opsos/pkg/merge"
	"github.com/stretchr/testify/assert"
)

func TestMergeBasic(t *testing.T) {
	map1 := map[string]any{"foo": "bar"}
	map2 := map[string]any{"baz": "bat"}

	inputs := []map[string]any{map1, map2}
	expected := map[string]any{"foo": "bar", "baz": "bat"}

	result, err := merge.Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMergeBasicOverride(t *testing.T) {
	map1 := map[string]any{"foo": "bar"}
	map2 := map[string]any{"baz": "bat"}
	map3 := map[string]any{"foo": "ood"}

	inputs := []map[string]any{map1, map2, map3}
	expected := map[string]any{"foo": "ood", "baz": "bat"}

	result, err := merge.Merge(inputs)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
