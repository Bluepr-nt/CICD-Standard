package main

import (
	"ccs/cmd/initialize"
	"ccs/cmd/validate"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRunCommand(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "run [task-name]",
		Short: "Execute a specific task or all tasks in the pipeline",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the run command
			cmd.Println("Running the specified task or all tasks...")
		},
	}
}

func NewStatusCommand(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of the pipeline, including task status and build artifacts",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the status command
			cmd.Println("Displaying pipeline status...")
		},
	}
}

func NewRootCommand(output io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cicd",
		Short: "CICD tool for managing CI/CD pipelines",
		Long: `CICD tool helps manage CI/CD pipelines by providing an interface
to create, manage, and execute tasks based on a standardized YAML configuration.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Bind environment variables
			viper.AutomaticEnv()

			// Read the configuration file
			viper.SetConfigName("ccs")
			viper.AddConfigPath(".")

			if err := viper.ReadInConfig(); err != nil {
				cmd.Println("Warning reading config file:", err)
			}
		},
	}

	cmd.AddCommand(initialize.NewInitCommand(nil, output), NewRunCommand(output), NewStatusCommand(output), validate.NewValidateCommand(output))

	return cmd
}

func main() {
	rootCmd := NewRootCommand(nil)
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
