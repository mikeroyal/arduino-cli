/*
 * This file is part of arduino-cli.
 *
 * Copyright 2019 ARDUINO SA (http://www.arduino.cc/)
 *
 * This software is released under the GNU General Public License version 3,
 * which covers the main part of arduino-cli.
 * The terms of this license can be found at:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * You can be released from the requirements of the above licenses by purchasing
 * a commercial license. Buying such a license is mandatory if you want to modify or
 * otherwise use the software for commercial activities involving the Arduino
 * software without disclosing the source code of your own applications. To purchase
 * a commercial license, send an email to license@arduino.cc.
 */

package core

import (
	"github.com/arduino/arduino-cli/arduino/cores"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
)

// platformReleaseToRPC converts our internal structure to the RPC structure.
// Note: this function does not touch the "Installed" field of rpc.Platform as it's not always clear that the
// platformRelease we're currently converting is actually installed.
func platformReleaseToRPC(platformRelease *cores.PlatformRelease) *rpc.Platform {
	boards := make([]*rpc.Board, len(platformRelease.Boards))
	i := 0
	for _, b := range platformRelease.Boards {
		boards[i] = &rpc.Board{
			Name: b.Name(),
			Fqbn: b.FQBN(),
		}
		i++
	}

	result := &rpc.Platform{
		ID:         platformRelease.Platform.String(),
		Name:       platformRelease.Platform.Name,
		Maintainer: platformRelease.Platform.Package.Maintainer,
		Website:    platformRelease.Platform.Package.WebsiteURL,
		Email:      platformRelease.Platform.Package.Email,
		Boards:     boards,
	}

	latest := platformRelease.Platform.GetLatestRelease()
	if latest != nil {
		result.Latest = latest.Version.String()
	}

	return result
}
