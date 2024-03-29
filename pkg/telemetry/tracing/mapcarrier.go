/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/propagation"
)

// MapCarrier is a TextMapCarrier that uses a map held in memory as a storage
// medium for propagated key-value pairs.
type MapCarrier map[string]interface{}

// Compile time check that MapCarrier implements the TextMapCarrier.
var _ propagation.TextMapCarrier = MapCarrier{}

// Get returns the value associated with the passed key.
func (c MapCarrier) Get(key string) string {
	return fmt.Sprintf("%s", c[key])
}

// Set stores the key-value pair.
func (c MapCarrier) Set(key, value string) {
	c[key] = value
}

// Keys lists the keys stored in this carrier.
func (c MapCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
