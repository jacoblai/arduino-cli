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

package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jacoblai/arduino-cli/configuration"
	"github.com/jacoblai/arduino-cli/i18n"
	"github.com/jacoblai/arduino-cli/patch/cli/board"
	"github.com/jacoblai/arduino-cli/patch/cli/burnbootloader"
	"github.com/jacoblai/arduino-cli/patch/cli/cache"
	"github.com/jacoblai/arduino-cli/patch/cli/compile"
	"github.com/jacoblai/arduino-cli/patch/cli/completion"
	"github.com/jacoblai/arduino-cli/patch/cli/config"
	"github.com/jacoblai/arduino-cli/patch/cli/core"
	"github.com/jacoblai/arduino-cli/patch/cli/daemon"
	"github.com/jacoblai/arduino-cli/patch/cli/debug"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback"
	"github.com/jacoblai/arduino-cli/patch/cli/generatedocs"
	"github.com/jacoblai/arduino-cli/patch/cli/lib"
	"github.com/jacoblai/arduino-cli/patch/cli/monitor"
	"github.com/jacoblai/arduino-cli/patch/cli/outdated"
	"github.com/jacoblai/arduino-cli/patch/cli/sketch"
	"github.com/jacoblai/arduino-cli/patch/cli/update"
	"github.com/jacoblai/arduino-cli/patch/cli/updater"
	"github.com/jacoblai/arduino-cli/patch/cli/upgrade"
	"github.com/jacoblai/arduino-cli/patch/cli/upload"
	"github.com/jacoblai/arduino-cli/patch/cli/version"
	"github.com/jacoblai/arduino-cli/patch/inventory"
	versioninfo "github.com/jacoblai/arduino-cli/version"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	semver "go.bug.st/relaxed-semver"
)

var (
	verbose            bool
	outputFormat       string
	configFile         string
	updaterMessageChan chan *semver.Version = make(chan *semver.Version)
)

// NewCommand creates a new ArduinoCli command root
func NewCommand() *cobra.Command {
	cobra.AddTemplateFunc("tr", i18n.Tr)

	// ArduinoCli is the root command
	arduinoCli := &cobra.Command{
		Use:               "arduino-cli",
		Short:             tr("Arduino CLI."),
		Long:              tr("Arduino Command Line Interface (arduino-cli)."),
		Example:           fmt.Sprintf("  %s <%s> [%s...]", os.Args[0], tr("command"), tr("flags")),
		PersistentPreRun:  preRun,
		PersistentPostRun: postRun,
	}

	arduinoCli.SetUsageTemplate(getUsageTemplate())

	createCliCommandTree(arduinoCli)

	return arduinoCli
}

// this is here only for testing
func createCliCommandTree(cmd *cobra.Command) {
	cmd.AddCommand(board.NewCommand())
	cmd.AddCommand(cache.NewCommand())
	cmd.AddCommand(compile.NewCommand())
	cmd.AddCommand(completion.NewCommand())
	cmd.AddCommand(config.NewCommand())
	cmd.AddCommand(core.NewCommand())
	cmd.AddCommand(daemon.NewCommand())
	cmd.AddCommand(generatedocs.NewCommand())
	cmd.AddCommand(lib.NewCommand())
	cmd.AddCommand(monitor.NewCommand())
	cmd.AddCommand(outdated.NewCommand())
	cmd.AddCommand(sketch.NewCommand())
	cmd.AddCommand(update.NewCommand())
	cmd.AddCommand(upgrade.NewCommand())
	cmd.AddCommand(upload.NewCommand())
	cmd.AddCommand(debug.NewCommand())
	cmd.AddCommand(burnbootloader.NewCommand())
	cmd.AddCommand(version.NewCommand())
	cmd.AddCommand(feedback.NewCommand())

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, tr("Print the logs on the standard output."))
	cmd.Flag("verbose").Hidden = true
	cmd.PersistentFlags().BoolVar(&verbose, "log", false, tr("Print the logs on the standard output."))
	validLogLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	cmd.PersistentFlags().String("log-level", "", tr("Messages with this level and above will be logged. Valid levels are: %s", strings.Join(validLogLevels, ", ")))
	cmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return validLogLevels, cobra.ShellCompDirectiveDefault
	})
	cmd.PersistentFlags().String("log-file", "", tr("Path to the file where logs will be written."))
	validLogFormats := []string{"text", "json"}
	cmd.PersistentFlags().String("log-format", "", tr("The output format for the logs, can be: %s", strings.Join(validLogFormats, ", ")))
	cmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return validLogFormats, cobra.ShellCompDirectiveDefault
	})
	validOutputFormats := []string{"text", "json", "jsonmini", "yaml"}
	cmd.PersistentFlags().StringVar(&outputFormat, "format", "text", tr("The output format for the logs, can be: %s", strings.Join(validOutputFormats, ", ")))
	cmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return validOutputFormats, cobra.ShellCompDirectiveDefault
	})
	cmd.PersistentFlags().StringVar(&configFile, "config-file", "", tr("The custom config file (if not specified the default will be used)."))
	cmd.PersistentFlags().StringSlice("additional-urls", []string{}, tr("Comma-separated list of additional URLs for the Boards Manager."))
	cmd.PersistentFlags().Bool("no-color", false, "Disable colored output.")
	configuration.BindFlags(cmd, configuration.Settings)
}

// convert the string passed to the `--log-level` option to the corresponding
// logrus formal level.
func toLogLevel(s string) (t logrus.Level, found bool) {
	t, found = map[string]logrus.Level{
		"trace": logrus.TraceLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}[s]

	return
}

func preRun(cmd *cobra.Command, args []string) {
	configFile := configuration.Settings.ConfigFileUsed()

	// initialize inventory
	err := inventory.Init(configuration.DataDir(configuration.Settings).String())
	if err != nil {
		feedback.Fatal(fmt.Sprintf("Error: %v", err), feedback.ErrInitializingInventory)
	}

	// https://no-color.org/
	color.NoColor = configuration.Settings.GetBool("output.no_color") || os.Getenv("NO_COLOR") != ""

	// Set default feedback output to colorable
	feedback.SetOut(colorable.NewColorableStdout())
	feedback.SetErr(colorable.NewColorableStderr())

	updaterMessageChan = make(chan *semver.Version)
	go func() {
		if cmd.Name() == "version" {
			// The version command checks by itself if there's a new available version
			updaterMessageChan <- nil
		}
		// Starts checking for updates
		currentVersion, err := semver.Parse(versioninfo.VersionInfo.VersionString)
		if err != nil {
			updaterMessageChan <- nil
		}
		updaterMessageChan <- updater.CheckForUpdate(currentVersion)
	}()

	//
	// Prepare logging
	//

	// decide whether we should log to stdout
	if verbose {
		// if we print on stdout, do it in full colors
		logrus.SetOutput(colorable.NewColorableStdout())
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: color.NoColor,
		})
	} else {
		logrus.SetOutput(io.Discard)
	}

	// set the Logger format
	logFormat := strings.ToLower(configuration.Settings.GetString("logging.format"))
	if logFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// should we log to file?
	logFile := configuration.Settings.GetString("logging.file")
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			feedback.Fatal(tr("Unable to open file for logging: %s", logFile), feedback.ErrGeneric)
		}

		// we use a hook so we don't get color codes in the log file
		if logFormat == "json" {
			logrus.AddHook(lfshook.NewHook(file, &logrus.JSONFormatter{}))
		} else {
			logrus.AddHook(lfshook.NewHook(file, &logrus.TextFormatter{}))
		}
	}

	// configure logging filter
	if lvl, found := toLogLevel(configuration.Settings.GetString("logging.level")); !found {
		feedback.Fatal(tr("Invalid option for --log-level: %s", configuration.Settings.GetString("logging.level")), feedback.ErrBadArgument)
	} else {
		logrus.SetLevel(lvl)
	}

	//
	// Prepare the Feedback system
	//

	// check the right output format was passed
	format, found := feedback.ParseOutputFormat(outputFormat)
	if !found {
		feedback.Fatal(tr("Invalid output format: %s", outputFormat), feedback.ErrBadArgument)
	}

	// use the output format to configure the Feedback
	feedback.SetFormat(format)

	//
	// Print some status info and check command is consistent
	//

	if configFile != "" {
		logrus.Infof("Using config file: %s", configFile)
	} else {
		logrus.Info("Config file not found, using default values")
	}

	logrus.Info(versioninfo.VersionInfo.Application + " version " + versioninfo.VersionInfo.VersionString)

	if outputFormat != "text" {
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			logrus.Warn("Calling help on JSON format")
			feedback.Fatal(tr("Invalid Call : should show Help, but it is available only in TEXT mode."), feedback.ErrBadArgument)
		})
	}
}

func postRun(cmd *cobra.Command, args []string) {
	latestVersion := <-updaterMessageChan
	if latestVersion != nil {
		// Notify the user a new version is available
		updater.NotifyNewVersionIsAvailable(latestVersion.String())
	}
}
