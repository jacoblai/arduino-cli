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

package lib

import (
	"context"
	"os"

	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/patch/cli/feedback"
	"github.com/arduino/arduino-cli/patch/cli/instance"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initUpdateIndexCommand() *cobra.Command {
	updateIndexCommand := &cobra.Command{
		Use:     "update-index",
		Short:   tr("Updates the libraries index."),
		Long:    tr("Updates the libraries index to the latest version."),
		Example: "  " + os.Args[0] + " lib update-index",
		Args:    cobra.NoArgs,
		Run:     runUpdateIndexCommand,
	}
	return updateIndexCommand
}

func runUpdateIndexCommand(cmd *cobra.Command, args []string) {
	inst := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli lib update-index`")
	UpdateIndex(inst)
}

// UpdateIndex updates the index of libraries.
func UpdateIndex(inst *rpc.Instance) {
	err := commands.UpdateLibrariesIndex(context.Background(), &rpc.UpdateLibrariesIndexRequest{
		Instance: inst,
	}, feedback.ProgressBar())
	if err != nil {
		feedback.Fatal(tr("Error updating library index: %v", err), feedback.ErrGeneric)
	}
}
