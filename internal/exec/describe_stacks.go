package exec

import (
	"fmt"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/cobra"
)

type DescribeStackOptions struct {
	Format     string
	OutputFile string
	Stack      string
}

type describeStackOutput struct {
	Name       string
	Components map[string]stack.ComponentConfigMap
}

type describeStacksOutput struct {
	Stacks map[string]describeStackOutput `yaml:",inline" json:",inline"`
}

// ExecuteDescribeStacks executes `describe stacks` command
func ExecuteDescribeStacks(cmd *cobra.Command, options DescribeStackOptions) error {
	ctx := cmd.Context()
	conf := config.GetConfig(ctx)

	stackProcessor, err := stack.NewStackProcessorFromConfig(conf)
	if err != nil {
		return err
	}

	stackNames, err := stackProcessor.GetStackNames(ctx)
	if err != nil {
		return err
	}

	var stacks []*stack.Stack
	if options.Stack != "" {
		if !utils.StringInSlice(options.Stack, stackNames) {
			return fmt.Errorf("stack %s not found", options.Stack)
		} else {
			stk, err := stackProcessor.GetStack(ctx, options.Stack)
			if err != nil {
				return err
			}
			stacks = []*stack.Stack{stk}
		}
	} else {
		stacks, err = stackProcessor.GetStacks(ctx, stackNames)
		if err != nil {
			return err
		}
	}

	output := describeStacksOutput{Stacks: make(map[string]describeStackOutput)}

	for _, stk := range stacks {
		output.Stacks[stk.Id] = describeStackOutput{
			Name:       stk.Name,
			Components: stk.ComponentTypes,
		}
	}

	err = utils.PrintOrWriteToFile(options.Format, options.OutputFile, &output, 0644)
	if err != nil {
		return err
	}

	return nil
}
