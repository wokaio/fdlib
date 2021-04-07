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
	"sync"
)

// JSONDataStream interface
type JSONDataStream struct {
	buffer []byte
	comma  bool
}

// Array JSONDataStream interface
type Array JSONDataStream

// Object JSONDataStream interface
type Object JSONDataStream

func (p *JSONDataStream) writeArray(b []byte) {
	p.buffer = append(p.buffer, b...)
}

func (p *JSONDataStream) write(b byte) {
	p.buffer = append(p.buffer, b)
}

func (p *JSONDataStream) reset() {
	p.buffer = p.buffer[:0]
	p.comma = false
}

// NewJSONDataStream JSON object
func NewJSONDataStream() *JSONDataStream {
	js := &JSONDataStream{buffer: make([]byte, 0, 1024*1024)}
	js.reset()
	return js
}

func (p *JSONDataStream) end() []byte {
	return p.buffer
}

func (p *JSONDataStream) consumeComma() {
	if p.comma {
		p.write(',')
		p.comma = false
	}
}

func (p *JSONDataStream) putComma() {
	if p.comma {
		p.write(',')
	}
}

// OpenArray to JSON object
func (p *JSONDataStream) OpenArray() {
	p.consumeComma()
	p.write('[')
}

// OpenObject to JSON object
func (p *JSONDataStream) OpenObject() {
	p.consumeComma()
	p.write('{')
}

// CloseArray to JSON object
func (p *JSONDataStream) CloseArray() {
	p.write(']')
	p.comma = true
}

// CloseObject to JSON object
func (p *JSONDataStream) CloseObject() {
	p.write('}')
	p.comma = true
}

// PutKey to JSON object
func (p *JSONDataStream) PutKey(key []byte) {
	p.consumeComma()
	p.write('"')
	p.escapedCopy(key)
	p.write('"')
	p.write(':')
}

// PutInt to JSON object
func (p *JSONDataStream) PutInt(value int) {
	p.putComma()
	p.comma = true
	p.writeArray([]byte(fmt.Sprintf("%d", value)))
}

// PutFloat64 to JSON object
func (p *JSONDataStream) PutFloat64(value float64) {
	p.putComma()
	p.comma = true
	p.writeArray([]byte(fmt.Sprintf("%f", value)))
}

// PutNull to JSON object
func (p *JSONDataStream) PutNull() {
	p.putComma()
	p.comma = true
	p.writeArray([]byte("null"))
}

// PutBoolean to JSON object
func (p *JSONDataStream) PutBoolean(value bool) {
	p.putComma()
	p.comma = true
	if value {
		p.writeArray([]byte("true"))
	} else {
		p.writeArray([]byte("false"))
	}
}

func (p *JSONDataStream) escapedCopy(value []byte) {
	for i := 0; i < len(value); i++ {
		if value[i] != '\\' && (value[i] > '"' || value[i] == ' ') {
			p.write(value[i])
		} else if value[i] == '"' {
			p.write('\\')
			p.write('"')
		} else if value[i] == '\\' {
			p.write('\\')
			p.write('"')
		} else if value[i] == '\n' {
			p.write('\\')
			p.write('n')
		} else if value[i] == '\r' {
			p.write('\\')
			p.write('r')
		} else if value[i] == '\t' {
			p.write('\\')
			p.write('t')
		} else if value[i] == '\f' {
			p.write('\\')
			p.write('f')
		} else if value[i] == '\b' {
			p.write('\\')
			p.write('b')
		} else {
			p.write(value[i])
		}
	}
}

// PutString to JSON object
func (p *JSONDataStream) PutString(value []byte) {
	p.putComma()
	p.comma = true
	p.write('"')
	p.escapedCopy(value)
	p.write('"')
}

func (p *JSONDataStream) routeValueType(value interface{}) {
	switch v := value.(type) {
	case string:
		p.PutString([]byte(v))
	case int:
		p.PutInt(v)
	case float64:
		p.PutFloat64(v)
	case bool:
		p.PutBoolean(v)
	case func(array *Array):
		p.OpenArray()
		v((*Array)(p))
		p.CloseArray()
	case func(object *Object):
		p.OpenObject()
		v((*Object)(p))
		p.CloseObject()
	default:
		p.PutNull()
	}
}

// Put to JSON object
func (p *Array) Put(value interface{}) {
	(*JSONDataStream)(p).routeValueType(value)
}

// Put to JSON object
func (p *Object) Put(key string, value interface{}) {
	(*JSONDataStream)(p).PutKey([]byte(key))
	(*JSONDataStream)(p).routeValueType(value)
}

var streamPool = sync.Pool{
	New: func() interface{} {
		return NewJSONDataStream()
	},
}

// Marshal JSON to interface
func Marshal(value interface{}) []byte {
	var js *JSONDataStream = streamPool.Get().(*JSONDataStream)
	js.reset()
	js.routeValueType(value)

	var ret []byte = js.end()
	streamPool.Put(js)
	return ret
}
