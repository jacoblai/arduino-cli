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
	"strings"

	"github.com/arduino/arduino-cli/commands/board"
	"github.com/arduino/arduino-cli/patch/cli/feedback"
	"github.com/arduino/arduino-cli/patch/cli/feedback/result"
	"github.com/arduino/arduino-cli/patch/cli/feedback/table"
	"github.com/arduino/arduino-cli/patch/cli/instance"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initSearchCommand() *cobra.Command {
	var searchCommand = &cobra.Command{
		Use:   fmt.Sprintf("search [%s]", tr("boardname")),
		Short: tr("Search for a board in the Boards Manager."),
		Long:  tr(`Search for a board in the Boards Manager using the specified keywords.`),
		Example: "" +
			"  " + os.Args[0] + " board search\n" +
			"  " + os.Args[0] + " board search zero",
		Args: cobra.ArbitraryArgs,
		Run:  runSearchCommand,
	}
	searchCommand.Flags().BoolVarP(&showHiddenBoard, "show-hidden", "a", false, tr("Show also boards marked as 'hidden' in the platform"))
	return searchCommand
}

func runSearchCommand(cmd *cobra.Command, args []string) {
	inst := instance.CreateAndInit()

	logrus.Info("Executing `arduino-cli board search`")

	res, err := board.Search(context.Background(), &rpc.BoardSearchRequest{
		Instance:            inst,
		SearchArgs:          strings.Join(args, " "),
		IncludeHiddenBoards: showHiddenBoard,
	})
	if err != nil {
		feedback.Fatal(tr("Error searching boards: %v", err), feedback.ErrGeneric)
	}

	feedback.PrintResult(searchResults{result.NewBoardListItems(res.GetBoards())})
}

// output from this command requires special formatting so we create a dedicated
// feedback.Result implementation
type searchResults struct {
	Boards []*result.BoardListItem `json:"boards"`
}

func (r searchResults) Data() interface{} {
	return r
}

func (r searchResults) String() string {
	if len(r.Boards) == 0 {
		return ""
	}

	t := table.New()
	t.SetHeader(tr("Board Name"), tr("FQBN"), tr("Platform ID"), "")

	sort.Slice(r.Boards, func(i, j int) bool {
		return r.Boards[i].Name < r.Boards[j].Name
	})

	for _, item := range r.Boards {
		hidden := ""
		if item.IsHidden {
			hidden = tr("(hidden)")
		}
		t.AddRow(item.Name, item.Fqbn, item.Platform.Metadata.Id, hidden)
	}
	return t.Render()
}
