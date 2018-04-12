/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 ARDUINO AG (http://www.arduino.cc/)
 */

package releases

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IsCached returns a bool representing if the release has already been downloaded
func (r *DownloadResource) IsCached() (bool, error) {
	archivePath, err := r.ArchivePath()
	if err != nil {
		return false, fmt.Errorf("getting archive path: %s", err)
	}

	_, err = os.Stat(archivePath)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("checking archive existence: %s", err)
	}

	return !os.IsNotExist(err), nil
}

// TestLocalArchiveChecksum test if the checksum of the local archive match the checksum of the DownloadResource
func (r *DownloadResource) TestLocalArchiveChecksum() (bool, error) {
	split := strings.SplitN(r.Checksum, ":", 2)
	if len(split) != 2 {
		return false, fmt.Errorf("invalid checksum format: %s", r.Checksum)
	}
	digest, err := hex.DecodeString(split[1])
	if err != nil {
		return false, fmt.Errorf("invalid hash '%s': %s", split[1], err)
	}

	// names based on: https://docs.oracle.com/javase/8/docs/technotes/guides/security/StandardNames.html#MessageDigest
	var algo hash.Hash
	switch split[0] {
	case "SHA-256":
		algo = crypto.SHA256.New()
	case "SHA-1":
		algo = crypto.SHA1.New()
	case "MD5":
		algo = crypto.MD5.New()
	default:
		return false, fmt.Errorf("unsupported hash algorithm: %s", split[0])
	}

	filePath, err := r.ArchivePath()
	if err != nil {
		return false, fmt.Errorf("getting archive path: %s", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("opening archive file: %s", err)
	}
	defer file.Close()
	if _, err := io.Copy(algo, file); err != nil {
		return false, fmt.Errorf("computing hash: %s", err)
	}
	return bytes.Compare(algo.Sum(nil), digest) == 0, nil
}

// TestLocalArchiveSize test if the local archive size match the DownloadResource size
func (r *DownloadResource) TestLocalArchiveSize() (bool, error) {
	filePath, err := r.ArchivePath()
	if err != nil {
		return false, fmt.Errorf("getting archive path: %s", err)
	}
	info, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("getting archive info: %s", err)
	}
	return info.Size() != r.Size, nil
}

// TestLocalArchiveIntegrity checks for integrity of the local archive.
func (r *DownloadResource) TestLocalArchiveIntegrity() (bool, error) {
	if cached, err := r.IsCached(); err != nil {
		return false, fmt.Errorf("testing if archive is cached: %s", err)
	} else if !cached {
		return false, nil
	}

	if ok, err := r.TestLocalArchiveSize(); err != nil {
		return false, fmt.Errorf("teting archive size: %s", err)
	} else if !ok {
		return false, nil
	}

	ok, err := r.TestLocalArchiveChecksum()
	if err != nil {
		return false, fmt.Errorf("testing archive checksum: %s", err)
	}
	return ok, nil
}

const (
	filePermissions = 0644
	packageFileName = "package.json"
)

type packageFile struct {
	Checksum string `json:"checksum"`
}

func computeDirChecksum(root string) (string, error) {
	hash := sha256.New()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || (info.Name() == packageFileName && filepath.Dir(path) == root) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		if _, err := io.Copy(hash, f); err != nil {
			return fmt.Errorf("failed to compute hash of file \"%s\"", info.Name())
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func createPackageFile(root string) error {
	checksum, err := computeDirChecksum(root)
	if err != nil {
		return err
	}

	packageJSON, _ := json.Marshal(packageFile{checksum})
	err = ioutil.WriteFile(filepath.Join(root, packageFileName), packageJSON, filePermissions)
	if err != nil {
		return err
	}
	return nil
}

// CheckDirChecksum reads checksum from the package.json and compares it with a recomputed value.
func CheckDirChecksum(root string) (bool, error) {
	packageJSON, err := ioutil.ReadFile(filepath.Join(root, packageFileName))
	if err != nil {
		return false, err
	}
	var file packageFile
	json.Unmarshal(packageJSON, &file)
	checksum, err := computeDirChecksum(root)
	if err != nil {
		return false, err
	}
	return file.Checksum == checksum, nil
}
