package terraform

import (
	"context"
	"path"
	"testing"

	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/stretchr/testify/require"
)

func TestProvider_ProcessComponentConfig(t *testing.T) {
	stackConfig := schema.StackConfig{
		KubeConfigProvider: "",
		Components: schema.ComponentsConfig{
			Types: map[string]map[string]schema.ComponentConfig{
				ComponentType: {
					"NoComponentOverride": {},
					"ComponentOverrides": {
						Vars: map[string]any{
							"key2": "val-override-2",
							"key3": "val-3",
						},
						Envs: map[string]string{
							"env1": "val-override-1",
							"env3": "val3",
						},
						BackendType: "remote",
					},
					"OverrideComponent1": {
						Component: "ComponentOverrides",
						Vars: map[string]any{
							"key3": "val-override-3",
							"key4": "val-4",
						},
						Envs: map[string]string{
							"env4": "val-4",
						},
						BackendType: "local",
					},
					"OverrideComponent2": {
						Component: "OverrideComponent1",
						Vars: map[string]any{
							"key4": "val-override-4",
						},
						Envs: map[string]string{
							"env4": "val-override-4",
						},
						BackendType: "s3",
						Backend: map[string]any{
							"path": "terraform1.tfstate",
						},
					},
					"infra/overrides/OverrideComponent3": {
						Component: "OverrideComponent1",
						Vars: map[string]any{
							"key4": "val-override-4",
						},
						Envs: map[string]string{
							"env4": "val-override-4",
						},
						BackendType: "s3",
						Backend: map[string]any{
							"path": "terraform1.tfstate",
						},
					},
				},
			},
		},
		ComponentTypeSettings: map[string]schema.ComponentTypeSettings{
			ComponentType: {
				Vars:        nil,
				Envs:        nil,
				BackendType: "local",
			},
		},
	}
	testContext := context.WithValue(
		context.WithValue(context.Background(), stackConfigKey, &stackConfig),
		terraformStackConfigKey, map[string]any{
			varsField: map[string]any{
				"key1": "val1",
				"key2": "val2",
			},
			envsField: map[string]string{
				"env1": "val1",
				"env2": "val2",
			},
			backendTypeField: "local",
			backendField: map[string]any{
				"path": "terraform.tfstate",
			},
		})

	componentPathFunc := func(componentName string) (string, error) {
		return path.Join("a/b/c/", componentName), nil
	}

	tests := []struct {
		componentName         string
		expectedError         bool
		expectedComponentInfo map[string]any
	}{
		{
			componentName: "MissingComponent",
			expectedError: true,
		},
		{
			componentName: "NoComponentOverride",
			expectedComponentInfo: map[string]any{
				varsField: map[string]any{
					"key1": "val1",
					"key2": "val2",
				},
				envsField: map[string]any{
					"env1": "val1",
					"env2": "val2",
				},
				backendTypeField: "local",
				backendField: map[string]any{
					"path": "terraform.tfstate",
				},
			},
		},
		{
			componentName: "ComponentOverrides",
			expectedComponentInfo: map[string]any{
				varsField: map[string]any{
					"key1": "val1",
					"key2": "val-override-2",
					"key3": "val-3",
				},
				envsField: map[string]any{
					"env1": "val-override-1",
					"env2": "val2",
					"env3": "val3",
				},
				backendTypeField: "remote",
				backendField: map[string]any{
					"path": "terraform.tfstate",
				},
			},
		},
		{
			componentName: "OverrideComponent1",
			expectedComponentInfo: map[string]any{
				varsField: map[string]any{
					"key1": "val1",
					"key2": "val-override-2",
					"key3": "val-override-3",
					"key4": "val-4",
				},
				envsField: map[string]any{
					"env1": "val-override-1",
					"env2": "val2",
					"env3": "val3",
					"env4": "val-4",
				},
				backendTypeField: "local",
				backendField: map[string]any{
					"path": "terraform.tfstate",
				},
			},
		},
		{
			componentName: "OverrideComponent2",
			expectedComponentInfo: map[string]any{
				varsField: map[string]any{
					"key1": "val1",
					"key2": "val-override-2",
					"key3": "val-override-3",
					"key4": "val-override-4",
				},
				envsField: map[string]any{
					"env1": "val-override-1",
					"env2": "val2",
					"env3": "val3",
					"env4": "val-override-4",
				},
				backendTypeField: "s3",
				backendField: map[string]any{
					"path": "terraform1.tfstate",
				},
			},
		},
		{
			componentName: "infra/overrides/OverrideComponent3",
			expectedComponentInfo: map[string]any{
				varsField: map[string]any{
					"key1": "val1",
					"key2": "val-override-2",
					"key3": "val-override-3",
					"key4": "val-override-4",
				},
				envsField: map[string]any{
					"env1": "val-override-1",
					"env2": "val2",
					"env3": "val3",
					"env4": "val-override-4",
				},
				backendTypeField: "s3",
				backendField: map[string]any{
					"path": "terraform1.tfstate",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.componentName, func(t *testing.T) {
			testCase := tc
			t.Parallel()
			componentInfo, err := processComponent(testContext, testCase.componentName, componentPathFunc)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.expectedComponentInfo, componentInfo)
		})
	}
}
