package terraform

import "github.com/neermitt/opsos/pkg/stack"

type ExecutionContext struct {
	Stack           *stack.Stack
	ComponentName   string
	ComponentConfig stack.ConfigWithMetadata
	WorkingDir      string
	DryRun          bool
}
