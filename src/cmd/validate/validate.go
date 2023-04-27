package validate

import (
	"ccs/pkg/cicd"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewValidateCommand() *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the pipeline configuration against the CI/CD standard",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, err := cmd.Flags().GetString("file")
			if err != nil {
				return fmt.Errorf("error reading file flag: %v", err)
			}

			fileContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}

			var pipeline cicd.Pipeline
			err = yaml.Unmarshal(fileContent, &pipeline)
			if err != nil {
				return fmt.Errorf("error unmarshalling YAML: %v", err)
			}

			if err := validatePipeline(&pipeline); err != nil {
				return fmt.Errorf("pipeline validation failed: %v", err)
			}

			cmd.Println("Pipeline validation succeeded")
			return nil
		},
	}

	validateCmd.Flags().StringP("file", "f", "ccs.yaml", "Path to the pipeline configuration file")

	return validateCmd
}

func validatePipeline(p *cicd.Pipeline) error {
	if p.Spec == nil {
		return fmt.Errorf("pipeline spec must not be empty")
	}

	if len(p.Spec.Tasks) == 0 {
		return fmt.Errorf("tasks must not be empty")
	}

	taskNames := make(map[string]*cicd.Task)

	for _, task := range p.Spec.Tasks {
		if _, exists := taskNames[task.Name]; exists {
			return fmt.Errorf("duplicate task name: %s", task.Name)
		}
		taskNames[task.Name] = &task

		if err := validateTask(&task, taskNames); err != nil {
			return fmt.Errorf("task validation failed: %w", err)
		}
	}

	return nil
}

func validateTask(t *cicd.Task, tasks map[string]*cicd.Task) error {
	if t.Name == "" {
		return fmt.Errorf("task name must be set")
	}

	if err := validateTaskType(t.Type); err != nil {
		return fmt.Errorf("task type validation failed: %w", err)
	}

	if t.Needs != nil {
		for _, dependency := range t.Needs {
			if _, exists := tasks[dependency]; !exists {
				return fmt.Errorf("task %q has an invalid dependency: %q", t.Name, dependency)
			}
		}
	}

	return nil
}

func validateTaskType(taskType string) error {
	allowedTaskTypes := []string{"build", "release", "deployment"}

	for _, allowedType := range allowedTaskTypes {
		if taskType == allowedType {
			return nil
		}
	}

	return fmt.Errorf("invalid task type: %s, allowed types: %v", taskType, allowedTaskTypes)
}
