package godoist

import (
	"fmt"
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
	Short: "Adds a new task, given a task name. Creates on Inbox by default",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		display := NewCobraDisplay(cmd, args)
		taskName := display.Prompt("Task name")
		task, err := CreateTask(taskName)
		if err != nil {
			display.Error(err.Error())
			os.Exit(1)
		}
		display.Info(fmt.Sprintf("Task created! Id: %d", task.Id))
	},
}
var useGUI bool = false

func init() {
	LoadConfig()

	rootCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")
	listCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")
	addCmd.Flags().BoolVarP(&useGUI, "use-gui", "g", false, "Uses GUI instead of CLI")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
