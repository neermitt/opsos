package v1_test

import (
	"testing"

	"github.com/neermitt/opsos/api/common"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigRead(t *testing.T) {
	component, err := v1.ReadAndMergeConfigsFromDirs([]string{"../.."})
	require.NoError(t, err)

	assert.Equal(t, &v1.Config{
		Object: common.Object{
			ApiVersion: "opsos/v1",
			Kind:       "Configuration",
			Metadata: common.ObjectMetadata{
				Name:        "opsos-test-config",
				Description: "OPSOS Test Configuration",
			},
		},
		Spec: v1.ConfigSpec{
			BasePath: stringPtr("examples/complete"),
			Stacks: &v1.StacksSpec{
				BasePath: stringPtr("stacks"),
				IncludedPaths: []string{
					"orgs/**/*",
				},
				ExcludedPaths: []string{
					"**/_defaults.yaml",
				},
				NamePattern: stringPtr("{{.tenant}}-{{.environment}}-{{.stage}}"),
			},
			Providers: map[string]v1.ProviderSettings{
				"helmfile": {
					"base_path":       "components/helmfile",
					"kubeconfig_path": "/dev/shm",
				},
				"terraform": {
					"base_path":                  "components/terraform",
					"apply_auto_approve":         false,
					"deploy_run_init":            true,
					"init_run_reconfigure":       true,
					"auto_generate_backend_file": false,
				},
				"kind": {
					"cluster_name_pattern": "{{.namespace}}-{{.tenant}}-{{.environment}}-{{.stage}}",
				},
			},
		},
	}, component)
}

func stringPtr(s string) *string {
	return &s
}
