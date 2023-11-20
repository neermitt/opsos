package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

var (
	describeStackOptins exec.DescribeStackOptions

	stackDescribeCmd = &cobra.Command{
		Use:   "describe [<stack>]",
		Short: "Execute 'stack describe' command",
		Long:  `This command describes the stack components: opsos stack describe [<stack>]`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				describeStackOptins.Stack = args[0]
			}
			return exec.ExecuteDescribeStacks(cmd, describeStackOptins)
		},
	}
)

// stackDescribeCmd command describes the stack

func init() {
	stackDescribeCmd.PersistentFlags().StringVar(&describeStackOptins.OutputFile, "file", "", "Write the result to file: opsos describe stacks --file=stacks.yaml")
	stackDescribeCmd.PersistentFlags().StringVar(&describeStackOptins.Format, "format", "yaml", "Specify output format: opsos describe stacks --format=yaml/json ('yaml' is default)")
	stackDescribeCmd.PersistentFlags().StringArrayVar(&describeStackOptins.Components, "components", nil, "Filter by specific components: opsos describe stacks --components=<component1>,<component2>")
	stackDescribeCmd.PersistentFlags().StringArrayVar(&describeStackOptins.ComponentTypes, "component-types", nil, "Filter by specific component types: opsos describe stacks --component-types=terraform,helmfile, Available component types: terraform, helmfile")
	stackDescribeCmd.PersistentFlags().StringArrayVar(&describeStackOptins.PrintSections, "sections", nil, "Output only these component sections: opsos describe stacks --sections=vars,settings. Available component sections: backend, backend_type, deps, env, inheritance, metadata, remote_state_backend, remote_state_backend_type, settings, vars")

	stackCmd.AddCommand(stackDescribeCmd)
}
