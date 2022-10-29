package components

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	component1 = ComponentConfig{
		ApiVersion: "opsos/Component",
		Kind:       "Component",
		Metadata: ComponentMetadata{
			Name:        "TestComponent1",
			Description: "test component 1",
		},
		Spec: ComponentSpec{
			Source: VendorComponentSource{
				Uri:     "github.com/cloudposse/terraform-aws-components.git//modules/account-map?ref={{.Version}}",
				Version: "0.196.1",
				IncludedPaths: []string{
					"**/*.tf",
					"**/*.tfvars",
					"**/*.md",
					"**/modules/**",
					"**/modules/**/*.tf",
					"**/modules/**/*.tfvars",
					"**/modules/**/*.md",
				},
				ExcludedPaths: []string{},
			},
			Mixins: nil,
		},
	}
)

func TestGetComponent(t *testing.T) {
	tmpDir := t.TempDir()
	err := PrepareComponentBySpec(context.Background(), tmpDir, component1.Spec)
	require.NoError(t, err)
}
