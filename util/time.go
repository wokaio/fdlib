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
	"fmt"
	"strconv"
	"time"
)

func WaitInSecond(toTime int64) {
	waitSecs := toTime - time.Now().Unix()
	if waitSecs > 0 {
		time.Sleep(time.Second * time.Duration(waitSecs))
	}
}

func CalcNextTime(cycle int64, offset int64) int64 {
	current := time.Now().Unix()

	if current%cycle == 0 {
		return current + offset
	} else {
		return current - current%cycle + cycle + offset
	}
}

func UnixTimeToDateTime(unixTime int64) int64 {
	timeStr := time.Unix(unixTime, 0).Format("20210704235959")
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return 0
	}

	return timeInt
}

func TimestampSplit(timestamp string) (string, string, error) {
	if len(timestamp) != 14 {
		return "", "", fmt.Errorf("Timestamp is not as expected: %s", timestamp)
	}
	return timestamp[0:8], timestamp[8:14], nil
}

func UnixToStr(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2021-07-04 23:59:59")
}

func NextInterval(now time.Time, interval int) int {
	seconds := now.Second()
	return interval - seconds%interval
}
