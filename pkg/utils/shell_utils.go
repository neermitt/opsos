package utils

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/neermitt/opsos/pkg/logging"
)

type ExecOptions struct {
	DryRun           bool
	Env              []string
	WorkingDirectory string
	StdOut           io.Writer
}

func ExecuteShellCommand(ctx context.Context, command string, args []string, options ExecOptions) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = append(os.Environ(), options.Env...)
	cmd.Dir = options.WorkingDirectory
	cmd.Stdin = os.Stdin
	if options.StdOut != nil {
		cmd.Stdout = options.StdOut
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	log.Printf("[TRACE] Executing command: %s", logging.Indent(logging.Indent(cmd.String())))

	if options.DryRun {
		return nil
	}

	return cmd.Run()
}

func SetExecOptions(ctx context.Context, component ExecOptions) context.Context {
	return context.WithValue(ctx, "exec-options", component)
}

func GetExecOptions(ctx context.Context) ExecOptions {
	return ctx.Value("exec-options").(ExecOptions)
}
