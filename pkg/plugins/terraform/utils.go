package terraform

import (
	"fmt"
	"strings"

	"github.com/neermitt/opsos/pkg/stack"
)

func constructVarfileName(stack *stack.Stack, componentName string) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(componentName, "/", "-")
	return fmt.Sprintf("%s-%s.terraform.tfvars.json", stack.Name, fmtdComponentFolderPrefix)
}

// constructPlanfileName constructs the planfile name for a terraform component in a stack
func constructPlanfileName(stack *stack.Stack, componentName string) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(componentName, "/", "-")
	return fmt.Sprintf("%s-%s.planfile", stack.Name, fmtdComponentFolderPrefix)
}
