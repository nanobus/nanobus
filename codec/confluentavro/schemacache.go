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

package confluentavro

import (
	"fmt"

	"github.com/golang/groupcache/lru"
	"github.com/hamba/avro"
)

type (
	// SchemaCache wraps the calls to the schema registry with a LRU cache.
	SchemaCache struct {
		schemaRegistry schemaRegistry
		schemaCache    *lru.Cache
	}

	// Schema encapsulates an identifier from the schema registry and avro schema.
	Schema struct {
		id       int
		typeName string
		schema   avro.Schema
	}

	// schemaRegistry retrieves Avro schemas from a Confluent Schema Registry by ID.
	schemaRegistry interface {
		GetSchemaByID(id int) (string, error)
	}
)

// ParseSchema parses an Avro schema.
func ParseSchema(id int, schema string) (*Schema, error) {
	avroSchema, err := avro.Parse(schema)
	if err != nil {
		return nil, err
	}
	var typeName string
	if namedSchema, ok := avroSchema.(avro.NamedSchema); ok {
		typeName = namedSchema.FullName()
	}
	return &Schema{
		id:       id,
		typeName: typeName,
		schema:   avroSchema,
	}, nil
}

// ID returns the identifier in the schema registry for this schema.
func (s *Schema) ID() int {
	return s.id
}

// String returns the schema in JSON.
func (s *Schema) String() string {
	return s.schema.String()
}

// NewSchemaCache creates a `SchemaCache`.
func NewSchemaCache(schemaRegistry schemaRegistry, cacheSize int) *SchemaCache {
	cache := lru.New(cacheSize)
	return &SchemaCache{
		schemaRegistry: schemaRegistry,
		schemaCache:    cache,
	}
}

// GetSchema returns the schema for some id.
// The schema registry only provides the schema itself, not the id, subject or version.
func (s *SchemaCache) GetSchema(id int) (*Schema, error) {
	schemaI, ok := s.schemaCache.Get(id)
	var schema *Schema

	if ok {
		schema, ok = schemaI.(*Schema)
	}

	if !ok {
		schemaString, err := s.schemaRegistry.GetSchemaByID(id)
		if err != nil {
			return nil, fmt.Errorf("error getting schema ID %d: %w", id, err)
		}

		schema, err = ParseSchema(id, schemaString)
		if err != nil {
			return nil, fmt.Errorf("error parsing schema ID %d: %w", id, err)
		}
		s.schemaCache.Add(id, schema)
	}

	return schema, nil
}
