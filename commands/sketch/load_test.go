// This file is part of arduino-cli.
//
// Copyright 2023 ARDUINO SA (http://www.arduino.cc/)
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

package sketch

import (
	"context"
	"testing"

	"github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/stretchr/testify/require"
)

func TestLoadSketchProfiles(t *testing.T) {
	loadResp, err := LoadSketch(context.Background(), &commands.LoadSketchRequest{
		SketchPath: "./testdata/sketch_with_profile",
	})
	require.NoError(t, err)
	require.Len(t, loadResp.GetProfiles(), 2)
	require.Equal(t, loadResp.GetDefaultProfile().GetName(), "nanorp")
}
