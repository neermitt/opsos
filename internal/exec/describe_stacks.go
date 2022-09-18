package exec

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/formatters"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type describeStackOutput struct {
}

type describeStacksOutput struct {
	Stacks map[string]describeStackOutput `yaml:",inline" json:",inline"`
}

// ExecuteDescribeStacks executes `describe stacks` command
func ExecuteDescribeStacks(cmd *cobra.Command, args []string) error {

	flags := cmd.Flags()

	format, err := flags.GetString("format")
	if err != nil {
		return err
	}

	file, err := flags.GetString("file")
	if err != nil {
		return err
	}

	conf := cmd.Context().Value("config").(*config.Configuration)

	stacksBasePath := path.Join(conf.BasePath, conf.Stacks.BasePath)
	stacksBaseAbsPath, err := filepath.Abs(stacksBasePath)
	if err != nil {
		return err
	}

	stackFS := afero.NewBasePathFs(afero.NewOsFs(), stacksBaseAbsPath)

	stackProcessor := stack.NewStackProcessor(stackFS, conf.Stacks.IncludedPaths, conf.Stacks.ExcludedPaths)
	names, err := stackProcessor.GetStackNames()
	if err != nil {
		return err
	}

	stacks, err := stackProcessor.GetStacks(names)
	if err != nil {
		return err
	}

	output := describeStacksOutput{Stacks: make(map[string]describeStackOutput)}

	for _, stk := range stacks {
		output.Stacks[stk.Name()] = describeStackOutput{}
	}

	var w io.Writer = os.Stdout
	if file != "" {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	err = formatters.Get(format)(w, &output)
	if err != nil {
		return err
	}

	return nil
}
