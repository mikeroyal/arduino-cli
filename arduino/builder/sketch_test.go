// This file is part of arduino-cli.
//
// Copyright 2019 ARDUINO SA (http://www.arduino.cc/)
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

package builder_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/arduino/arduino-cli/arduino/builder"
	"github.com/stretchr/testify/assert"
)

func TestSaveSketch(t *testing.T) {
	sketchName := t.Name() + ".ino"
	outName := sketchName + ".cpp"
	sketchFile := filepath.Join("testdata", sketchName)
	tmp := tmpDirOrDie()
	defer os.RemoveAll(tmp)
	source, err := ioutil.ReadFile(sketchFile)
	if err != nil {
		t.Fatalf("unable to read golden file %s: %v", sketchFile, err)
	}

	builder.SaveSketch(sketchName, string(source), tmp)

	out, err := ioutil.ReadFile(filepath.Join(tmp, outName))
	if err != nil {
		t.Fatalf("unable to read output file %s: %v", outName, err)
	}

	assert.Equal(t, source, out)
}