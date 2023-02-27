/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/resource"
)

type ConnectionConfig struct {
	Address  string `mapstructure:"address" validate:"required"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Connection is the NamedLoader for a redis connection.
func Connection() (string, resource.Loader) {
	return "redis", ConnectionLoader
}

func ConnectionLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (interface{}, error) {
	var c ConnectionConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password,
		DB:       c.DB,
	})

	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)

	return client, err
}
