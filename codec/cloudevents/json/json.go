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

package json

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type (
	Config struct {
		SpecVersion string `mapstructure:"specversion"`
		Source      string `mapstructre:"source"`
	}

	// Codec encodes and decodes Avro records.
	Codec struct {
		config *Config
	}
)

// CloudEventsJSON is the NamedLoader for this codec.
func CloudEventsJSON() (string, bool, codec.Loader) {
	return "cloudevents+json", true, Loader
}

func Loader(with interface{}, resolver resolve.ResolveAs) (codec.Codec, error) {
	c := Config{
		SpecVersion: "1.0",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return NewCodec(&c), nil
}

// NewCodec creates a `Codec`.
func NewCodec(c *Config) *Codec {
	return &Codec{
		config: c,
	}
}

func (c *Codec) ContentType() string {
	return "application/cloudevents+json"
}

// Decode decodes JSON bytes to a value.
func (c *Codec) Decode(msgValue []byte, args ...interface{}) (interface{}, string, error) {
	var data map[string]interface{}
	if err := coalesce.JSONUnmarshal(msgValue, &data); err != nil {
		return nil, "", err
	}

	var typeValue string
	if typeField, ok := data["type"]; ok {
		typeValue, _ = typeField.(string)
	}

	return data, typeValue, nil
}

// Encode encodes a value into JSON encoded bytes.
func (c *Codec) Encode(value interface{}, args ...interface{}) ([]byte, error) {
	if m, ok := value.(map[string]interface{}); ok {
		if c.config.SpecVersion != "" {
			if _, exists := m["specversion"]; !exists {
				m["specversion"] = c.config.SpecVersion
			}
		}
		if _, exists := m["id"]; !exists {
			m["id"] = uuid.New().String()
		}
		if c.config.Source != "" {
			if _, exists := m["source"]; !exists {
				m["source"] = c.config.Source
			}
		}
		if _, exists := m["datacontenttype"]; !exists {
			data := m["data"]
			if _, ok := data.([]byte); ok {
				m["datacontenttype"] = "application/octet-stream"
			} else {
				m["datacontenttype"] = "application/json"
			}
		}
		if _, exists := m["time"]; !exists {
			m["time"] = time.Now().UTC().Format(time.RFC3339Nano)
		}
	}

	return json.Marshal(value)
}
