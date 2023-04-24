package validate

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidateCommand(t *testing.T) {
	tests := []struct {
		name           string
		configFilename string
		configContent  string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "valid pipeline",
			configFilename: "valid_pipeline.yaml",
			configContent: `apiVersion: v1
kind: Pipeline
spec:
  tasks:
  - name: build_app_a
    type: build
    needs: []
    build:
      environment: my-docker-image
      command: docker build $REPOSITORY_DIR`,
			expectedOutput: "Pipeline validation succeeded\n",
		},
		{
			name:           "invalid pipeline",
			configFilename: "invalid_pipeline.yaml",
			configContent: `apiVersion: v1
kind: Pipeline
spec:
  tasks:
  - name: build_app_a
    type: build
    needs: []
    build:
      environment: my-docker-image
      command: docker build $REPOSITORY_DIR
  - name: build_app_a
    type: build
    needs: []
    build:
      environment: my-docker-image
      command: docker build $REPOSITORY_DIR`,
			expectedOutput: `Error: pipeline validation failed: duplicate task name: build_app_a
Usage:
  validate [flags]

Flags:
  -f, --file string   Path to the pipeline configuration file (default "ccs.yaml")
  -h, --help          help for validate

`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for config files
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, tt.configFilename)

			err := ioutil.WriteFile(configPath, []byte(tt.configContent), 0644)
			assert.NoError(t, err)

			validateCmd := NewValidateCommand(nil)
			validateCmd.Flags().Set("file", configPath)

			output := new(bytes.Buffer)
			validateCmd.SetOutput(output)

			err = validateCmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}

func TestNewValidateCommand_NeedsValidation(t *testing.T) {
	// Create a temporary YAML file with an invalid 'Needs' field
	fileContent := []byte(`
spec:
  tasks:
    - name: build_app_a
      type: build
      build:
        environment: my-docker-image
        command: docker build $REPOSITORY_DIR
    - name: build_app_b
      type: build
      needs: ["non_existent_task"]
      build:
        environment: my-docker-image
        command: docker build $REPOSITORY_DIR
`)

	tmpfile, err := ioutil.TempFile("", "ccs.yaml")
	if err != nil {
		t.Fatal("Error creating temporary file:", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(fileContent); err != nil {
		t.Fatal("Error writing to temporary file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Error closing temporary file:", err)
	}

	// Create the Validate command
	validateCmd := NewValidateCommand(nil)
	validateCmd.Flags().Set("file", tmpfile.Name())

	// Capture the output of the command
	outputBuffer := new(bytes.Buffer)
	validateCmd.SetOutput(outputBuffer)

	err = validateCmd.Execute()

	// Check if the command fails with the expected error
	if err == nil {
		t.Error("Expected the command to fail due to invalid 'Needs' field, but it did not")
	} else {
		expectedError := "pipeline validation failed: task validation failed: task \"build_app_b\" has an invalid dependency: \"non_existent_task\""
		if err.Error() != expectedError {
			t.Errorf("Expected error: %q, got: %q", expectedError, err.Error())
		}
	}
}

func TestNewValidateCommand_TaskTypeValidation(t *testing.T) {
	// Create a temporary YAML file with an invalid 'type' field
	fileContent := []byte(`
spec:
  tasks:
    - name: build_app_a
      type: invalid_type
      build:
        environment: my-docker-image
        command: docker build $REPOSITORY_DIR
`)

	tmpfile, err := ioutil.TempFile("", "ccs.yaml")
	if err != nil {
		t.Fatal("Error creating temporary file:", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(fileContent); err != nil {
		t.Fatal("Error writing to temporary file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Error closing temporary file:", err)
	}

	validateCmd := NewValidateCommand(nil)
	validateCmd.Flags().Set("file", tmpfile.Name())

	outputBuffer := new(bytes.Buffer)
	validateCmd.SetOutput(outputBuffer)

	err = validateCmd.Execute()

	expectedError := "pipeline validation failed: task validation failed: task type validation failed: invalid task type: invalid_type, allowed types: [build release deployment]"
	if assert.Error(t, err, "Expected the command to fail due to invalid 'type' field, but it did not") {
		assert.EqualError(t, err, expectedError, "Expected error: %q, got: %q", expectedError, err.Error())
	}
}

func TestNewValidateCommand_MissingTaskType(t *testing.T) {
	fileContent := []byte(`
spec:
  tasks:
    - name: build_app_a
      build:
        environment: my-docker-image
        command: docker build $REPOSITORY_DIR
`)

	tmpfile, err := ioutil.TempFile("", "ccs.yaml")
	if err != nil {
		t.Fatal("Error creating temporary file:", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(fileContent); err != nil {
		t.Fatal("Error writing to temporary file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Error closing temporary file:", err)
	}

	validateCmd := NewValidateCommand(nil)
	validateCmd.Flags().Set("file", tmpfile.Name())

	outputBuffer := new(bytes.Buffer)
	validateCmd.SetOutput(outputBuffer)

	err = validateCmd.Execute()

	expectedError := "pipeline validation failed: task validation failed: task type validation failed: invalid task type: , allowed types: [build release deployment]"
	if assert.Error(t, err, "Expected the command to fail due to missing task type, but it did not") {
		assert.EqualError(t, err, expectedError, "Expected error: %q, got: %q", expectedError, err.Error())
	}
}
