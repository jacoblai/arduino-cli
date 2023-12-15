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

package core

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jacoblai/arduino-cli/arduino"
	"github.com/jacoblai/arduino-cli/commands/core"
	"github.com/jacoblai/arduino-cli/inter/cli/arguments"
	"github.com/jacoblai/arduino-cli/inter/cli/feedback"
	"github.com/jacoblai/arduino-cli/inter/cli/instance"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initUpgradeCommand() *cobra.Command {
	var postInstallFlags arguments.PostInstallFlags
	upgradeCommand := &cobra.Command{
		Use:   fmt.Sprintf("upgrade [%s:%s] ...", tr("PACKAGER"), tr("ARCH")),
		Short: tr("Upgrades one or all installed platforms to the latest version."),
		Long:  tr("Upgrades one or all installed platforms to the latest version."),
		Example: "" +
			"  # " + tr("upgrade everything to the latest version") + "\n" +
			"  " + os.Args[0] + " core upgrade\n\n" +
			"  # " + tr("upgrade arduino:samd to the latest version") + "\n" +
			"  " + os.Args[0] + " core upgrade arduino:samd",
		Run: func(cmd *cobra.Command, args []string) {
			runUpgradeCommand(args, postInstallFlags.DetectSkipPostInstallValue())
		},
	}
	postInstallFlags.AddToCommand(upgradeCommand)
	return upgradeCommand
}

func runUpgradeCommand(args []string, skipPostInstall bool) {
	inst := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli core upgrade`")
	Upgrade(inst, args, skipPostInstall)
}

// Upgrade upgrades one or all installed platforms to the latest version.
func Upgrade(inst *rpc.Instance, args []string, skipPostInstall bool) {
	// if no platform was passed, upgrade allthethings
	if len(args) == 0 {
		targets, err := core.PlatformList(&rpc.PlatformListRequest{
			Instance:      inst,
			UpdatableOnly: true,
		})
		if err != nil {
			feedback.Fatal(tr("Error retrieving core list: %v", err), feedback.ErrGeneric)
		}

		if len(targets.InstalledPlatforms) == 0 {
			feedback.Print(tr("All the cores are already at the latest version"))
			return
		}

		for _, t := range targets.InstalledPlatforms {
			args = append(args, t.Id)
		}
	}

	warningMissingIndex := func(response *rpc.PlatformUpgradeResponse) {
		if response == nil || response.Platform == nil {
			return
		}
		if !response.Platform.Indexed {
			feedback.Warning(tr("missing package index for %s, future updates cannot be guaranteed", response.Platform.Id))
		}
	}

	// proceed upgrading, if anything is upgradable
	platformsRefs, err := arguments.ParseReferences(args)
	if err != nil {
		feedback.Fatal(tr("Invalid argument passed: %v", err), feedback.ErrBadArgument)
	}

	hasBadArguments := false
	for i, platformRef := range platformsRefs {
		if platformRef.Version != "" {
			feedback.Warning(tr("Invalid item %s", args[i]))
			hasBadArguments = true
			continue
		}

		r := &rpc.PlatformUpgradeRequest{
			Instance:        inst,
			PlatformPackage: platformRef.PackageName,
			Architecture:    platformRef.Architecture,
			SkipPostInstall: skipPostInstall,
		}
		response, err := core.PlatformUpgrade(context.Background(), r, feedback.ProgressBar(), feedback.TaskProgress())
		warningMissingIndex(response)
		if err != nil {
			var alreadyAtLatestVersionErr *arduino.PlatformAlreadyAtTheLatestVersionError
			if errors.As(err, &alreadyAtLatestVersionErr) {
				feedback.Warning(err.Error())
				continue
			}

			feedback.Fatal(tr("Error during upgrade: %v", err), feedback.ErrGeneric)
		}
	}

	if hasBadArguments {
		feedback.Fatal(tr("Some upgrades failed, please check the output for details."), feedback.ErrBadArgument)
	}

	feedback.PrintResult(&platformUpgradeResult{})
}

// This is needed so we can print warning messages in case users use --format json
type platformUpgradeResult struct{}

// Data implements feedback.Result.
func (r *platformUpgradeResult) Data() interface{} {
	return r
}

// String implements feedback.Result.
func (r *platformUpgradeResult) String() string {
	return ""
}
