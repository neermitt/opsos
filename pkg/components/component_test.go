package components_test

import (
	"context"
	"os"
	"testing"

	"github.com/neermitt/opsos/api/common"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/components"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	component1 = v1.Component{
		Object: common.Object{
			ApiVersion: "opsos/v1",
			Kind:       "Component",
			Metadata: common.ObjectMetadata{
				Name:        "TestComponent1",
				Description: "test component 1",
			},
		},
		Spec: v1.ComponentSpec{
			Source: v1.ComponentSource{
				Uri:     "github.com/cloudposse/terraform-aws-components.git//modules/account-map?ref={{.Version}}",
				Version: "0.196.1",
				IncludedPaths: []string{
					"**/*.tf",
					"**/*.tfvars",
					"**/*.md",
				},
			},
			Mixins: nil,
		},
	}
)

func TestComponentRead(t *testing.T) {
	file, err := os.Open("../../examples/complete/components/terraform/infra/account-map/component.yaml")
	require.NoError(t, err)

	defer file.Close()

	component, err := components.ReadComponent(file)
	require.NoError(t, err)

	assert.Equal(t, &v1.Component{
		Object: common.Object{
			ApiVersion: "opsos/v1",
			Kind:       "Component",
			Metadata: common.ObjectMetadata{
				Name:        "account-map-vendor-config",
				Description: "Source and mixins config for building 'vpc-flow-logs-bucket' component",
			},
		},
		Spec: v1.ComponentSpec{
			Source: v1.ComponentSource{
				Uri:     "github.com/cloudposse/terraform-aws-components.git//modules/account-map?ref={{.Version}}",
				Version: "0.196.1",
				IncludedPaths: []string{
					"**/*.tf",
					"**/*.tfvars",
					"**/*.md",
				},
			},
			Mixins: []v1.ComponentMixins{
				{
					Uri:      "https://raw.githubusercontent.com/cloudposse/terraform-aws-components/{{.Version}}/modules/datadog-agent/introspection.mixin.tf",
					Version:  "0.196.1",
					Filename: "introspection.mixin.tf",
				},
			},
		},
	}, component)
}

func TestPrepareComponentBySpec(t *testing.T) {
	tmpDir := t.TempDir()
	err := components.PrepareComponentBySpec(context.Background(), tmpDir, tmpDir, component1.Spec, components.PrepareComponentOptions{DryRun: false})
	require.NoError(t, err)
}
