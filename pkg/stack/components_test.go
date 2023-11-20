package stack

import (
	"testing"

	"github.com/neermitt/opsos/pkg/stack/schema"

	"github.com/stretchr/testify/require"
)

func TestProcessComponentConfigs(t *testing.T) {

	backendTypeS3 := "s3"
	backendTypeRemote := "remote"
	backendTypeAzurerm := "azurerm"
	backendTypeStatic := "static"

	baseConfig := schema.Config{
		Vars: map[string]any{
			"key1": "val1",
			"key2": "val2",
		},
		Envs: map[string]string{
			"env1": "val1",
			"env2": "val2",
		},
		BackendType: &backendTypeS3,
		BackendConfigs: map[string]any{
			backendTypeS3: map[string]any{
				"encrypt":        true,
				"bucket":         "cp-ue2-root-tfstate",
				"key":            "terraform.tfstate",
				"dynamodb_table": "cp-ue2-root-tfstate-lock",
				"acl":            "bucket-owner-full-control",
				"region":         "us-east-2",
				"role_arn":       nil,
			},
			backendTypeAzurerm: map[string]any{
				"subscription_id":      "88888-8888-8888-8888-8888888888",
				"resource_group_name":  "rg-terraform-state",
				"storage_account_name": "staterraformstate",
				"container_name":       "dev-tfstate",
				"key":                  "dev.opsos",
			},
			backendTypeRemote: nil,
			"vault":           nil,
		},
		RemoteStateBackendConfigs: map[string]any{
			backendTypeS3: map[string]any{
				"role_arn": "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
			},
		},
	}

	componentInfraVPC := "infra/vpc"
	componentTestComponent := "test/test-component"
	componentOverrides := "ComponentOverrides"
	componentOverrideComponent1 := "OverrideComponent1"
	componentTestComponentOverride := "test/test-component-override"
	componentNoOverride := "NoComponentOverride"

	componentMetadataTypeReal := "real"
	componentMetadataTypeAbstract := "abstract"
	terraformWorkspaceOverride := "test-component-override-workspace-override"
	terraformWorkspacePattern := "{{.tenant}}-{{.environment}}-{{.stage}}-{{.component}}"

	commandOverride := "/usr/local/bin/terraform"
	componentsConfigMap := map[string]schema.ConfigWithMetadata{
		componentNoOverride: {},
		componentOverrides: {
			Config: schema.Config{
				Command: &commandOverride,
				Vars: map[string]any{
					"key2": "val-override-2",
					"key3": "val-3",
				},
				Envs: map[string]string{
					"env1": "val-override-1",
					"env3": "val3",
				},
				BackendType: &backendTypeRemote,
				BackendConfigs: map[string]any{
					backendTypeS3: map[string]any{
						"bucket": "cp-ue2-root-tfstate-override",
					},
				},
				RemoteStateBackendType: &backendTypeAzurerm,
				RemoteStateBackendConfigs: map[string]any{
					backendTypeAzurerm: map[string]any{
						"subscription_id": "99999-9999-9999-9999-9999999999",
					},
				},
			},
		},
		componentOverrideComponent1: {
			Config: schema.Config{
				Component: &componentOverrides,
				Vars: map[string]any{
					"key3": "val-override-3",
					"key4": "val-4",
				},
				Envs: map[string]string{
					"env4": "val-4",
				},
				BackendType: &backendTypeS3,
			},
		},
		"OverrideComponent2": {
			Config: schema.Config{
				Component: &componentOverrideComponent1,
				Vars: map[string]any{
					"key4": "val-override-4",
				},
				Envs: map[string]string{
					"env4": "val-override-4",
				},
			},
		},
		"metadata/component": {
			Metadata: &schema.Metadata{
				Component: &componentInfraVPC,
			},
		},
		componentTestComponent: {
			Config: schema.Config{
				Vars: map[string]any{
					"enabled": true,
				},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1",
					"TEST_ENV_VAR2": "val2",
					"TEST_ENV_VAR3": "val3",
				},
			},
			Metadata: &schema.Metadata{
				Type: &componentMetadataTypeReal,
			},
		},
		componentTestComponentOverride: {
			Config: schema.Config{
				Component: &componentTestComponent,
				Vars:      map[string]any{},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1-override",
					"TEST_ENV_VAR3": "val3-override",
					"TEST_ENV_VAR4": "val4",
				},
				RemoteStateBackendType: &backendTypeStatic,
				RemoteStateBackendConfigs: map[string]any{
					backendTypeStatic: map[string]any{
						"val1": true,
						"val2": "2",
						"val3": 3,
						"val4": nil,
					},
				},
			},
			Metadata: &schema.Metadata{
				TerraformWorkspace: &terraformWorkspaceOverride,
			},
		},
		"test/test-component-override-2": {
			Config: schema.Config{
				Component: &componentTestComponentOverride,
				Vars:      map[string]any{},
				Envs: map[string]string{
					"TEST_ENV_VAR1": "val1-override-2",
					"TEST_ENV_VAR2": "val2-override-2",
					"TEST_ENV_VAR4": "val4-override-2",
				},
				RemoteStateBackendType: &backendTypeStatic,
				RemoteStateBackendConfigs: map[string]any{
					backendTypeStatic: map[string]any{
						"val1": true,
						"val2": "5",
						"val3": 7,
						"val4": nil,
					},
				},
			},
			Metadata: &schema.Metadata{
				TerraformWorkspacePattern: &terraformWorkspacePattern,
			},
		},
		"mixin/test-1": {
			Config: schema.Config{
				Vars: map[string]any{
					"service_1_name": "mixin-1",
				},
			},
			Metadata: &schema.Metadata{
				Type: &componentMetadataTypeAbstract,
			},
		},
		"mixin/test-2": {
			Config: schema.Config{
				Vars: map[string]any{
					"service_1_name": "mixin-2",
				},
			},
			Metadata: &schema.Metadata{
				Type: &componentMetadataTypeAbstract,
			},
		},
		"metadata/inherit-1": {
			Metadata: &schema.Metadata{
				Component: &componentTestComponent,
				Inherits: []string{
					componentTestComponentOverride,
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
		expectedComponentInfo *schema.ConfigWithMetadata
	}{
		{
			componentName: "MissingComponent",
			expectedError: true,
		},
		{
			componentName: componentNoOverride,
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Component: &componentNoOverride,
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val2",
					},
					Envs: map[string]string{
						"env1": "val1",
						"env2": "val2",
					},
					BackendType: &backendTypeS3,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeS3,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
				},
			},
		},
		{
			componentName: componentOverrides,
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Command:   &commandOverride,
					Component: &componentOverrides,
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
					BackendType: &backendTypeRemote,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeAzurerm,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
				},
			},
		},
		{
			componentName: componentOverrideComponent1,
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Command:   &commandOverride,
					Component: &componentOverrides,
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
					BackendType: &backendTypeS3,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeAzurerm,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
				},
			},
		},
		{
			componentName: "OverrideComponent2",
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Command:   &commandOverride,
					Component: &componentOverrides,
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
					BackendType: &backendTypeS3,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeAzurerm,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate-override",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "99999-9999-9999-9999-9999999999",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
				},
			},
		},
		{
			componentName: "metadata/component",
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Component: &componentInfraVPC,
					Vars: map[string]any{
						"key1": "val1",
						"key2": "val2",
					},
					Envs: map[string]string{
						"env1": "val1",
						"env2": "val2",
					},
					BackendType: &backendTypeS3,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeS3,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
				},
				Metadata: &schema.Metadata{
					Component: &componentInfraVPC,
				},
			},
		},
		{
			componentName: "metadata/inherit-1",
			expectedComponentInfo: &schema.ConfigWithMetadata{
				Config: schema.Config{
					Component: &componentTestComponent,
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
					BackendType: &backendTypeS3,
					BackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       nil,
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
					},
					RemoteStateBackendType: &backendTypeStatic,
					RemoteStateBackendConfigs: map[string]any{
						backendTypeS3: map[string]any{
							"encrypt":        true,
							"bucket":         "cp-ue2-root-tfstate",
							"key":            "terraform.tfstate",
							"dynamodb_table": "cp-ue2-root-tfstate-lock",
							"acl":            "bucket-owner-full-control",
							"region":         "us-east-2",
							"role_arn":       "arn:aws:iam::123456789012:role/cp-gbl-root-terraform",
						},
						backendTypeAzurerm: map[string]any{
							"subscription_id":      "88888-8888-8888-8888-8888888888",
							"resource_group_name":  "rg-terraform-state",
							"storage_account_name": "staterraformstate",
							"container_name":       "dev-tfstate",
							"key":                  "dev.opsos",
						},
						backendTypeRemote: nil,
						"vault":           nil,
						backendTypeStatic: map[string]any{
							"val1": true,
							"val2": "5",
							"val3": 7,
							"val4": nil,
						},
					},
				},
				Metadata: &schema.Metadata{
					Component: &componentTestComponent,
					Inherits: []string{
						componentTestComponentOverride,
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
