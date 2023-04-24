package initialize

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewInitCommand(input io.Reader, output io.Writer) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project with a default YAML configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			interactive, err := cmd.Flags().GetBool("interactive")
			if err != nil {
				return fmt.Errorf("error reading interactive flag: %v", err)
			}

			config := ""

			if interactive {
				if input == nil {
					input = os.Stdin
				}

				config = generateInteractiveConfig(input, cmd)
			} else {
				config = getDefaultConfig()
			}

			err = writeConfigFile(config)
			if err != nil {
				return fmt.Errorf("error writing configuration file: %v", err)
			}

			cmd.Println("Project initialized with configuration file: ccs.yaml")
			return nil
		},
	}

	initCmd.Flags().BoolP("interactive", "i", false, "Use interactive mode to configure the project")

	return initCmd
}

func getDefaultConfig() string {
	return `apiVersion: v1
kind: Pipeline
spec:
  tasks:
  - name: build
    type: build
    build:
      environment: my-docker-image
      command: docker build $REPOSITORY_DIR
  - name: release
    type: release
    release:
      environment: my-release-environment
      command: release-command
  - name: deployment
    type: deployment
    deployment:
      environment: my-deployment-environment
      command: deployment-command
`
}

func generateInteractiveConfig(input io.Reader, cmd *cobra.Command) string {
	reader := bufio.NewReader(os.Stdin)
	if input != nil {
		reader = bufio.NewReader(input)
	}

	cmd.Println("Enter task information for the build task:")
	buildEnvironment := prompt(reader, "Environment:")
	buildCommand := prompt(reader, "Command:")

	cmd.Println("Enter task information for the release task:")
	releaseEnvironment := prompt(reader, "Environment:")
	releaseCommand := prompt(reader, "Command:")

	cmd.Println("Enter task information for the deployment task:")
	deploymentEnvironment := prompt(reader, "Environment:")
	deploymentCommand := prompt(reader, "Command:")

	return fmt.Sprintf(`apiVersion: v1
kind: Pipeline
spec:
  tasks:
  - name: build
    type: build
    build:
      environment: %s
      command: %s
  - name: release
    type: release
    release:
      environment: %s
      command: %s
  - name: deployment
    type: deployment
    deployment:
      environment: %s
      command: %s
`, buildEnvironment, buildCommand, releaseEnvironment, releaseCommand, deploymentEnvironment, deploymentCommand)
}

func prompt(reader *bufio.Reader, question string) string {
	fmt.Print(question + " ")
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}

func writeConfigFile(config string) error {
	return os.WriteFile("ccs.yaml", []byte(config), 0644)
}
