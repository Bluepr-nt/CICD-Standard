package cicd

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PipelineSpec is the spec for a CI/CD pipeline
type PipelineSpec struct {
	ProductData ProductData `json:"product_data"`
	Tasks       []Task      `json:"tasks"`
}

// ProductData represents the product information
type ProductData struct {
	Name string `json:"name"`
}

// TaskAction represents a task action, either Build, Release or Deployment
type TaskAction interface{}

// Build represents a build task configuration
type Build struct {
	Environment string `json:"environment"`
	Command     string `json:"command"`
}

// Release represents a release task configuration
type Release struct {
	Level    string                 `json:"level"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Deployment represents a deployment task configuration
type Deployment struct {
	Environment string `json:"environment"`
	Release     struct {
		Task string `json:"task"`
	} `json:"release"`
}

// Task represents a CI/CD task
type Task struct {
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	Needs  []string   `json:"needs"`
	Action TaskAction `json:"-"`
}

// UnmarshalJSON is a custom JSON unmarshaller for Task
func (t *Task) UnmarshalJSON(data []byte) error {
	type taskAlias Task
	var taskJSON struct {
		taskAlias
		Build      *Build      `json:"build,omitempty"`
		Release    *Release    `json:"release,omitempty"`
		Deployment *Deployment `json:"deployment,omitempty"`
	}

	if err := json.Unmarshal(data, &taskJSON); err != nil {
		return err
	}

	*t = Task(taskJSON.taskAlias)

	definedCount := 0
	if taskJSON.Build != nil {
		t.Action = taskJSON.Build
		definedCount++
	}
	if taskJSON.Release != nil {
		t.Action = taskJSON.Release
		definedCount++
	}
	if taskJSON.Deployment != nil {
		t.Action = taskJSON.Deployment
		definedCount++
	} else {
		return fmt.Errorf("task %q must have one action type defined", t.Name)
	}

	if definedCount > 1 {
		return fmt.Errorf("task %q has more than one action type defined", t.Name)
	}

	return nil
}

// Pipeline is the top-level Kubernetes-style struct
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineSpec `json:"spec"`
}
