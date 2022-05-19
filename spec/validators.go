package spec

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

var validators = map[string]ValidationLoader{
	"url": func(t *TypeRef, f *Field, a *Annotation) (Validation, error) {
		return func(v interface{}) ([]ValidationError, error) {
			val := validator.New()
			value := cast.ToString(v)

			if err := val.Var(value, "url"); err != nil {
				return []ValidationError{
					{
						Fields:  []string{f.Name},
						Message: fmt.Sprintf("%q is an invalid URL", f.Name),
					},
				}, nil
			}

			return nil, nil
		}, nil
	},
}
