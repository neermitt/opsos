package terraform

import (
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
)

type ExecutionContext struct {
	Config          *config.Configuration
	Stack           *stack.Stack
	ComponentName   string
	ComponentConfig stack.ConfigWithMetadata
	WorkingDir      string
	DryRun          bool
}
