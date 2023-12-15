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

package sketch

import (
	"context"

	paths "github.com/arduino/go-paths-helper"
	"github.com/jacoblai/arduino-cli/arduino"
	"github.com/jacoblai/arduino-cli/arduino/sketch"
	rpc "github.com/jacoblai/arduino-cli/rpc/cc/arduino/cli/commands/v1"
)

// SetSketchDefaults updates the sketch project file (sketch.yaml) with the given defaults
// for the values `default_fqbn`, `default_port`, and `default_protocol`.
func SetSketchDefaults(ctx context.Context, req *rpc.SetSketchDefaultsRequest) (*rpc.SetSketchDefaultsResponse, error) {
	sk, err := sketch.New(paths.New(req.SketchPath))
	if err != nil {
		return nil, &arduino.CantOpenSketchError{Cause: err}
	}

	oldAddress, oldProtocol := sk.GetDefaultPortAddressAndProtocol()
	res := &rpc.SetSketchDefaultsResponse{
		DefaultFqbn:         sk.GetDefaultFQBN(),
		DefaultPortAddress:  oldAddress,
		DefaultPortProtocol: oldProtocol,
	}

	if fqbn := req.GetDefaultFqbn(); fqbn != "" {
		if err := sk.SetDefaultFQBN(fqbn); err != nil {
			return nil, &arduino.CantUpdateSketchError{Cause: err}
		}
		res.DefaultFqbn = fqbn
	}
	if newAddress, newProtocol := req.GetDefaultPortAddress(), req.GetDefaultPortProtocol(); newAddress != "" {
		if err := sk.SetDefaultPort(newAddress, newProtocol); err != nil {
			return nil, &arduino.CantUpdateSketchError{Cause: err}
		}
		res.DefaultPortAddress = newAddress
		res.DefaultPortProtocol = newProtocol
	}

	return res, nil
}
