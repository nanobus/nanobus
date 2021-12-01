package errorz

import (
	"fmt"
	"time"
)

type Error struct {
	// Type is a textual type for the error.
	Type string `json:"type,omitempty" yaml:"type,omitempty" msgpack:"type,omitempty"`
	// Code is the numeric error code to return.
	Code ErrCode `json:"code" yaml:"code" msgpack:"code"`
	// Is the transport specific status code for the
	// error type/code.
	Status int `json:"status,omitempty" yaml:"status,omitempty" msgpack:"status,omitempty"`
	// Title is a short message of the error.
	Title string `json:"title,omitempty" yaml:"title,omitempty" msgpack:"title,omitempty"`
	// Message is a descriptive message of the error.
	Message string `json:"message,omitempty" yaml:"message,omitempty" msgpack:"message,omitempty"`
	// Details are user-defined additional details.
	Details interface{} `json:"details,omitempty" yaml:"details,omitempty" msgpack:"details,omitempty"`
	// Metadata is debugging information structured as key-value
	// pairs. Metadata is not exposed to external clients.
	Metadata Metadata `json:"-" yaml:"-" msgpack:"-"`
	// Instance is a URI that identifies the specific
	// occurrence of the error.
	Instance string `json:"instance,omitempty" yaml:"instance,omitempty" msgpack:"instance,omitempty"`
	// Err is the underlying error if any.
	Err error `json:"-" yaml:"-" msgpack:"-"`
	// Errors encapsulate multiple errors that occurred.
	Errors []*Error `json:"errors,omitempty" yaml:"errors,omitempty" msgpack:"errors,omitempty"`
	// Timestamp is the time in which the error occurred in UTC.
	Timestamp time.Time `json:"timestamp" yaml:"timestamp" msgpack:"timestamp"`
}

type Metadata map[string]interface{}

func From(err error) *Error {
	if err == nil {
		return nil
	}
	if errz, ok := err.(*Error); ok {
		return errz
	}
	return Wrap(err, Unknown)
}

func New(code ErrCode, message ...string) *Error {
	var messageItem string
	if len(message) > 0 {
		messageItem = message[0]
	}
	return &Error{
		Type:      code.String(),
		Code:      code,
		Status:    code.HTTPStatus(),
		Message:   messageItem,
		Timestamp: time.Now().UTC(),
	}
}

func Wrap(err error, code ErrCode, message ...string) *Error {
	var messageItem string
	if len(message) > 0 {
		messageItem = message[0]
	}
	return &Error{
		Type:      code.String(),
		Code:      code,
		Status:    code.HTTPStatus(),
		Message:   messageItem,
		Err:       err,
		Timestamp: time.Now().UTC(),
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) WithType(t string) *Error {
	e.Type = t
	return e
}

func (e *Error) WithTitle(format string, args ...interface{}) *Error {
	e.Title = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) WithMessage(format string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) WithDetails(details interface{}) *Error {
	e.Details = details
	return e
}

func (e *Error) WithMetadata(metadata Metadata) *Error {
	e.Metadata = metadata
	return e
}

func (e *Error) WithError(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) WithInstance(instance string) *Error {
	e.Instance = instance
	return e
}
