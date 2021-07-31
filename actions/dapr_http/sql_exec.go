package dapr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/runtime"
)

type SQLExecConfig struct {
	// Name is name of SQL binding to invoke.
	Name string `mapstructure:"name"`
	// SQL
	SQL string `mapstructure:"sql"`
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
}

// SQLExec is the NamedLoader for Dapr output bindings
func SQLExec() (string, actions.Loader) {
	return "@dapr/sql_exec", SQLExecLoader
}

func SQLExecLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c SQLExecConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var httpClient HTTPClient
	var env runtime.Environment
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient,
		"os:env", &env); err != nil {
		return nil, err
	}

	return SQLExecConfigAction(httpClient, env, &c), nil
}

func SQLExecConfigAction(httpClient HTTPClient, env runtime.Environment, config *SQLExecConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		type requestBody struct {
			Operation string            `json:"operation"`
			Data      interface{}       `json:"data,omitempty"`
			Metadata  map[string]string `json:"metadata,omitempty"`
		}

		var err error
		variables := data
		if config.Data != nil {
			var input interface{}
			if input, err = config.Data.Eval(variables); err != nil {
				return nil, err
			}
			variables.Clone()
			variables["input"] = input
		}

		sql, err := replacements(config.SQL, func(exprString string) (string, error) {
			var valueExpr expr.ValueExpr
			if err := valueExpr.DecodeString(exprString); err != nil {
				return "", err
			}
			value, err := valueExpr.Eval(variables)
			if err != nil {
				return "", err
			}
			switch v := value.(type) {
			case string:
				return "'" + strings.ReplaceAll(v, "'", "''") + "'", nil
			default:
				return fmt.Sprintf("%v", v), nil
			}
		})
		if err != nil {
			return nil, err
		}

		r := requestBody{
			Operation: "exec",
			Metadata: map[string]string{
				"sql": sql,
			},
		}

		u, err := url.Parse(daprBaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "v1.0/bindings", config.Name)

		var response interface{}
		err = POST(ctx, httpClient,
			u.String(),
			func() ([]byte, error) {
				return json.Marshal(&r)
			}, func(data []byte) error {
				return coalesce.JSONUnmarshal(data, &response)
			})

		return response, err
	}
}

func replacements(s string, replacer func(string) (string, error)) (string, error) {
	var inVar bool
	var prev byte
	i := 0
	t := len(s)

	var resultBuf bytes.Buffer
	var varBuf bytes.Buffer

	for i < t {
		b := s[i]
		if b != ':' {
			if inVar {
				if b >= 'a' && b <= 'z' ||
					b >= 'A' && b <= 'Z' ||
					b >= '0' && b <= '9' ||
					b == '.' || b == '_' || b == '$' {
					varBuf.WriteByte(b)
				} else {
					pn := varBuf.String()
					if len(pn) > 0 {
						value, err := replacer(pn)
						if err != nil {
							return "", err
						}
						resultBuf.WriteString(value)
					}
					varBuf.Reset()
					inVar = false
					resultBuf.WriteByte(b)
				}
			} else {
				resultBuf.WriteByte(b)
			}
			prev = s[i]
		} else {
			inVar = prev != ':'
			prev = b
		}
		i++
	}

	if inVar {
		pn := varBuf.String()
		if len(pn) > 0 {
			value, err := replacer(pn)
			if err != nil {
				return "", err
			}
			resultBuf.WriteString(value)
		}
	}

	return resultBuf.String(), nil
}
