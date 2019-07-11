//
// This file is part of arduino-cli.
//
// Copyright 2018 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to modify or
// otherwise use the software for commercial activities involving the Arduino
// software without disclosing the source code of your own applications. To purchase
// a commercial license, send an email to license@arduino.cc.
//

package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arduino/arduino-cli/cli/errorcodes"
	"github.com/arduino/arduino-cli/cli/globals"
	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/common/formatter"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	colorable "github.com/mattn/go-colorable"
)

// TODO: Feed text output into colorable stdOut
var _ = colorable.NewColorableStdout()

// JSONOrElse outputs the JSON encoding of v if the JSON output format has been
// selected by the user and returns false. Otherwise no output is produced and the
// function returns true.
func JSONOrElse(v interface{}) bool {
	if !globals.OutputJSON {
		return true
	}
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		formatter.PrintError(err, "Error during JSON encoding of the output")
		os.Exit(errorcodes.ErrGeneric)
	}
	fmt.Print(string(d))
	return false
}

// ProgressBar returns a DownloadProgressCB that prints a progress bar.
// If JSON output format has been selected, the callback outputs nothing.
func ProgressBar() commands.DownloadProgressCB {
	if !globals.OutputJSON {
		return NewDownloadProgressBarCB()
	}
	return func(curr *rpc.DownloadProgress) {
		// XXX: Output progress in JSON?
	}
}

// TaskProgress returns a TaskProgressCB that prints the task progress.
// If JSON output format has been selected, the callback outputs nothing.
func TaskProgress() commands.TaskProgressCB {
	if !globals.OutputJSON {
		return NewTaskProgressCB()
	}
	return func(curr *rpc.TaskProgress) {
		// XXX: Output progress in JSON?
	}
}
