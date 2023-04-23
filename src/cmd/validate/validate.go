package validate

import (
	"ccs/pkg/cicd"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("file", "f", "ccs.yaml", "Path to the pipeline configuration file")
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the pipeline configuration against the CI/CD standard",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error reading file flag:", err)
			os.Exit(1)
		}

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		var pipeline cicd.Pipeline
		err = yaml.Unmarshal(fileContent, &pipeline)
		if err != nil {
			fmt.Println("Error unmarshalling YAML:", err)
			os.Exit(1)
		}

		if err := validatePipeline(&pipeline); err != nil {
			fmt.Println("Pipeline validation failed:", err)
			os.Exit(1)
		}

		fmt.Println("Pipeline validation succeeded")
	},
}

func validatePipeline(p *cicd.Pipeline) error {
	taskNames := make(map[string]bool)

	for _, task := range p.Spec.Tasks {
		if _, exists := taskNames[task.Name]; exists {
			return fmt.Errorf("duplicate task name: %s", task.Name)
		}
		taskNames[task.Name] = true

		if err := validateTask(&task); err != nil {
			return fmt.Errorf("task validation failed: %w", err)
		}
	}

	return nil
}

func validateTask(t *cicd.Task) error {
	if t.Name == "" {
		return fmt.Errorf("task name must be set")
	}

	if t.Needs != nil {
		for _, dependency := range t.Needs {
			if dependency == "" {
				return fmt.Errorf("task %q has an empty dependency", t.Name)
			}
		}
	}

	return nil
}
