package components

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessComponentConfigs(t *testing.T) {

	baseConfig := Config{
		Vars: map[string]any{
			"key1": "val1",
			"key2": "val2",
		},
		Envs: map[string]string{
			"env1": "val1",
			"env2": "val2",
		},
		BackendType: "s3",
		BackendConfigs: map[string]any{
			"s3": map[string]any{
				"encrypt":        true,
				"bucket":         "cp-ue2-root-tfstate",
				"key":            "terraform.tfstate",
				"dynamodb_table": "cp-ue2-root-tfstate-lock",
				"acl":            "bucket-owner-full-control",
				"region":         "us-east-2",
				"role_arn":       nil,
			},
			"azurerm": map[string]any{
				"subscription_id":      "88888-8888-8888-8888-8888888888",
				"resource_group_name":  "rg-terraform-state",
				"storage_account_name": "staterraformstate",
				"container_name":       "dev-tfstate",
				"key":                  "dev.atmos",
			},
			"remote": nil,
			"vault":  nil,
		},
		RemoteStateBackendConfigs: map[string]any{
			"s3": map[string]any{
				"role_arn": "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
			},
		},
	}

	componentsConfigMap := map[string]ConfigWithMetadata{
		"NoComponentOverride": {},
		"ComponentOverrides": {
			Config: Config{
				Vars: map[string]any{
					"key2": "val-override-2",
					"key3": "val-3",
				},
				Envs: map[string]string{
					"env1": "val-override-1",
					"env3": "val3",
				},
				BackendType: "remote",
				BackendConfigs: map[string]any{
					"s3": map[string]any{
						"bucket": "cp-ue2-root-tfstate-override",
					},
				},
				RemoteStateBackendType: "azurerm",
				RemoteStateBackendConfigs: map[string]any{
					"azurerm": map[string]any{
						"subscription_id": "99999-9999-9999-9999-9999999999",
					},
				},
				Settings: map[string]any{
					"spacelift": map[string]any{
						"workspace_enabled": true,
					},
				},
			},
		},
		"OverrideComponent1": {
			Config: Config{
				Component: "ComponentOverrides",
				Vars: map[string]any{
					"key3": "val-override-3",
					"key4": "val-4",
				},
				Envs: map[string]string{
					"env4": "val-4",
				},
				BackendType: "s3",
			},
		},
		"OverrideComponent2": {
			Config: Config{
				Component: "OverrideComponent1",
				Vars: map[string]any{
					"key4": "val-override-4",
				},
				Envs: map[string]string{
					"env4": "val-override-4",
				},
			},
		},
		"metadata/component": {
			Metadata: Metadata{
				Component: "infra/vpc",
			},
		},
		"test/test-component": {
			Config: Config{
				Vars: map[string]any{
					"enabled": true,
				},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1",
					"TEST_ENV_VAR2": "val2",
					"TEST_ENV_VAR3": "val3",
				},
				Settings: map[string]any{
					"spacelift": map[string]any{
						"workspace_enabled": true,
					},
				},
			},
			Metadata: Metadata{
				Type: "real",
			},
		},
		"test/test-component-override": {
			Config: Config{
				Component: "test/test-component",
				Vars:      map[string]any{},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1-override",
					"TEST_ENV_VAR3": "val3-override",
					"TEST_ENV_VAR4": "val4",
				},
				RemoteStateBackendType: "static",
				RemoteStateBackendConfigs: map[string]any{
					"static": map[string]any{
						"val1": true,
						"val2": "2",
						"val3": 3,
						"val4": nil,
					},
				},
				Settings: map[string]any{
					"spacelift": map[string]any{
						"workspace_enabled": true,
					},
				},
			},
			Metadata: Metadata{
				TerraformWorkspace: "test-component-override-workspace-override",
			},
		},
		"test/test-component-override-2": {
			Config: Config{
				Component: "test/test-component-override",
				Vars:      map[string]any{},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1-override-2",
					"TEST_ENV_VAR2": "val2-override-2",
					"TEST_ENV_VAR4": "val4-override-2",
				},
				RemoteStateBackendType: "static",
				RemoteStateBackendConfigs: map[string]any{
					"static": map[string]any{
						"val1": true,
						"val2": "5",
						"val3": 7,
						"val4": nil,
					},
				},
				Settings: map[string]any{
					"spacelift": map[string]any{
						"workspace_enabled":  true,
						"stack_name_pattern": "{{.tenant}}-{{.environment}}-{{.stage}}-new-component",
					},
				},
			},
			Metadata: Metadata{
				TerraformWorkspacePattern: "{{.tenant}}-{{.environment}}-{{.stage}}-{{.component}}",
			},
		},
		"mixin/test-1": {
			Config: Config{
				Vars: map[string]any{
					"service_1_name": "mixin-1",
				},
			},
			Metadata: Metadata{
				Type: "abstract",
			},
		},
		"mixin/test-2": {
			Config: Config{
				Vars: map[string]any{
					"service_1_name": "mixin-2",
				},
			},
			Metadata: Metadata{
				Type: "abstract",
			},
		},
		"metadata/inherit-1": {
			Metadata: Metadata{
				Component: "test/test-component",
				Inherits: []string{
					"test/test-component-override",
					"test/test-component-override-2",
					"mixin/test-1",
					"mixin/test-2",
				},
			},
		},
	}

	tests := []struct {
		componentName         string
		expectedError         bool
		expectedComponentInfo *ConfigWithMetadata
	}{
		{
			componentName: "MissingComponent",
			expectedError: true,
		},
		{
			componentName: "NoComponentOverride",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "NoComponentOverride",
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val2",
					},
					Envs: map[string]string{
						"env1": "val1",
						"env2": "val2",
					},
					BackendType: "s3",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "s3",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					Settings: map[string]any{},
				},
			},
		},
		{
			componentName: "ComponentOverrides",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "ComponentOverrides",
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val-override-2",
						"key3": "val-3",
					},
					Envs: map[string]string{
						"env1": "val-override-1",
						"env2": "val2",
						"env3": "val3",
					},
					BackendType: "remote",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "azurerm",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					Settings: map[string]any{
						"spacelift": map[string]any{
							"workspace_enabled": true,
						},
					},
				},
			},
		},
		{
			componentName: "OverrideComponent1",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "ComponentOverrides",
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val-override-2",
						"key3": "val-override-3",
						"key4": "val-4",
					},
					Envs: map[string]string{
						"env1": "val-override-1",
						"env2": "val2",
						"env3": "val3",
						"env4": "val-4",
					},
					BackendType: "s3",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "azurerm",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					Settings: map[string]any{
						"spacelift": map[string]any{
							"workspace_enabled": true,
						},
					},
				},
			},
		},
		{
			componentName: "OverrideComponent2",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "ComponentOverrides",
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val-override-2",
						"key3": "val-override-3",
						"key4": "val-override-4",
					},
					Envs: map[string]string{
						"env1": "val-override-1",
						"env2": "val2",
						"env3": "val3",
						"env4": "val-override-4",
					},
					BackendType: "s3",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "azurerm",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					Settings: map[string]any{
						"spacelift": map[string]any{
							"workspace_enabled": true,
						},
					},
				},
			},
		},
		{
			componentName: "metadata/component",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "infra/vpc",
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val2",
					},
					Envs: map[string]string{
						"env1": "val1",
						"env2": "val2",
					},
					BackendType: "s3",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "s3",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					Settings: map[string]any{},
				},
				Metadata: Metadata{
					Component: "infra/vpc",
				},
			},
		},
		{
			componentName: "metadata/inherit-1",
			expectedComponentInfo: &ConfigWithMetadata{
				Config: Config{
					Component: "test/test-component",
					Vars: map[string]any{
						"enabled":        true,
						"key1":           "val1",
						"key2":           "val2",
						"service_1_name": "mixin-2",
					},
					Envs: map[string]string{
						"TEST_ENV_VAR1": "val1-override-2",
						"TEST_ENV_VAR2": "val2-override-2",
						"TEST_ENV_VAR3": "val3-override",
						"TEST_ENV_VAR4": "val4-override-2",
						"env1":          "val1",
						"env2":          "val2",
					},
					BackendType: "s3",
					BackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
					},
					RemoteStateBackendType: "static",
					RemoteStateBackendConfigs: map[string]any{
						"s3": map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						"azurerm": map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.atmos",
						},
						"remote": nil,
						"vault":  nil,
						"static": map[string]any{
							"val1": true,
							"val2": "5",
							"val3": 7,
							"val4": nil,
						},
					},
					Settings: map[string]any{
						"spacelift": map[string]any{
							"workspace_enabled":  true,
							"stack_name_pattern": "{{.tenant}}-{{.environment}}-{{.stage}}-new-component",
						},
					},
				},
				Metadata: Metadata{
					Component: "test/test-component",
					Inherits: []string{
						"test/test-component-override",
						"test/test-component-override-2",
						"mixin/test-1",
						"mixin/test-2",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.componentName, func(t *testing.T) {
			testCase := tc
			t.Parallel()
			componentInfo, err := processComponentConfigs("testStack", baseConfig, componentsConfigMap, testCase.componentName)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.expectedComponentInfo, componentInfo)
		})
	}

}
