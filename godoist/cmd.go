package godoist

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godoist",
	Short: "go and doist",
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tasks for a given project name.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		display := NewCobraDisplay(cmd, args)
		projectName := display.Prompt("Project name")
		tasks, err := FetchTasksByProjectName(projectName)
		if err != nil {
			display.Error(err.Error())
			os.Exit(1)
		}
		display.Debug(tasks)
	},
}
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new task to a specified project. Use --project flag or select interactively",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		display := NewCobraDisplay(cmd, args)
		projectName, _ := cmd.Flags().GetString("project")

		// If no project flag, prompt for project with Inbox pre-filled
		if projectName == "" {
			projectName = display.PromptWithDefault("Project name", "Inbox")
		}

		parsedTask := display.PromptForTask(projectName)
		_, err := CreateTaskInProjectWithDue(parsedTask.Content, projectName, parsedTask.DueString)
		if err != nil {
			display.Error(err.Error())
			os.Exit(1)
		}
		// Task created successfully - exit silently
	},
}
var useGUI bool = false

func init() {
	LoadConfig()

	rootCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")
	listCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")
	addCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")
	addCmd.Flags().StringP("project", "p", "", "Project name to add task to")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
