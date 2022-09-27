package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteShellCommand(ctx context.Context, command string, args []string, workingDir string, env []string, dryRun bool) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(cmd.String())

	if dryRun {
		return nil
	}

	return cmd.Run()
}
