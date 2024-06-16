package godoist

import (
	"fmt"
	"log"

	"github.com/alecthomas/repr"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
)

type CobraDisplay interface {
	Prompt(msg string) string
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
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")

	b, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	l, err := gtk.LabelNew("Hello, gotk3 2!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	b.Add(l)

	e, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}

	b.Add(e)

	win.Add(b)

	win.SetDefaultSize(800, 600)
	win.SetKeepAbove(true)
	win.SetSkipTaskbarHint(true)

	win.ShowAll()
	win.Present()
	win.GrabFocus()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()

	// resp, err := zenity.Entry(msg)
	// if err != nil {
	// 	zenity.Error(err.Error())
	// 	os.Exit(1)
	// }
	// return resp
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
