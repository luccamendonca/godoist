package godoist

import (
	"fmt"
	"os"

	"github.com/alecthomas/repr"
	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
)

type CobraDisplay interface {
	Prompt(msg string) string
	PromptWithDefault(msg string, defaultValue string) string
	PromptForTask(projectName string) string
	Error(msg string)
	Info(msg string)
	Debug(any interface{})
}

type DisplayCLI struct {
	cmd  *cobra.Command
	args []string
}
type DisplayGUI struct {
	cmd  *cobra.Command
	args []string
}

func NewCobraDisplay(cmd *cobra.Command, args []string) CobraDisplay {
	useGUI := cmd.Flag("use-gui").Value.String() == "true"
	if useGUI {
		return DisplayGUI{cmd, args}
	}
	return DisplayCLI{cmd, args}
}

// DisplayCLI
func (cli DisplayCLI) Prompt(msg string) string {
	return cli.args[0]
}
func (cli DisplayCLI) PromptWithDefault(msg string, defaultValue string) string {
	return defaultValue
}
func (cli DisplayCLI) PromptForTask(projectName string) string {
	if len(cli.args) > 0 {
		return cli.args[0]
	}
	return ""
}
func (cli DisplayCLI) Error(msg string) {
	cli.Info(msg)
}
func (cli DisplayCLI) Info(msg string) {
	fmt.Println(msg)
}
func (cli DisplayCLI) Debug(any interface{}) {
	repr.Println(any)
}

// DisplayGUI
func (gui DisplayGUI) Prompt(msg string) string {
	resp, err := zenity.Entry(msg)
	if err != nil {
		// User cancelled - exit silently
		os.Exit(0)
	}
	return resp
}
func (gui DisplayGUI) PromptWithDefault(msg string, defaultValue string) string {
	resp, err := zenity.Entry(msg, zenity.EntryText(defaultValue))
	if err != nil {
		// User cancelled - exit silently
		os.Exit(0)
	}
	return resp
}
func (gui DisplayGUI) PromptForTask(projectName string) string {
	msg := fmt.Sprintf("Task name (Project: %s)", projectName)
	resp, err := zenity.Entry(msg)
	if err != nil {
		// User cancelled - exit silently
		os.Exit(0)
	}
	return resp
}
func (gui DisplayGUI) Error(msg string) {
	zenity.Error(msg)
}
func (gui DisplayGUI) Info(msg string) {
	zenity.Info(msg)
}
func (gui DisplayGUI) Debug(any interface{}) {
	zenity.Info(repr.String(any))
}
