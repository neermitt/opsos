package terraform

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

func Clean(ectx ExecutionContext, clearDataDir bool) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), ectx.WorkingDir)
	fmt.Println("Deleting '.terraform' folder")
	err := fs.RemoveAll(".terraform")
	if err != nil {
		return err
	}

	fmt.Println("Deleting '.terraform.lock.hcl' file")
	err = fs.RemoveAll(".terraform.lock.hcl")
	if err != nil {
		return err
	}

	fmt.Printf("Deleting terraform varfile: %s\n", "*.terraform.tfvars*")
	matches, err := afero.Glob(fs, "*.terraform.tfvars*")
	if err != nil {
		return err
	}
	for _, m := range matches {
		err = fs.Remove(m)
		if err != nil {
			return err
		}
	}

	// If `auto_generate_backend_file` is `true` (we are auto-generating backend files), remove `backend.tf.json`
	if ectx.Config.Components.Terraform.AutoGenerateBackendFile {
		fmt.Println("Deleting 'backend.tf*' file")
		matches, err := afero.Glob(fs, "backend.tf*")
		if err != nil {
			return err
		}
		for _, m := range matches {
			err = fs.Remove(m)
			if err != nil {
				return err
			}
		}
	}

	tfDataDir := os.Getenv("TF_DATA_DIR")
	if len(tfDataDir) > 0 && tfDataDir != "." && tfDataDir != "/" && tfDataDir != "./" {
		if clearDataDir {
			fmt.Printf("Found ENV var TF_DATA_DIR=%s", tfDataDir)
			fmt.Printf("Deleting folder '%s'\n", tfDataDir)
			err = fs.RemoveAll(tfDataDir)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
