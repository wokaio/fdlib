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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PrettyDuration time.Duration
type PrettyAge time.Time

var prettyDurationRe = regexp.MustCompile(`\.[0-9]+`)
var ageUnits = []struct {
	Size   time.Duration
	Symbol string
}{
	{12 * 30 * 24 * time.Hour, "y"},
	{30 * 24 * time.Hour, "mo"},
	{7 * 24 * time.Hour, "w"},
	{24 * time.Hour, "d"},
	{time.Hour, "h"},
	{time.Minute, "m"},
	{time.Second, "s"},
}

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

func (d PrettyDuration) String() string {
	label := fmt.Sprintf("%v", time.Duration(d))
	if match := prettyDurationRe.FindString(label); len(match) > 4 {
		label = strings.Replace(label, match, match[:4], 1)
	}
	return label
}

func (t PrettyAge) String() string {
	// Calculate the time difference and handle the 0 cornercase
	diff := time.Since(time.Time(t))
	if diff < time.Second {
		return "0"
	}
	// Accumulate a precision of 3 components before returning
	result, prec := "", 0

	for _, unit := range ageUnits {
		if diff > unit.Size {
			result = fmt.Sprintf("%s%d%s", result, diff/unit.Size, unit.Symbol)
			diff %= unit.Size

			if prec += 1; prec >= 3 {
				break
			}
		}
	}
	return result
}
