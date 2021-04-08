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

type CounterNumber uint64

func (c *CounterNumber) Inc(delta uint) {
	if c == nil {
		return
	}
	atomic.AddUint64((*uint64)(c), uint64(delta))
}

func (c *CounterNumber) Get() int64 {
	if c == nil {
		return 0
	}
	return int64(atomic.LoadUint64((*uint64)(c)))
}

func (c *CounterNumber) Type() string {
	return TypeCounter
}
