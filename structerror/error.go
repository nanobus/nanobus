/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package structerror

import (
	"strings"
)

type Error struct {
	code     string
	kvpairs  []kvpair
	metadata map[string]string
}

type kvpair struct {
	key   string
	value string
}

func New(code string, kvpairs ...string) *Error {
	var pairs []kvpair
	var metadata map[string]string
	lenMD := len(kvpairs)

	if lenMD > 1 {
		pairs = make([]kvpair, lenMD/2)
		metadata = make(map[string]string, lenMD/2)

		for i := 0; i < lenMD; i += 2 {
			pairs[i/2] = kvpair{
				key:   kvpairs[i],
				value: kvpairs[i+1],
			}
			metadata[kvpairs[i]] = kvpairs[i+1]
		}
	}

	return &Error{
		code:     code,
		kvpairs:  pairs,
		metadata: metadata,
	}
}

func Parse(contents string) *Error {
	i := strings.IndexRune(contents, '\n')
	if i < 0 {
		return &Error{
			code: contents,
		}
	}

	code := contents[0:i]
	contents = contents[i+1:]
	lines := strings.Split(contents, "\n")
	kvpairs := make([]kvpair, 0, len(lines))
	metadata := make(map[string]string, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(line, "[") {
			continue
		}

		line = strings.TrimPrefix(line, "[")
		i = strings.IndexRune(line, ']')
		k := strings.TrimSpace(line[0:i])
		v := strings.TrimSpace(line[i+1:])
		kvpairs = append(kvpairs, kvpair{
			key:   k,
			value: v,
		})
		metadata[k] = v
	}

	return &Error{
		code:     code,
		kvpairs:  kvpairs,
		metadata: metadata,
	}
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(e.code)

	for _, kv := range e.kvpairs {
		b.WriteString("\n[")
		b.WriteString(kv.key)
		b.WriteString("] ")
		b.WriteString(kv.value)
	}

	return b.String()
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Metadata() map[string]string {
	return e.metadata
}
