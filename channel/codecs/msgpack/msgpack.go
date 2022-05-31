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

package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"
)

type Codec struct{}

func New() *Codec {
	return &Codec{}
}

func (c *Codec) ContentType() string {
	return "application/msgpack"
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (c *Codec) Decode(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
