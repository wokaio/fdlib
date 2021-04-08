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

package metric

import "sync/atomic"

type GaugeNumber int64

func (c *GaugeNumber) Inc(delta uint) {
	if c == nil {
		return
	}
	atomic.AddInt64((*int64)(c), int64(delta))
}

func (c *GaugeNumber) Dec(delta uint) {
	if c == nil {
		return
	}
	atomic.AddInt64((*int64)(c), int64(-delta))
}

func (c *GaugeNumber) Get() int64 {
	if c == nil {
		return 0
	}
	return atomic.LoadInt64((*int64)(c))
}

func (c *GaugeNumber) Set(v int64) {
	if c == nil {
		return
	}
	atomic.StoreInt64((*int64)(c), v)
}

func (c *GaugeNumber) Type() string {
	return TypeGauge
}
