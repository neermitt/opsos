package exec

import (
	"fmt"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/cobra"
)

type DescribeStackOptions struct {
	Format         string
	OutputFile     string
	Stack          string
	Components     []string
	ComponentTypes []string
	PrintSections  []string
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

	stackNames, err := stackProcessor.GetStackNames()
	if err != nil {
		return err
	}

	getStackOptions := stack.GetStackOptions{
		Components:     options.Components,
		ComponentTypes: options.ComponentTypes,
	}

	var stacks []*stack.Stack
	if options.Stack != "" {
		if !utils.StringInSlice(options.Stack, stackNames) {
			return fmt.Errorf("stack %s not found", options.Stack)
		} else {
			stk, err := stackProcessor.GetStack(options.Stack, getStackOptions)
			if err != nil {
				return err
			}
			stacks = []*stack.Stack{stk}
		}
	} else {
		stacks, err = stackProcessor.GetStacks(stackNames, getStackOptions)
		if err != nil {
			return err
		}
	}

	output := describeStacksOutput{Stacks: make(map[string]describeStackOutput)}

	for _, stk := range stacks {
		filterAbstractComponents(stk)
		output.Stacks[stk.Id] = describeStackOutput{
			Name:       stk.Name,
			Components: filterComponentSections(stk.Components, options.PrintSections),
		}
	}

	err = utils.PrintOrWriteToFile(options.Format, options.OutputFile, &output, 0644)
	if err != nil {
		return err
	}

	return nil
}

func filterComponentSections(components map[string]stack.ComponentConfigMap, sections []string) map[string]stack.ComponentConfigMap {
	if len(sections) == 0 {
		return components
	}
	var filteredResult = make(map[string]stack.ComponentConfigMap, len(components))
	for componentType, componentMap := range components {
		filteredResult[componentType] = filterComponentMapForSections(componentMap, sections)
	}

	return filteredResult
}

func filterComponentMapForSections(componentMap stack.ComponentConfigMap, sections []string) stack.ComponentConfigMap {
	var filteredResult = make(stack.ComponentConfigMap, len(componentMap))
	for componentName, componentConfig := range componentMap {
		filteredResult[componentName] = filterComponentConfigForSections(componentConfig, sections)
	}
	return filteredResult
}

func filterComponentConfigForSections(sourceConfig stack.ConfigWithMetadata, sections []string) stack.ConfigWithMetadata {
	var destinationConfig stack.ConfigWithMetadata
	for _, section := range sections {
		switch section {
		case "vars":
			destinationConfig.Vars = sourceConfig.Vars
		case "env":
			destinationConfig.Envs = sourceConfig.Envs
		case "backend_type":
			destinationConfig.BackendType = sourceConfig.BackendType
		case "backend":
			destinationConfig.Backend = sourceConfig.Backend
		case "remote_state_backend":
			destinationConfig.RemoteStateBackend = sourceConfig.RemoteStateBackend
		case "remote_state_backend_type":
			destinationConfig.RemoteStateBackendType = sourceConfig.RemoteStateBackendType
		case "settings":
			destinationConfig.Settings = sourceConfig.Settings
		case "metadata":
			destinationConfig.Metadata = sourceConfig.Metadata
		}
	}

	return destinationConfig
}

func filterAbstractComponents(stk *stack.Stack) {
	for ct, components := range stk.Components {
		filteredComponents := make(map[string]stack.ConfigWithMetadata, len(components))
		for componentName, c := range components {
			if c.Metadata != nil && c.Metadata.Type != nil && *c.Metadata.Type == "abstract" {
				continue
			}
			filteredComponents[componentName] = c
		}
		stk.Components[ct] = filteredComponents
	}
}
