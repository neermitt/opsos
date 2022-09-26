package exec

import (
	"fmt"
	"io"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/formatters"
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
	Components map[string]stack.ComponentConfigMap
}

type describeStacksOutput struct {
	Stacks map[string]describeStackOutput `yaml:",inline" json:",inline"`
}

// ExecuteDescribeStacks executes `describe stacks` command
func ExecuteDescribeStacks(cmd *cobra.Command, options DescribeStackOptions) error {
	conf := config.GetConfig(cmd.Context())

	stackProcessor, err := stack.NewStackProcessorFromConfig(conf)
	if err != nil {
		return err
	}

	stackNames, err := stackProcessor.GetStackNames()
	if err != nil {
		return err
	}

	var stacks []*stack.Stack
	if options.Stack != "" {
		if !utils.StringInSlice(options.Stack, stackNames) {
			return fmt.Errorf("stack %s not found", options.Stack)
		} else {
			stk, err := stackProcessor.GetStack(options.Stack)
			if err != nil {
				return err
			}
			stacks = []*stack.Stack{stk}
		}
	} else {
		stacks, err = stackProcessor.GetStacks(stackNames)
		if err != nil {
			return err
		}
	}

	output := describeStacksOutput{Stacks: make(map[string]describeStackOutput)}

	for _, stk := range stacks {
		output.Stacks[stk.Name] = describeStackOutput{
			Components: stk.ComponentTypes,
		}
	}

	var w io.Writer = os.Stdout
	if options.OutputFile != "" {
		f, err := os.OpenFile(options.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	err = formatters.Get(options.Format)(w, &output)
	if err != nil {
		return err
	}

	return nil
}
