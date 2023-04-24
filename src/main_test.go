package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func(cmd *cobra.Command, args []string)) string {
	// Save the original standard output
	origStdout := os.Stdout

	// Create a pipe to capture the output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a buffer to store the captured output
	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		io.Copy(&buf, r)
		done <- struct{}{}
	}()

	// Execute the function with the modified standard output
	f(nil, nil)

	// Close the pipe and restore the original standard output
	w.Close()
	os.Stdout = origStdout
	<-done

	return buf.String()
}

func TestRunCommandOutput(t *testing.T) {
	runCmd := NewRunCommand(nil)
	output := new(bytes.Buffer)
	runCmd.SetOutput(output)

	expectedOutput := "Running the specified task or all tasks...\n"

	if err := runCmd.Execute(); err != nil {
		t.Errorf("No error expected, got: %v", err)
	}
	if output.String() != expectedOutput {
		t.Errorf("Expected output: %q, got: %q", expectedOutput, output)
	}
}

func TestStatusCommandOutput(t *testing.T) {
	statusCmd := NewStatusCommand(nil)
	output := new(bytes.Buffer)
	statusCmd.SetOutput(output)

	expectedOutput := "Displaying pipeline status...\n"

	if err := statusCmd.Execute(); err != nil {
		t.Errorf("No error expected, got: %v", err)
	}

	if output.String() != expectedOutput {
		t.Errorf("Expected output: %q, got: %q", expectedOutput, output)
	}
}

func TestRootCommand(t *testing.T) {
	rootCmd := NewRootCommand(nil)

	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `CICD tool helps manage CI/CD pipelines by providing an interface
to create, manage, and execute tasks based on a standardized YAML configuration.

Usage:
	cicd [command]

Available Commands:
	completion  Generate the autocompletion script for the specified shell
	help        Help about any command
	init        Initialize a new project with a default YAML configuration file
	run         Execute a specific task or all tasks in the pipeline
	status      Show the status of the pipeline, including task status and build artifacts
	validate    Validate the pipeline configuration against the CI/CD standard

Flags:
	-h, --help   help for cicd

Use "cicd [command] --help" for more information about a command.
`

	assert.Equal(t, output.String(), expected, "Output does not match expected value. Got:\n%s\n\nExpected:\n%s", output, expected)

}
