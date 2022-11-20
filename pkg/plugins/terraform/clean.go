package terraform

import (
	"context"
	"log"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/afero"
)

func Clean(ctx context.Context, clearDataDir bool) error {
	execOptions := utils.GetExecOptions(ctx)

	fs := afero.NewBasePathFs(afero.NewOsFs(), execOptions.WorkingDirectory)
	log.Println("[INFO] (terraform) Deleting '.terraform' folder")
	if err := fs.RemoveAll(".terraform"); err != nil {
		return err
	}

	log.Printf("[INFO] (terraform) Deleting lock file %s", ".terraform.lock.hcl")
	if err := fs.RemoveAll(".terraform.lock.hcl"); err != nil {
		return err
	}

	log.Printf("[INFO] (terraform) Deleting terraform varfile: %s", "*.terraform.tfvars*")
	if err := removeAllFiles(fs, "*.terraform.tfvars*"); err != nil {
		return err
	}

	if err := cleanPlanFileFromComponentDir(execOptions.WorkingDirectory); err != nil {
		return err
	}

	conf := config.GetConfig(ctx)
	var terraformConfig Config
	err := utils.FromMap(conf.Providers[ComponentType], &terraformConfig)
	if err != nil {
		return err
	}

	// If `auto_generate_backend_file` is `true` (we are auto-generating backend files), remove `backend.tf.json`
	if terraformConfig.AutoGenerateBackendFile {
		log.Printf("[INFO] (terraform) Deleting backend file: %s", "backend.tf*")
		if err := removeAllFiles(fs, "backend.tf*"); err != nil {
			return err
		}
	}

	tfDataDir := os.Getenv("TF_DATA_DIR")
	if len(tfDataDir) > 0 && tfDataDir != "." && tfDataDir != "/" && tfDataDir != "./" {
		if clearDataDir {
			log.Printf("[INFO] (terraform) Found ENV var: %s=%s", "TF_DATA_DIR", tfDataDir)
			log.Printf("[INFO] (terraform) Deleting folder '%s'", tfDataDir)
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
	log.Printf("[INFO] (terraform) Deleting terraform plan: %s", "*.planfile")
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
