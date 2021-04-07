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
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"

	uuid_ "github.com/linuxpham/go.uuid"
	"github.com/speps/go-hashids"
	h3 "github.com/uber/h3-go/v3"
)

func String(msg string) *string {
	return &msg
}

func GetH3GeoString(latNum float64, lonNum float64, resolutionNum int) string {
	geo := h3.GeoCoord{
		Latitude:  latNum,
		Longitude: lonNum,
	}
	return fmt.Sprintf("%#x", h3.FromGeo(geo, resolutionNum))
}

func StringToFloat64(numStr string) (float64, error) {
	return strconv.ParseFloat(numStr, 64)
}

func StringToInt64(numStr string) (int64, error) {
	return strconv.ParseInt(numStr, 10, 64)
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetUUID() string {
	id, _ := uuid_.NewV4()
	return fmt.Sprintf("%s", id.String())
}

func GetUnixTimestamp() int32 {
	return int32(time.Now().Unix())
}

func ByteArrayToInt64(buf []byte) int64 {
	x, _ := binary.Varint(buf)
	return x
}

func Int64ToByteArray(num int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, num)
	b := buf[:n]
	return b
}

func Int64ToEncodeString(salt string, num int64) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{num})
	return e
}

func StringToDecodeInt64(salt string, e string) int64 {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 30
	h, err := hashids.NewWithData(hd)
	if err == nil {
		return 0
	}

	d, err := h.DecodeInt64WithError(e)
	if err == nil {
		return 0
	}

	return d[0]
}

func Base64String(inBuffer *bytes.Buffer) string {
	return base64.URLEncoding.EncodeToString(inBuffer.Bytes())
}

func Base64StringToHex(s string) (string, error) {
	de, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", de)
	return hex, nil
}

func Deduplication(slice *[]string) {
	found := make(map[string]bool)
	total := 0
	for k, v := range *slice {
		if _, ok := found[v]; !ok {
			found[v] = true
			(*slice)[total] = (*slice)[k]
			total++
		}
	}
	*slice = (*slice)[:total]
}
