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

package outdated

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jacoblai/arduino-cli/i18n"
	"github.com/jacoblai/arduino-cli/patch/cli/core"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback/result"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback/table"
	"github.com/jacoblai/arduino-cli/patch/cli/instance"
	"github.com/jacoblai/arduino-cli/patch/cli/lib"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tr = i18n.Tr

// NewCommand creates a new `outdated` command
func NewCommand() *cobra.Command {
	outdatedCommand := &cobra.Command{
		Use:   "outdated",
		Short: tr("Lists cores and libraries that can be upgraded"),
		Long: tr(`This commands shows a list of installed cores and/or libraries
that can be upgraded. If nothing needs to be updated the output is empty.`),
		Example: "  " + os.Args[0] + " outdated\n",
		Args:    cobra.NoArgs,
		Run:     runOutdatedCommand,
	}
	return outdatedCommand
}

func runOutdatedCommand(cmd *cobra.Command, args []string) {
	inst := instance.CreateAndInit()
	logrus.Info("Executing `arduino-cli outdated`")
	Outdated(inst)
}

// Outdated prints a list of outdated platforms and libraries
func Outdated(inst *rpc.Instance) {
	feedback.PrintResult(
		newOutdatedResult(core.GetList(inst, false, true), lib.GetList(inst, []string{}, false, true)),
	)
}

// output from this command requires special formatting, let's create a dedicated
// feedback.Result implementation
type outdatedResult struct {
	Platforms     []*result.PlatformSummary  `json:"platforms,omitempty"`
	InstalledLibs []*result.InstalledLibrary `json:"libraries,omitempty"`
}

func newOutdatedResult(inPlatforms []*rpc.PlatformSummary, inLibraries []*rpc.InstalledLibrary) *outdatedResult {
	res := &outdatedResult{
		Platforms:     make([]*result.PlatformSummary, len(inPlatforms)),
		InstalledLibs: make([]*result.InstalledLibrary, len(inLibraries)),
	}
	for i, v := range inPlatforms {
		res.Platforms[i] = result.NewPlatformSummary(v)
	}
	for i, v := range inLibraries {
		res.InstalledLibs[i] = result.NewInstalledLibrary(v)
	}
	return res
}

func (ir outdatedResult) Data() interface{} {
	return &ir
}

func (ir outdatedResult) String() string {
	if len(ir.Platforms) == 0 && len(ir.InstalledLibs) == 0 {
		return tr("No outdated platforms or libraries found.")
	}

	// A table useful both for platforms and libraries, where some of the fields will be blank.
	t := table.New()
	t.SetHeader(
		tr("ID"),
		tr("Name"),
		tr("Installed"),
		tr("Latest"),
		tr("Location"),
		tr("Description"),
	)
	t.SetColumnWidthMode(2, table.Average)
	t.SetColumnWidthMode(3, table.Average)
	t.SetColumnWidthMode(5, table.Average)

	// Based on internal/cli/core/list.go
	for _, p := range ir.Platforms {
		name := ""
		if latest := p.GetLatestRelease(); latest != nil {
			name = latest.Name
		}
		if p.Deprecated {
			name = fmt.Sprintf("[%s] %s", tr("DEPRECATED"), name)
		}
		t.AddRow(p.Id, name, p.InstalledVersion, p.LatestVersion, "", "")
	}

	// Based on internal/cli/lib/list.go
	sort.Slice(ir.InstalledLibs, func(i, j int) bool {
		return strings.ToLower(
			ir.InstalledLibs[i].Library.Name,
		) < strings.ToLower(
			ir.InstalledLibs[j].Library.Name,
		) ||
			strings.ToLower(
				ir.InstalledLibs[i].Library.ContainerPlatform,
			) < strings.ToLower(
				ir.InstalledLibs[j].Library.ContainerPlatform,
			)
	})
	lastName := ""
	for _, libMeta := range ir.InstalledLibs {
		lib := libMeta.Library
		name := lib.Name
		if name == lastName {
			name = ` "`
		} else {
			lastName = name
		}

		location := string(lib.Location)
		if lib.ContainerPlatform != "" {
			location = lib.ContainerPlatform
		}

		available := ""
		sentence := ""
		if libMeta.Release != nil {
			available = libMeta.Release.Version
			sentence = lib.Sentence
		}

		if available == "" {
			available = "-"
		}
		if sentence == "" {
			sentence = "-"
		} else if len(sentence) > 40 {
			sentence = sentence[:37] + "..."
		}
		t.AddRow("", name, lib.Version, available, location, sentence)
	}

	return t.Render()
}
