package exec

import (
	"context"
	"fmt"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
)

type TerraformOptions struct {
	UsePlan     bool
	AutoApprove bool
	Destroy     bool
}

func ExecuteTerraformInit(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := prepareExecCtx(ctx, stackName, component, additionalArgs)
	if err != nil {
		return err
	}

	if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
		return err
	}

	if exeCtx.Config.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}
	return terraform.ExecuteInit(exeCtx, terraform.InitOptions{Reconfigure: exeCtx.Config.Components.Terraform.InitRunReconfigure})
}

func ExecuteTerraformPlan(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := prepareExecCtx(ctx, stackName, component, additionalArgs)
	if err != nil {
		return err
	}
	if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
		return err
	}
	if exeCtx.Config.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}
	if err := terraform.ExecuteInit(exeCtx, terraform.InitOptions{Reconfigure: exeCtx.Config.Components.Terraform.InitRunReconfigure}); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	return terraform.ExecutePlan(exeCtx)
}

func ExecuteTerraformApply(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := prepareExecCtx(ctx, stackName, component, additionalArgs)
	if err != nil {
		return err
	}
	if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
		return err
	}
	if exeCtx.Config.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}
	if err := terraform.ExecuteInit(exeCtx, terraform.InitOptions{Reconfigure: exeCtx.Config.Components.Terraform.InitRunReconfigure}); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	if err := terraform.ExecuteApply(exeCtx, terraform.ApplyOptions{UsePlan: options.UsePlan, AutoApprove: options.AutoApprove, Destroy: options.Destroy}); err != nil {
		return err
	}

	return terraform.CleanPlanFile(exeCtx)
}

func ExecuteTerraformImport(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := prepareExecCtx(ctx, stackName, component, additionalArgs)
	if err != nil {
		return err
	}
	if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
		return err
	}
	if exeCtx.Config.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}
	if err := terraform.ExecuteInit(exeCtx, terraform.InitOptions{Reconfigure: exeCtx.Config.Components.Terraform.InitRunReconfigure}); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	if err := terraform.ExecuteImport(exeCtx); err != nil {
		return err
	}

	return terraform.CleanPlanFile(exeCtx)
}

func prepareExecCtx(ctx context.Context, stackName string, component string, additionalArgs []string) (terraform.ExecutionContext, error) {
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, ComponentType: terraformComponentType, ComponentName: component})
	if err != nil {
		return terraform.ExecutionContext{}, err
	}

	terraformComponents, found := stk.Components[terraformComponentType]
	if !found {
		return terraform.ExecutionContext{}, fmt.Errorf("no terraform component found")
	}
	componentConfig, found := terraformComponents[component]
	if !found {
		return terraform.ExecutionContext{}, fmt.Errorf("terraform component %s not found", component)
	}

	if err != nil {
		return terraform.ExecutionContext{}, err
	}

	conf := config.GetConfig(ctx)
	workingDir, _, err := getComponentWorkingDirectory(conf, terraformComponentType, componentConfig)
	if err != nil {
		return terraform.ExecutionContext{}, err
	}

	exeCtx := terraform.ExecutionContext{
		Context:         ctx,
		Config:          conf,
		Stack:           stk,
		ComponentName:   component,
		ComponentConfig: componentConfig,
		WorkingDir:      workingDir,
		DryRun:          false,
		AdditionalArgs:  additionalArgs,
	}
	return exeCtx, nil
}
