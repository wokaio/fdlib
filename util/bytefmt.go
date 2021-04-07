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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	Byte     = 1.0
	Kilobyte = 1024 * Byte
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
	Terabyte = 1024 * Gigabyte
)

var bytesPattern = regexp.MustCompile(`(?i)^(-?\d+(?:\.\d+)?)([KMGT]B?|B)$`)
var invalidByteQuantityError = errors.New("Byte quantity: M, MB, G, or GB")

func ByteSize(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= Terabyte:
		unit = "T"
		value = value / Terabyte
	case bytes >= Gigabyte:
		unit = "G"
		value = value / Gigabyte
	case bytes >= Megabyte:
		unit = "M"
		value = value / Megabyte
	case bytes >= Kilobyte:
		unit = "K"
		value = value / Kilobyte
	case bytes >= Byte:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

func ToMegabytes(s string) (uint64, error) {
	bytes, err := ToBytes(s)
	if err != nil {
		return 0, err
	}

	return bytes / Megabyte, nil
}

func ToBytes(s string) (uint64, error) {
	parts := bytesPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return 0, invalidByteQuantityError
	}

	value, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || value <= 0 {
		return 0, invalidByteQuantityError
	}

	var bytes uint64
	unit := strings.ToUpper(parts[2])
	switch unit[:1] {
	case "T":
		bytes = uint64(value * Terabyte)
	case "G":
		bytes = uint64(value * Gigabyte)
	case "M":
		bytes = uint64(value * Megabyte)
	case "K":
		bytes = uint64(value * Kilobyte)
	case "B":
		bytes = uint64(value * Byte)
	default:
		bytes = uint64(value * Byte)
	}

	return bytes, nil
}
