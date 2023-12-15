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
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/fatih/color"
	"github.com/jacoblai/arduino-cli/commands/lib"
	"github.com/jacoblai/arduino-cli/inter/cli/arguments"
	"github.com/jacoblai/arduino-cli/inter/cli/feedback"
	"github.com/jacoblai/arduino-cli/inter/cli/instance"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	fqbn arguments.Fqbn
)

func initExamplesCommand() *cobra.Command {
	examplesCommand := &cobra.Command{
		Use:     fmt.Sprintf("examples [%s]", tr("LIBRARY_NAME")),
		Short:   tr("Shows the list of the examples for libraries."),
		Long:    tr("Shows the list of the examples for libraries. A name may be given as argument to search a specific library."),
		Example: "  " + os.Args[0] + " lib examples Wire",
		Args:    cobra.MaximumNArgs(1),
		Run:     runExamplesCommand,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return arguments.GetInstalledLibraries(), cobra.ShellCompDirectiveDefault
		},
	}
	fqbn.AddToCommand(examplesCommand)
	return examplesCommand
}

func runExamplesCommand(cmd *cobra.Command, args []string) {
	instance := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli lib examples`")

	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	res, err := lib.LibraryList(context.Background(), &rpc.LibraryListRequest{
		Instance: instance,
		All:      true,
		Name:     name,
		Fqbn:     fqbn.String(),
	})
	if err != nil {
		feedback.Fatal(tr("Error getting libraries info: %v", err), feedback.ErrGeneric)
	}

	found := []*libraryExamples{}
	for _, lib := range res.GetInstalledLibraries() {
		found = append(found, &libraryExamples{
			Library:  lib.Library,
			Examples: lib.Library.Examples,
		})
	}

	feedback.PrintResult(libraryExamplesResult{found})
	logrus.Info("Done")
}

// output from this command requires special formatting, let's create a dedicated
// feedback.Result implementation

type libraryExamples struct {
	Library  *rpc.Library `json:"library"`
	Examples []string     `json:"examples"`
}

type libraryExamplesResult struct {
	Examples []*libraryExamples
}

func (ir libraryExamplesResult) Data() interface{} {
	return ir.Examples
}

func (ir libraryExamplesResult) String() string {
	if ir.Examples == nil || len(ir.Examples) == 0 {
		return tr("No libraries found.")
	}

	sort.Slice(ir.Examples, func(i, j int) bool {
		return strings.ToLower(ir.Examples[i].Library.Name) < strings.ToLower(ir.Examples[j].Library.Name)
	})

	res := []string{}
	for _, lib := range ir.Examples {
		name := lib.Library.Name
		if lib.Library.ContainerPlatform != "" {
			name += " (" + lib.Library.GetContainerPlatform() + ")"
		} else if lib.Library.Location != rpc.LibraryLocation_LIBRARY_LOCATION_USER {
			name += " (" + lib.Library.GetLocation().String() + ")"
		}
		r := tr("Examples for library %s", color.GreenString("%s", name)) + "\n"
		sort.Slice(lib.Examples, func(i, j int) bool {
			return strings.ToLower(lib.Examples[i]) < strings.ToLower(lib.Examples[j])
		})
		for _, example := range lib.Examples {
			examplePath := paths.New(example)
			r += fmt.Sprintf("  - %s%s\n",
				color.New(color.Faint).Sprintf("%s%c", examplePath.Parent(), os.PathSeparator),
				examplePath.Base())
		}
		res = append(res, r)
	}

	return strings.Join(res, "\n")
}
