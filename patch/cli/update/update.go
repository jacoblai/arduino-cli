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

package update

import (
	"os"

	"github.com/jacoblai/arduino-cli/i18n"
	"github.com/jacoblai/arduino-cli/patch/cli/core"
	"github.com/jacoblai/arduino-cli/patch/cli/instance"
	"github.com/jacoblai/arduino-cli/patch/cli/lib"
	"github.com/jacoblai/arduino-cli/patch/cli/outdated"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tr = i18n.Tr

// NewCommand creates a new `update` command
func NewCommand() *cobra.Command {
	var showOutdated bool
	updateCommand := &cobra.Command{
		Use:     "update",
		Short:   tr("Updates the index of cores and libraries"),
		Long:    tr("Updates the index of cores and libraries to the latest versions."),
		Example: "  " + os.Args[0] + " update",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateCommand(showOutdated)
		},
	}
	updateCommand.Flags().BoolVar(&showOutdated, "show-outdated", false, tr("Show outdated cores and libraries after index update"))
	return updateCommand
}

func runUpdateCommand(showOutdated bool) {
	inst := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli update`")
	lib.UpdateIndex(inst)
	core.UpdateIndex(inst)
	instance.Init(inst)
	if showOutdated {
		outdated.Outdated(inst)
	}
}
