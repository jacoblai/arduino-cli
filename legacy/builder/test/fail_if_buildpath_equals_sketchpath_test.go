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

package test

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/jacoblai/arduino-cli/arduino/sketch"
	"github.com/jacoblai/arduino-cli/legacy/builder"
	"github.com/jacoblai/arduino-cli/legacy/builder/types"
	"github.com/stretchr/testify/require"
)

func TestFailIfBuildPathEqualsSketchPath(t *testing.T) {
	ctx := &types.Context{
		Sketch:    &sketch.Sketch{FullPath: paths.New("buildPath")},
		BuildPath: paths.New("buildPath"),
	}

	command := builder.FailIfBuildPathEqualsSketchPath{}
	require.Error(t, command.Run(ctx))
}

func TestFailIfBuildPathEqualsSketchPathSketchPathDiffers(t *testing.T) {
	ctx := &types.Context{
		Sketch:    &sketch.Sketch{FullPath: paths.New("sketchPath")},
		BuildPath: paths.New("buildPath"),
	}

	command := builder.FailIfBuildPathEqualsSketchPath{}
	NoError(t, command.Run(ctx))
}
