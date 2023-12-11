// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package debug

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/jacoblai/arduino-cli/arduino"
	"github.com/jacoblai/arduino-cli/commands/debug"
	"github.com/jacoblai/arduino-cli/commands/sketch"
	"github.com/jacoblai/arduino-cli/i18n"
	"github.com/jacoblai/arduino-cli/patch/cli/arguments"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback/table"
	"github.com/jacoblai/arduino-cli/patch/cli/instance"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tr = i18n.Tr

// NewCommand created a new `upload` command
func NewCommand() *cobra.Command {
	var (
		fqbnArg     arguments.Fqbn
		portArgs    arguments.Port
		interpreter string
		importDir   string
		printInfo   bool
		programmer  arguments.Programmer
	)

	debugCommand := &cobra.Command{
		Use:     "debug",
		Short:   tr("Debug Arduino sketches."),
		Long:    tr("Debug Arduino sketches. (this command opens an interactive gdb session)"),
		Example: "  " + os.Args[0] + " debug -b arduino:samd:mkr1000 -P atmel_ice /home/user/Arduino/MySketch",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runDebugCommand(args, &portArgs, &fqbnArg, interpreter, importDir, &programmer, printInfo)
		},
	}

	debugCommand.AddCommand(newDebugCheckCommand())
	fqbnArg.AddToCommand(debugCommand)
	portArgs.AddToCommand(debugCommand)
	programmer.AddToCommand(debugCommand)
	debugCommand.Flags().StringVar(&interpreter, "interpreter", "console", tr("Debug interpreter e.g.: %s", "console, mi, mi1, mi2, mi3"))
	debugCommand.Flags().StringVarP(&importDir, "input-dir", "", "", tr("Directory containing binaries for debug."))
	debugCommand.Flags().BoolVarP(&printInfo, "info", "I", false, tr("Show metadata about the debug session instead of starting the debugger."))

	return debugCommand
}

func runDebugCommand(args []string, portArgs *arguments.Port, fqbnArg *arguments.Fqbn,
	interpreter string, importDir string, programmer *arguments.Programmer, printInfo bool) {
	instance := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli debug`")

	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	sketchPath := arguments.InitSketchPath(path, true)
	sk, err := sketch.LoadSketch(context.Background(), &rpc.LoadSketchRequest{SketchPath: sketchPath.String()})
	if err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}
	fqbn, port := arguments.CalculateFQBNAndPort(portArgs, fqbnArg, instance, sk.GetDefaultFqbn(), sk.GetDefaultPort(), sk.GetDefaultProtocol())
	debugConfigRequested := &rpc.GetDebugConfigRequest{
		Instance:    instance,
		Fqbn:        fqbn,
		SketchPath:  sketchPath.String(),
		Port:        port,
		Interpreter: interpreter,
		ImportDir:   importDir,
		Programmer:  programmer.String(instance, fqbn),
	}

	if printInfo {

		if res, err := debug.GetDebugConfig(context.Background(), debugConfigRequested); err != nil {
			errcode := feedback.ErrBadArgument
			if errors.Is(err, &arduino.MissingProgrammerError{}) {
				errcode = feedback.ErrMissingProgrammer
			}
			feedback.Fatal(tr("Error getting Debug info: %v", err), errcode)
		} else {
			feedback.PrintResult(newDebugInfoResult(res))
		}

	} else {

		// Intercept SIGINT and forward them to debug process
		ctrlc := make(chan os.Signal, 1)
		signal.Notify(ctrlc, os.Interrupt)

		in, out, err := feedback.InteractiveStreams()
		if err != nil {
			feedback.FatalError(err, feedback.ErrBadArgument)
		}
		if _, err := debug.Debug(context.Background(), debugConfigRequested, in, out, ctrlc); err != nil {
			errcode := feedback.ErrGeneric
			if errors.Is(err, &arduino.MissingProgrammerError{}) {
				errcode = feedback.ErrMissingProgrammer
			}
			feedback.Fatal(tr("Error during Debug: %v", err), errcode)
		}

	}
}

type debugInfoResult struct {
	Executable      string         `json:"executable,omitempty"`
	Toolchain       string         `json:"toolchain,omitempty"`
	ToolchainPath   string         `json:"toolchain_path,omitempty"`
	ToolchainPrefix string         `json:"toolchain_prefix,omitempty"`
	ToolchainConfig any            `json:"toolchain_configuration,omitempty"`
	Server          string         `json:"server,omitempty"`
	ServerPath      string         `json:"server_path,omitempty"`
	ServerConfig    any            `json:"server_configuration,omitempty"`
	SvdFile         string         `json:"svd_file,omitempty"`
	CustomConfigs   map[string]any `json:"custom_configs,omitempty"`
	Programmer      string         `json:"programmer"`
}

type openOcdServerConfigResult struct {
	Path       string   `json:"path,omitempty"`
	ScriptsDir string   `json:"scripts_dir,omitempty"`
	Scripts    []string `json:"scripts,omitempty"`
}

func newDebugInfoResult(info *rpc.GetDebugConfigResponse) *debugInfoResult {
	var toolchainConfig interface{}
	var serverConfig interface{}
	switch info.GetServer() {
	case "openocd":
		var openocdConf rpc.DebugOpenOCDServerConfiguration
		if err := info.GetServerConfiguration().UnmarshalTo(&openocdConf); err != nil {
			feedback.Fatal(tr("Error during Debug: %v", err), feedback.ErrGeneric)
		}
		serverConfig = &openOcdServerConfigResult{
			Path:       openocdConf.GetPath(),
			ScriptsDir: openocdConf.GetScriptsDir(),
			Scripts:    openocdConf.GetScripts(),
		}
	}
	customConfigs := map[string]any{}
	for id, configJson := range info.GetCustomConfigs() {
		var config any
		if err := json.Unmarshal([]byte(configJson), &config); err == nil {
			customConfigs[id] = config
		}
	}
	return &debugInfoResult{
		Executable:      info.GetExecutable(),
		Toolchain:       info.GetToolchain(),
		ToolchainPath:   info.GetToolchainPath(),
		ToolchainPrefix: info.GetToolchainPrefix(),
		ToolchainConfig: toolchainConfig,
		Server:          info.GetServer(),
		ServerPath:      info.GetServerPath(),
		ServerConfig:    serverConfig,
		SvdFile:         info.GetSvdFile(),
		CustomConfigs:   customConfigs,
		Programmer:      info.GetProgrammer(),
	}
}

func (r *debugInfoResult) Data() interface{} {
	return r
}

func (r *debugInfoResult) String() string {
	t := table.New()
	green := color.New(color.FgHiGreen)
	dimGreen := color.New(color.FgGreen)
	t.AddRow(tr("Executable to debug"), table.NewCell(r.Executable, green))
	t.AddRow(tr("Toolchain type"), table.NewCell(r.Toolchain, green))
	t.AddRow(tr("Toolchain path"), table.NewCell(r.ToolchainPath, dimGreen))
	t.AddRow(tr("Toolchain prefix"), table.NewCell(r.ToolchainPrefix, dimGreen))
	if r.SvdFile != "" {
		t.AddRow(tr("SVD file path"), table.NewCell(r.SvdFile, dimGreen))
	}
	switch r.Toolchain {
	case "gcc":
		// no options available at the moment...
	default:
	}
	t.AddRow(tr("Server type"), table.NewCell(r.Server, green))
	t.AddRow(tr("Server path"), table.NewCell(r.ServerPath, dimGreen))

	switch r.Server {
	case "openocd":
		t.AddRow(tr("Configuration options for %s", r.Server))
		openocdConf := r.ServerConfig.(*openOcdServerConfigResult)
		if openocdConf.Path != "" {
			t.AddRow(" - Path", table.NewCell(openocdConf.Path, dimGreen))
		}
		if openocdConf.ScriptsDir != "" {
			t.AddRow(" - Scripts Directory", table.NewCell(openocdConf.ScriptsDir, dimGreen))
		}
		for _, script := range openocdConf.Scripts {
			t.AddRow(" - Script", table.NewCell(script, dimGreen))
		}
	default:
	}
	if custom := r.CustomConfigs; custom != nil {
		for id, config := range custom {
			configJson, _ := json.MarshalIndent(config, "", "  ")
			t.AddRow(tr("Custom configuration for %s:", id))
			return t.Render() + "  " + string(configJson)
		}
	}
	return t.Render()
}
