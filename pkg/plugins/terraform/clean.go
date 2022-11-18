package terraform

import (
	"context"
	"fmt"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/afero"
)

func Clean(ctx context.Context, clearDataDir bool) error {
	execOptions := utils.GetExecOptions(ctx)

	fs := afero.NewBasePathFs(afero.NewOsFs(), execOptions.WorkingDirectory)
	fmt.Println("Deleting '.terraform' folder")
	if err := fs.RemoveAll(".terraform"); err != nil {
		return err
	}

	fmt.Println("Deleting '.terraform.lock.hcl' file")
	if err := fs.RemoveAll(".terraform.lock.hcl"); err != nil {
		return err
	}

	fmt.Printf("Deleting terraform varfile: %s\n", "*.terraform.tfvars*")
	if err := removeAllFiles(fs, "*.terraform.tfvars*"); err != nil {
		return err
	}

	if err := cleanPlanFileFromComponentDir(execOptions.WorkingDirectory); err != nil {
		return err
	}

	conf := config.GetConfig(ctx)
	// If `auto_generate_backend_file` is `true` (we are auto-generating backend files), remove `backend.tf.json`
	if conf.Components.Terraform.AutoGenerateBackendFile {
		fmt.Println("Deleting 'backend.tf*' file")
		if err := removeAllFiles(fs, "backend.tf*"); err != nil {
			return err
		}
	}

	tfDataDir := os.Getenv("TF_DATA_DIR")
	if len(tfDataDir) > 0 && tfDataDir != "." && tfDataDir != "/" && tfDataDir != "./" {
		if clearDataDir {
			fmt.Printf("Found ENV var TF_DATA_DIR=%s", tfDataDir)
			fmt.Printf("Deleting folder '%s'\n", tfDataDir)
			if err := fs.RemoveAll(tfDataDir); err != nil {
				return err
			}
		}

	}
	return nil
}

func CleanPlanFile(ctx context.Context) error {
	execOptions := utils.GetExecOptions(ctx)
	return cleanPlanFileFromComponentDir(execOptions.WorkingDirectory)
}

func cleanPlanFileFromComponentDir(componentDir string) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), componentDir)
	fmt.Printf("Deleting terraform plan: %s\n", "*.planfile")
	return removeAllFiles(fs, "*.planfile")
}

func removeAllFiles(fs afero.Fs, pattern string) error {
	matches, err := afero.Glob(fs, pattern)
	for _, m := range matches {
		err = fs.Remove(m)
		if err != nil {
			return err
		}
	}
	return nil
}
