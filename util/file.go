// Copyright (c) 2021 Miczone Asia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"path/filepath"
)

func CreateFile(fileName string) error {
	fd, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = fd.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreateFileWithSize(fileName string, size int64) error {
	fd, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Seek(size-1, 0)
	if err != nil {
		return err
	}

	_, err = fd.Write([]byte{0})
	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(fileName string) error {
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}
	return nil
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}
