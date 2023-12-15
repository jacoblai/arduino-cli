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

package board

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/jacoblai/arduino-cli/commands/board"
	"github.com/jacoblai/arduino-cli/inter/cli/feedback"
	"github.com/jacoblai/arduino-cli/inter/cli/instance"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/jacoblai/arduino-cli/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var showHiddenBoard bool

func initListAllCommand() *cobra.Command {
	var listAllCommand = &cobra.Command{
		Use:   fmt.Sprintf("listall [%s]", tr("boardname")),
		Short: tr("List all known boards and their corresponding FQBN."),
		Long: tr(`List all boards that have the support platform installed. You can search
for a specific board if you specify the board name`),
		Example: "" +
			"  " + os.Args[0] + " board listall\n" +
			"  " + os.Args[0] + " board listall zero",
		Args: cobra.ArbitraryArgs,
		Run:  runListAllCommand,
	}
	listAllCommand.Flags().BoolVarP(&showHiddenBoard, "show-hidden", "a", false, tr("Show also boards marked as 'hidden' in the platform"))
	return listAllCommand
}

// runListAllCommand list all installed boards
func runListAllCommand(cmd *cobra.Command, args []string) {
	inst := instance.CreateAndInit()

	logrus.Info("Executing `arduino-cli board listall`")

	list, err := board.ListAll(context.Background(), &rpc.BoardListAllRequest{
		Instance:            inst,
		SearchArgs:          args,
		IncludeHiddenBoards: showHiddenBoard,
	})
	if err != nil {
		feedback.Fatal(tr("Error listing boards: %v", err), feedback.ErrGeneric)
	}

	feedback.PrintResult(resultAll{list})
}

// output from this command requires special formatting, let's create a dedicated
// feedback.Result implementation
type resultAll struct {
	list *rpc.BoardListAllResponse
}

func (dr resultAll) Data() interface{} {
	return dr.list
}

func (dr resultAll) String() string {
	sort.Slice(dr.list.Boards, func(i, j int) bool {
		return dr.list.Boards[i].GetName() < dr.list.Boards[j].GetName()
	})

	t := table.New()
	t.SetHeader(tr("Board Name"), tr("FQBN"), "")
	for _, item := range dr.list.GetBoards() {
		hidden := ""
		if item.IsHidden {
			hidden = tr("(hidden)")
		}
		t.AddRow(item.GetName(), item.GetFqbn(), hidden)
	}
	return t.Render()
}
