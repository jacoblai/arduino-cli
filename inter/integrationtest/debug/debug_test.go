// This file is part of arduino-cli.
//
// Copyright 2022 ARDUINO SA (http://www.arduino.cc/)
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

package debug_test

import (
	"testing"

	"github.com/jacoblai/arduino-cli/inter/integrationtest"
	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	env, cli := integrationtest.CreateArduinoCLIWithEnvironment(t)
	defer env.CleanUp()

	// Init the environment explicitly
	_, _, err := cli.Run("core", "update-index")
	require.NoError(t, err)

	// Install cores
	_, _, err = cli.Run("core", "install", "arduino:samd")
	require.NoError(t, err)

	integrationtest.CLISubtests{
		{"Start", testDebuggerStarts},
		{"WithPdeSketchStarts", testDebuggerWithPdeSketchStarts},
	}.Run(t, env, cli)
}

func testDebuggerStarts(t *testing.T, env *integrationtest.Environment, cli *integrationtest.ArduinoCLI) {
	// Create sketch for testing
	sketchName := "DebuggerStartTest"
	sketchPath := cli.DataDir().Join(sketchName)
	defer sketchPath.RemoveAll()
	fqbn := "arduino:samd:mkr1000"

	_, _, err := cli.Run("sketch", "new", sketchPath.String())
	require.NoError(t, err)

	// Build sketch
	_, _, err = cli.Run("compile", "-b", fqbn, sketchPath.String())
	require.NoError(t, err)

	programmer := "atmel_ice"
	// Starts debugger
	_, _, err = cli.Run("debug", "-b", fqbn, "-P", programmer, sketchPath.String(), "--info")
	require.NoError(t, err)
}

func testDebuggerWithPdeSketchStarts(t *testing.T, env *integrationtest.Environment, cli *integrationtest.ArduinoCLI) {
	sketchName := "DebuggerPdeSketchStartTest"
	sketchPath := cli.DataDir().Join(sketchName)
	defer sketchPath.RemoveAll()
	fqbn := "arduino:samd:mkr1000"

	_, _, err := cli.Run("sketch", "new", sketchPath.String())
	require.NoError(t, err)

	// Looks for sketch file .ino
	pathDir, err := sketchPath.ReadDir()
	require.NoError(t, err)
	fileIno := pathDir[0]

	// Renames sketch file to pde
	filePde := sketchPath.Join(sketchName + ".pde")
	err = fileIno.Rename(filePde)
	require.NoError(t, err)

	// Build sketch
	_, _, err = cli.Run("compile", "-b", fqbn, filePde.String())
	require.NoError(t, err)

	programmer := "atmel_ice"
	// Starts debugger
	_, _, err = cli.Run("debug", "-b", fqbn, "-P", programmer, filePde.String(), "--info")
	require.NoError(t, err)
}
