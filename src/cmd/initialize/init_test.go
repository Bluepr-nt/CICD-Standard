package initialize

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInitCommand_DefaultConfig(t *testing.T) {
	// Prepare command
	cmd := NewInitCommand(nil)
	cmd.Flags().Set("interactive", "false")

	// Capture command output
	outputBuffer := new(bytes.Buffer)
	cmd.SetOut(outputBuffer)

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err)

	// Check output message
	expectedOutput := "Project initialized with configuration file: ccs.yaml\n"
	assert.Equal(t, expectedOutput, outputBuffer.String())

	// Check if the file was created and remove it after testing
	configContent, err := ioutil.ReadFile("ccs.yaml")
	assert.NoError(t, err)
	os.Remove("ccs.yaml")

	defaultConfig := getDefaultConfig()
	assert.Equal(t, defaultConfig, string(configContent))
}

func TestNewInitCommand_InteractiveConfig(t *testing.T) {
	// Prepare command
	input := strings.NewReader("my-docker-image\nmy-build-command\nmy-release-environment\nmy-release-command\nmy-deployment-environment\nmy-deployment-command\n")
	cmd := NewInitCommand(input)
	cmd.Flags().Set("interactive", "true")

	// Capture command output
	outputBuffer := new(bytes.Buffer)
	cmd.SetOut(outputBuffer)

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err)

	// Check output message
	expectedOutput := "Enter task information for the build task:\nEnter task information for the release task:\nEnter task information for the deployment task:\nProject initialized with configuration file: ccs.yaml\n"
	assert.Equal(t, expectedOutput, outputBuffer.String())

	// Check if the file was created and remove it after testing
	configContent, err := ioutil.ReadFile("ccs.yaml")
	assert.NoError(t, err)
	os.Remove("ccs.yaml")

	// Build expected interactive config content
	expectedConfigContent := `apiVersion: v1
kind: Pipeline
spec:
  tasks:
  - name: build
    type: build
    build:
      environment: my-docker-image
      command: my-build-command
  - name: release
    type: release
    release:
      environment: my-release-environment
      command: my-release-command
  - name: deployment
    type: deployment
    deployment:
      environment: my-deployment-environment
      command: my-deployment-command
`

	assert.Equal(t, expectedConfigContent, string(configContent))
}
