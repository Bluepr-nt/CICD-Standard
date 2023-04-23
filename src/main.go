package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project with a default YAML configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the init command
			fmt.Println("Initializing a new project...")
		},
	}
}

func NewRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run [task-name]",
		Short: "Execute a specific task or all tasks in the pipeline",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the run command
			fmt.Println("Running the specified task or all tasks...")
		},
	}
}

func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of the pipeline, including task status and build artifacts",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the status command
			fmt.Println("Displaying pipeline status...")
		},
	}
}

func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the pipeline configuration against the CI/CD standard",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the validate command
			fmt.Println("Validating pipeline configuration...")
		},
	}
}

func NewRootCommand() *cobra.Command {
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
				fmt.Println("Error reading config file:", err)
			}
		},
	}

	cmd.PersistentFlags().StringP("environment", "e", "development", "The environment in which to run the tasks")
	viper.BindPFlag("environment", cmd.PersistentFlags().Lookup("environment"))

	cmd.AddCommand(NewInitCommand(), NewRunCommand(), NewStatusCommand(), NewValidateCommand())

	return cmd
}

func main() {
	rootCmd := NewRootCommand()
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
