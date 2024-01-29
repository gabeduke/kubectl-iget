package cmdbuilder

import (
	"context"
	"fmt"
	"strings"

	execute "github.com/alexellis/go-execute/v2"
)

type CommandBuilder struct {
	// Add additional fields if needed
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

func (cb *CommandBuilder) BuildAndExecuteCommand(resourceType string, fields []string, filters []string, dryRun bool) (string, error) {
	args := []string{"get", resourceType}

	// Append fields and filters to arguments
	if len(fields) > 0 {
		args = append(args, fmt.Sprintf("--custom-columns=%s", strings.Join(fields, ",")))
	}
	if len(filters) > 0 {
		args = append(args, fmt.Sprintf("--field-selector=%s", strings.Join(filters, ",")))
	}

	// Handle dry-run option
	if dryRun {
		return fmt.Sprintf("kubectl %s", strings.Join(args, " ")), nil
	}

	// Execute the command
	cmd := execute.ExecTask{
		Command:     "kubectl",
		Args:        args,
		StreamStdio: false,
	}

	res, err := cmd.Execute(context.Background())
	if err != nil {
		return "", fmt.Errorf("error executing command: %v", err)
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf("non-zero exit code: %s", res.Stderr)
	}

	return res.Stdout, nil
}
