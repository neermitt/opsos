package components_test

import (
	"context"
	"testing"

	"github.com/neermitt/opsos/api/common"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/components"
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

func TestPrepareComponentBySpec(t *testing.T) {
	tmpDir := t.TempDir()
	err := components.PrepareComponentBySpec(context.Background(), tmpDir, tmpDir, component1.Spec, components.PrepareComponentOptions{DryRun: false})
	require.NoError(t, err)
}
