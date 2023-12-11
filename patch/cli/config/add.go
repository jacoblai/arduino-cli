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

package config

import (
	"os"
	"reflect"

	"github.com/jacoblai/arduino-cli/configuration"
	"github.com/jacoblai/arduino-cli/patch/cli/feedback"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func uniquify[T comparable](s []T) []T {
	// use a map, which enforces unique keys
	inResult := make(map[T]bool)
	var result []T
	// loop through input slice **in order**,
	// to ensure output retains that order
	// (except that it removes duplicates)
	for i := 0; i < len(s); i++ {
		// attempt to use the element as a key
		if _, ok := inResult[s[i]]; !ok {
			// if key didn't exist in map,
			// add to map and append to result
			inResult[s[i]] = true
			result = append(result, s[i])
		}
	}
	return result
}

func initAddCommand() *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add",
		Short: tr("Adds one or more values to a setting."),
		Long:  tr("Adds one or more values to a setting."),
		Example: "" +
			"  " + os.Args[0] + " config add board_manager.additional_urls https://example.com/package_example_index.json\n" +
			"  " + os.Args[0] + " config add board_manager.additional_urls https://example.com/package_example_index.json https://another-url.com/package_another_index.json\n",
		Args: cobra.MinimumNArgs(2),
		Run:  runAddCommand,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return GetConfigurationKeys(), cobra.ShellCompDirectiveDefault
		},
	}
	return addCommand
}

func runAddCommand(cmd *cobra.Command, args []string) {
	logrus.Info("Executing `arduino-cli config add`")
	key := args[0]
	kind := validateKey(key)

	if kind != reflect.Slice {
		msg := tr("The key '%[1]v' is not a list of items, can't add to it.\nMaybe use '%[2]s'?", key, "config set")
		feedback.Fatal(msg, feedback.ErrGeneric)
	}

	v := configuration.Settings.GetStringSlice(key)
	v = append(v, args[1:]...)
	v = uniquify(v)
	configuration.Settings.Set(key, v)

	if err := configuration.Settings.WriteConfig(); err != nil {
		feedback.Fatal(tr("Can't write config file: %v", err), feedback.ErrGeneric)
	}
}
