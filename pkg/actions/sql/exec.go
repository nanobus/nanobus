/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sql

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/resource"
)

func ExecLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := ExecConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var resources resource.Resources
	if err := resolve.Resolve(resolver,
		"resource:lookup", &resources); err != nil {
		return nil, err
	}

	db, err := resource.Get[*sqlx.DB](resources, c.Resource)
	if err != nil {
		return nil, err
	}

	return ExecAction(db, &c), nil
}

func ExecAction(
	db *sqlx.DB,
	config *ExecConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var err error
		var input interface{} = map[string]interface{}(data)
		if config.Data != nil {
			input, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		if multi, ok := input.([]interface{}); ok {
			tx, err := db.Begin()
			if err != nil {
				return nil, err
			}
			defer tx.Rollback() // nolint:errcheck

			for _, item := range multi {
				if single, ok := item.(map[string]interface{}); ok {
					args := make([]interface{}, len(config.Args))
					for i, expr := range config.Args {
						var err error
						if args[i], err = expr.Eval(single); err != nil {
							return nil, err
						}
					}

					_, err := tx.Exec(config.SQL, args...)
					if err != nil {
						return nil, err
					}
					// if tag.RowsAffected() == 0 {
					// 	return errors.New("no rows effected")
					// }
				}
			}
		} else if single, ok := input.(map[string]interface{}); ok {
			args := make([]interface{}, len(config.Args))
			for i, expr := range config.Args {
				var err error
				if args[i], err = expr.Eval(single); err != nil {
					return nil, err
				}
			}

			_, err := db.Exec(config.SQL, args...)
			if err != nil {
				return nil, err
			}
			// if tag.RowsAffected() == 0 {
			// 	return nil, errors.New("no rows effected")
			// }
		}

		return nil, nil
	}
}
