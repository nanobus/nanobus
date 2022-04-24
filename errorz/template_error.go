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

package errorz

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

type TemplateError struct {
	Template string
	Metadata Metadata
}

func Return(template string, metadata ...Metadata) *TemplateError {
	var md Metadata
	l := len(metadata)

	if l == 1 {
		md = metadata[0]
	} else if l > 1 {
		md = make(Metadata)
		for _, m := range metadata {
			for k, v := range m {
				md[k] = v
			}
		}
	}

	return &TemplateError{
		Template: template,
		Metadata: md,
	}
}

func ParseTemplateError(message string) TemplateError {
	scanner := bufio.NewScanner(strings.NewReader(message))
	template := "unknown"
	var md Metadata

	if scanner.Scan() {
		template = scanner.Text()
	}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "[") {
			continue
		}
		line = line[1:]
		parts := strings.SplitN(line, "] ", 2)
		if len(parts) == 2 {
			if md == nil {
				md = make(Metadata)
			}
			md[parts[0]] = parts[1]
		}
	}

	return TemplateError{
		Template: template,
		Metadata: md,
	}
}

func (e *TemplateError) Error() string {
	var sb strings.Builder

	keys := make([]string, len(e.Metadata))
	i := 0
	for k := range e.Metadata {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	sb.WriteString(e.Template)
	for _, k := range keys {
		v := e.Metadata[k]
		sb.WriteString("\n[")
		sb.WriteString(k)
		sb.WriteString("] ")
		sb.WriteString(fmt.Sprintf("%v", v))
	}

	return sb.String()
}
