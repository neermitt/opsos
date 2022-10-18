package terraform

import (
	"fmt"
	"github.com/neermitt/opsos/pkg/config"
	"os"

	"github.com/spf13/afero"
)

func Clean(ectx ExecutionContext, clearDataDir bool) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), ectx.execOptions.WorkingDirectory)
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

	if err := CleanPlanFile(ectx); err != nil {
		return err
	}

	conf := config.GetConfig(ectx.Context)
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

func CleanPlanFile(exeCtx ExecutionContext) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), exeCtx.execOptions.WorkingDirectory)
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
