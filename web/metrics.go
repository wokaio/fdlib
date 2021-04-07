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

package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	DIntervalNumber = 15
	KindTotal       = "total"
	KindDelta       = "delta"
	TypeGauge       = "GaugeNumber"
	TypeCounter     = "CounterNumber"
	TypeState       = "StateNumber"
)

var (
	errStructPtrType   = errors.New("Metrics should be struct pointor")
	errStructFieldType = errors.New("Struct field shoule be *CounterNumber or *GaugeNumber or *StateNumber")
)

var (
	supportTypes = map[string]bool{TypeGauge: true, TypeCounter: true, TypeState: true}
)

type MetricStats struct {
	metricStruct interface{}
	metricPrefix string
	interval     int
	counterMap   map[string]*CounterNumber
	gaugeMap     map[string]*GaugeNumber
	stateMap     map[string]*StateNumber

	lock        sync.RWMutex
	metricsLast *MetricsData
	metricsDiff *MetricsData
}

// NewMetricStats returns a new, empty MetricStats
func NewMetricStats(metrics interface{}, prefix string, interval int) (*MetricStats, error) {
	m := new(MetricStats)
	if err := validateMetrics(metrics); err != nil {
		return m, err
	}

	if interval <= 0 {
		interval = DIntervalNumber
	}

	m.metricStruct = metrics
	m.metricPrefix = prefix
	m.interval = interval
	m.initMetrics(metrics)

	m.metricsLast = m.GetAll()
	m.metricsDiff = m.metricsLast.Diff(m.metricsLast)

	go m.handleCounterDiff(m.interval)
	return m, nil
}

func validateMetrics(metrics interface{}) error {
	// check type of counters is pointer to struct
	t := reflect.TypeOf(metrics)
	if t.Kind() != reflect.Ptr {
		return errStructPtrType
	}

	s := t.Elem()
	if s.Kind() != reflect.Struct {
		return errStructPtrType
	}

	// check type of struct field is *Counter || *Gauge
	for i := 0; i < s.NumField(); i++ {
		ft := s.Field(i).Type
		if ft.Kind() != reflect.Ptr {
			return errStructFieldType
		}

		fn := ft.Elem().Name()
		if _, ok := supportTypes[fn]; !ok {
			return errStructFieldType
		}
	}

	return nil
}

// GetAll gets absoulute values for all counters
func (m *MetricStats) GetAll() *MetricsData {
	d := NewMetricsData(m.metricPrefix, KindTotal)
	for k, c := range m.counterMap {
		d.CounterData[k] = int64(c.Get())
	}

	for k, c := range m.gaugeMap {
		d.GaugeData[k] = int64(c.Get())
	}

	for k, s := range m.stateMap {
		d.StateData[k] = s.Get()
	}

	return d
}

// GetDiff gets diff values for all counters
func (m *MetricStats) GetDiff() *MetricsData {
	m.lock.RLock()
	diff := m.metricsDiff
	m.lock.RUnlock()
	return diff
}

// initMetrics initializes metrics struct
func (m *MetricStats) initMetrics(s interface{}) {
	m.counterMap = make(map[string]*CounterNumber)
	m.gaugeMap = make(map[string]*GaugeNumber)
	m.stateMap = make(map[string]*StateNumber)

	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		name := m.convert(field.Name)

		switch mType := field.Type.Elem().Name(); mType {
		case TypeState:
			v := new(StateNumber)
			m.stateMap[name] = v
			value.Set(reflect.ValueOf(v))

		case TypeCounter:
			v := new(CounterNumber)
			m.counterMap[name] = v
			value.Set(reflect.ValueOf(v))

		case TypeGauge:
			v := new(GaugeNumber)
			m.gaugeMap[name] = v
			value.Set(reflect.ValueOf(v))
		}
	}
}

// convert converts name from CamelCase to UnderScoreCase
func (m *MetricStats) convert(name string) string {
	var b bytes.Buffer
	for i, c := range name {
		if unicode.IsUpper(c) {
			if i > 0 {
				b.WriteString("_")
			}
			b.WriteRune(c)
		} else {
			b.WriteRune(unicode.ToUpper(c))
		}
	}
	return b.String()
}

// handleCounterDiff is go-routine for periodically update counter diff
func (m *MetricStats) handleCounterDiff(interval int) {
	for {
		m.updateDiff()

		seconds := time.Now().Second()
		left := m.interval - seconds%m.interval
		time.Sleep(time.Duration(left) * time.Second)
	}
}

// updateDiff updates diff values for all counters
func (m *MetricStats) updateDiff() {
	var diff *MetricsData

	m.lock.RLock()
	last := m.metricsLast
	m.lock.RUnlock()

	current := m.GetAll()
	diff = current.Diff(last)

	m.lock.Lock()
	m.metricsLast = current
	m.metricsDiff = diff
	m.lock.Unlock()
}

type MetricsData struct {
	Prefix      string
	Kind        string
	GaugeData   map[string]int64
	CounterData map[string]int64
	StateData   map[string]string
}

func NewMetricsData(prefix string, kind string) *MetricsData {
	d := new(MetricsData)
	d.Prefix = prefix
	d.Kind = kind
	d.GaugeData = make(map[string]int64)
	d.CounterData = make(map[string]int64)
	d.StateData = make(map[string]string)
	return d
}

func (d *MetricsData) Diff(last *MetricsData) *MetricsData {
	diff := NewMetricsData(d.Prefix+"_diff", KindDelta)

	for k, v := range d.CounterData {

		if v2, ok := last.CounterData[k]; ok {
			diff.CounterData[k] = v - v2
		} else {
			diff.CounterData[k] = v
		}

	}
	return diff
}

func (d *MetricsData) Sum(d2 *MetricsData) *MetricsData {
	for k, v := range d2.CounterData {
		if v0, ok := d.CounterData[k]; ok {
			d.CounterData[k] = v0 + v
		} else {
			d.CounterData[k] = v
		}
	}
	return d
}

func (d *MetricsData) KeyValueFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.CounterData {
		line := fmt.Sprintf("%s_%s: %d\n", d.Prefix, k, v)
		b.WriteString(line)
	}

	for k, v := range d.GaugeData {
		line := fmt.Sprintf("%s_%s: %d\n", d.Prefix, k, v)
		b.WriteString(line)
	}

	for k, v := range d.StateData {
		line := fmt.Sprintf("%s_%s: %s\n", d.Prefix, k, v)
		b.WriteString(line)
	}
	return b.Bytes()
}

func (d *MetricsData) PrometheusFormat() []byte {
	var b bytes.Buffer
	for k, v := range d.CounterData {
		key := fmt.Sprintf("%s_%s", d.Prefix, k)
		line := fmt.Sprintf("# TYPE %s %s\n%s %d\n", key, strings.ToLower(TypeCounter), key, v)
		b.WriteString(line)
	}

	for k, v := range d.GaugeData {
		key := fmt.Sprintf("%s_%s", d.Prefix, k)
		line := fmt.Sprintf("# TYPE %s %s\n%s %d\n", key, strings.ToLower(TypeGauge), key, v)
		b.WriteString(line)
	}
	return b.Bytes()
}

func (d *MetricsData) Format(params map[string][]string) ([]byte, error) {
	format, err := GetParamValue(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json":
		return json.Marshal(d)
	case "kv", "noah":
		return d.KeyValueFormat(), nil
	case "prometheus":
		return d.PrometheusFormat(), nil
	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}
}
