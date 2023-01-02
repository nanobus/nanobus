// Code generated by @apexlang/codegen. DO NOT EDIT.

package runtime

import (
	"encoding/json"
	"fmt"

	"github.com/nanobus/nanobus/pkg/expr"
)

// A mapping of target resource name to source name.
type ResourceLinks map[string]string

// A mapping of interface to operation authorizations.
type Authorizations map[string]AuthOperations

// A mapping of operation to authorization rules.
type AuthOperations map[string]Authorization

// A map of interface name to operation pipelines.
type Interfaces map[string]Operations

// A map of operation name to pipeline.
type Operations map[string]Pipeline

// The resource configuration that is tied to an application at runtime.
type ResourcesConfig struct {
	// A mapping of resource name to component configuration.
	Resources map[string]Component `json:"resources" yaml:"resources" msgpack:"resources" mapstructure:"resources" validate:"dive"`
	// A mapping of Itoa to resource links.
	ResourceLinks map[string]ResourceLinks `json:"resourceLinks,omitempty" yaml:"resourceLinks,omitempty" msgpack:"resourceLinks,omitempty" mapstructure:"resourceLinks" validate:"dive"`
}

// The main Iota component/application configuration.
type BusConfig struct {
	// The Iota identifier.
	ID string `json:"id" yaml:"id" msgpack:"id" mapstructure:"id" validate:"required"`
	// The Iota version.
	Version string `json:"version" yaml:"version" msgpack:"version" mapstructure:"version" validate:"required"`
	// The main code binary (wasm or native).
	Main *string `json:"main,omitempty" yaml:"main,omitempty" msgpack:"main,omitempty" mapstructure:"main"`
	// The interface definition.
	Spec *string `json:"spec,omitempty" yaml:"spec,omitempty" msgpack:"spec,omitempty" mapstructure:"spec"`
	// Package specifies the contents of the Iota OCI image.
	Package *Package `json:"package,omitempty" yaml:"package,omitempty" msgpack:"package,omitempty" mapstructure:"package"`
	// Resources are externally configured sources and receivers of data (DB, REST
	// endpoint).
	Resources []string `json:"resources,omitempty" yaml:"resources,omitempty" msgpack:"resources,omitempty" mapstructure:"resources" validate:"dive"`
	// Other Iotas that this Iota depends on using.
	Includes map[string]Reference `json:"includes,omitempty" yaml:"includes,omitempty" msgpack:"includes,omitempty" mapstructure:"includes" validate:"dive"`
	// Tracing configures an Open Telemetry span exporter.
	Tracing *Component  `json:"tracing,omitempty" yaml:"tracing,omitempty" msgpack:"tracing,omitempty" mapstructure:"tracing"`
	Specs   []Component `json:"specs,omitempty" yaml:"specs,omitempty" msgpack:"specs,omitempty" mapstructure:"specs" validate:"dive"`
	Compute []Component `json:"compute,omitempty" yaml:"compute,omitempty" msgpack:"compute,omitempty" mapstructure:"compute" validate:"dive"`
	// Resiliency defines policies for fault tolerance.
	Resiliency *Resiliency `json:"resiliency,omitempty" yaml:"resiliency,omitempty" msgpack:"resiliency,omitempty" mapstructure:"resiliency"`
	// Codecs configure how data formats are encoded (and persisted) and decoded (and
	// received).
	Codecs map[string]Component `json:"codecs,omitempty" yaml:"codecs,omitempty" msgpack:"codecs,omitempty" mapstructure:"codecs" validate:"dive"`
	// Initializers handle startup and initialization tasks which execute and freed up.
	Initializers map[string]Component `json:"initializers,omitempty" yaml:"initializers,omitempty" msgpack:"initializers,omitempty" mapstructure:"initializers" validate:"dive"`
	// Transports configure inbound communication mechanisms from clients (e.g. HTTP or
	// gRPC) or event sources (e.g Message brokers).
	Transports map[string]Component `json:"transports,omitempty" yaml:"transports,omitempty" msgpack:"transports,omitempty" mapstructure:"transports" validate:"dive"`
	// Filters process received data immediately after the transports.
	Filters []Component `json:"filters,omitempty" yaml:"filters,omitempty" msgpack:"filters,omitempty" mapstructure:"filters" validate:"dive"`
	// Processing pipeline for pre-authoration actions.
	Preauth Interfaces `json:"preauth,omitempty" yaml:"preauth,omitempty" msgpack:"preauth,omitempty" mapstructure:"preauth"`
	// Authorization rules for interfaces.
	Authorization Authorizations `json:"authorization,omitempty" yaml:"authorization,omitempty" msgpack:"authorization,omitempty" mapstructure:"authorization"`
	// Processing pipeline for post-authoration actions.
	Postauth Interfaces `json:"postauth,omitempty" yaml:"postauth,omitempty" msgpack:"postauth,omitempty" mapstructure:"postauth"`
	// Interface composed from declarative pipelines.
	Interfaces Interfaces `json:"interfaces,omitempty" yaml:"interfaces,omitempty" msgpack:"interfaces,omitempty" mapstructure:"interfaces"`
	// Pipelines that preform data access (typically using resources) on behalf of the
	// application.
	Providers Interfaces               `json:"providers,omitempty" yaml:"providers,omitempty" msgpack:"providers,omitempty" mapstructure:"providers"`
	Errors    map[string]ErrorTemplate `json:"errors,omitempty" yaml:"errors,omitempty" msgpack:"errors,omitempty" mapstructure:"errors" validate:"dive"`
}

// Package defines the contents of the OCI image.
type Package struct {
	// The OCI registry host.
	Registry *string `json:"registry,omitempty" yaml:"registry,omitempty" msgpack:"registry,omitempty" mapstructure:"registry"`
	// The OCI registry organization
	Org *string `json:"org,omitempty" yaml:"org,omitempty" msgpack:"org,omitempty" mapstructure:"org"`
	// The files and directories to include.
	Add []string `json:"add" yaml:"add" msgpack:"add" mapstructure:"add" validate:"dive"`
}

// A reference to another Iota as a dependency.
type Reference struct {
	// The OCI reference or directory.
	Ref string `json:"ref" yaml:"ref" msgpack:"ref" mapstructure:"ref" validate:"required"`
	// Used to pass configured resources to Iota dependencies. It is a mapping of the
	// referred Iota resources to this Iotas resources.
	ResourceLinks map[string]string `json:"resourceLinks,omitempty" yaml:"resourceLinks,omitempty" msgpack:"resourceLinks,omitempty" mapstructure:"resourceLinks" validate:"dive"`
}

// A loadable component and its configuration.
type Component struct {
	// The component name.
	Uses string `json:"uses" yaml:"uses" msgpack:"uses" mapstructure:"uses" validate:"required"`
	// The component's configuration where the structure is defined by the component
	// implementation.
	With interface{} `json:"with,omitempty" yaml:"with,omitempty" msgpack:"with,omitempty" mapstructure:"with"`
}

// A section where resiliency policies are configured and given reusable reference
// names.
type Resiliency struct {
	// Timeout durations.
	Timeouts map[string]Duration `json:"timeouts,omitempty" yaml:"timeouts,omitempty" msgpack:"timeouts,omitempty" mapstructure:"timeouts" validate:"dive"`
	// Retry handling using constant and exponential backoff policies.
	Retries map[string]Backoff `json:"retries,omitempty" yaml:"retries,omitempty" msgpack:"retries,omitempty" mapstructure:"retries" validate:"dive"`
	// Prevent further calls to remote resources that are unavailable or operating
	// abnormally.
	CircuitBreakers map[string]CircuitBreaker `json:"circuitBreakers,omitempty" yaml:"circuitBreakers,omitempty" msgpack:"circuitBreakers,omitempty" mapstructure:"circuitBreakers" validate:"dive"`
}

// A backoff policy that always returns a fixed backoff delay.
type ConstantBackoff struct {
	// The duration to wait between retry attempts.
	Duration Duration `json:"duration" yaml:"duration" msgpack:"duration" mapstructure:"duration"`
	// The maximum number of retries to attempt. No value denotes an indefinite number
	// of retries.
	MaxRetries *uint32 `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty" msgpack:"maxRetries,omitempty" mapstructure:"maxRetries"`
}

// A backoff implementation that increases the backoff period for each retry
// attempt using a randomization function that grows exponentially.
//
// The exponential back-off window uses the following formula:
//
// ``` BackOffDuration = PreviousBackOffDuration * (Random value from 0.5 to 1.5) *
// 1.5 if BackOffDuration > maxInterval {   BackoffDuration = maxInterval } ```
type ExponentialBackoff struct {
	// The initial interval.
	InitialInterval     Duration `json:"initialInterval" yaml:"initialInterval" msgpack:"initialInterval" mapstructure:"initialInterval"`
	RandomizationFactor float64  `json:"randomizationFactor" yaml:"randomizationFactor" msgpack:"randomizationFactor" mapstructure:"randomizationFactor"`
	Multiplier          float64  `json:"multiplier" yaml:"multiplier" msgpack:"multiplier" mapstructure:"multiplier"`
	// Determines the maximum interval between retries to which the exponential
	// back-off policy can grow. Additional retries always occur after a duration of
	// `maxInterval`.
	MaxInterval    Duration `json:"maxInterval" yaml:"maxInterval" msgpack:"maxInterval" mapstructure:"maxInterval"`
	MaxElapsedTime Duration `json:"maxElapsedTime" yaml:"maxElapsedTime" msgpack:"maxElapsedTime" mapstructure:"maxElapsedTime"`
	// The maximum number of retries to attempt. No value denotes an indefinite number
	// of retries.
	MaxRetries *uint32 `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty" msgpack:"maxRetries,omitempty" mapstructure:"maxRetries"`
}

type CircuitBreaker struct {
	// The maximum number of requests allowed to pass through when the circuit breaker
	// is half-open (recovering from failure).
	MaxRequests uint32 `json:"maxRequests" yaml:"maxRequests" msgpack:"maxRequests" mapstructure:"maxRequests"`
	// The cyclical period of time used by the CB to clear its internal counts. If set
	// to 0 seconds, this never clears.
	Interval Duration `json:"interval" yaml:"interval" msgpack:"interval" mapstructure:"interval"`
	// The period of the open state (directly after failure) until the circuit breaker
	// switches to half-open. Defaults to 60s.
	Timeout Duration `json:"timeout" yaml:"timeout" msgpack:"timeout" mapstructure:"timeout"`
	// A Common Expression Language (CEL) statement that is evaluated by the circuit
	// breaker. When the statement evaluates to true, the CB trips and becomes open.
	// Default is consecutiveFailures > 5.
	Trip *expr.ValueExpr `json:"trip,omitempty" yaml:"trip,omitempty" msgpack:"trip,omitempty" mapstructure:"trip"`
}

// An authorization rule that can assert based on the existence or quality of a
// claim, or through authorization components.
type Authorization struct {
	// This flag must be explicitly set to `true` if unauthenticated/anonymous access
	// is allowed.
	Unauthenticated bool                   `json:"unauthenticated" yaml:"unauthenticated" msgpack:"unauthenticated" mapstructure:"unauthenticated"`
	Has             []string               `json:"has" yaml:"has" msgpack:"has" mapstructure:"has" validate:"dive"`
	Checks          map[string]interface{} `json:"checks" yaml:"checks" msgpack:"checks" mapstructure:"checks" validate:"dive"`
	Rules           []Component            `json:"rules" yaml:"rules" msgpack:"rules" mapstructure:"rules" validate:"dive"`
}

// A
type Pipeline struct {
	// The pipeline name.
	Name  string `json:"name" yaml:"name" msgpack:"name" mapstructure:"name" validate:"required"`
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty" msgpack:"steps,omitempty" mapstructure:"steps" validate:"dive"`
}

type Step struct {
	Name           string      `json:"name" yaml:"name" msgpack:"name" mapstructure:"name" validate:"required"`
	Call           *string     `json:"call,omitempty" yaml:"call,omitempty" msgpack:"call,omitempty" mapstructure:"call"`
	Uses           string      `json:"uses" yaml:"uses" msgpack:"uses" mapstructure:"uses" validate:"required"`
	With           interface{} `json:"with,omitempty" yaml:"with,omitempty" msgpack:"with,omitempty" mapstructure:"with"`
	Returns        *string     `json:"returns,omitempty" yaml:"returns,omitempty" msgpack:"returns,omitempty" mapstructure:"returns"`
	Timeout        *string     `json:"timeout,omitempty" yaml:"timeout,omitempty" msgpack:"timeout,omitempty" mapstructure:"timeout"`
	Retry          *string     `json:"retry,omitempty" yaml:"retry,omitempty" msgpack:"retry,omitempty" mapstructure:"retry"`
	CircuitBreaker *string     `json:"circuitBreaker,omitempty" yaml:"circuitBreaker,omitempty" msgpack:"circuitBreaker,omitempty" mapstructure:"circuitBreaker"`
	OnError        *Pipeline   `json:"onError,omitempty" yaml:"onError,omitempty" msgpack:"onError,omitempty" mapstructure:"onError"`
}

type ErrorTemplate struct {
	Type    string             `json:"type" yaml:"type" msgpack:"type" mapstructure:"type" validate:"required"`
	Code    ErrCode            `json:"code" yaml:"code" msgpack:"code" mapstructure:"code"`
	Status  uint32             `json:"status" yaml:"status" msgpack:"status" mapstructure:"status"`
	Title   *expr.Text         `json:"title" yaml:"title" msgpack:"title" mapstructure:"title"`
	Message *expr.Text         `json:"message" yaml:"message" msgpack:"message" mapstructure:"message"`
	Path    *string            `json:"path,omitempty" yaml:"path,omitempty" msgpack:"path,omitempty" mapstructure:"path"`
	Help    *expr.Text         `json:"help,omitempty" yaml:"help,omitempty" msgpack:"help,omitempty" mapstructure:"help"`
	Locales map[string]Strings `json:"locales,omitempty" yaml:"locales,omitempty" msgpack:"locales,omitempty" mapstructure:"locales" validate:"dive"`
}

type Strings struct {
	Title   *expr.Text `json:"title" yaml:"title" msgpack:"title" mapstructure:"title"`
	Message *expr.Text `json:"message" yaml:"message" msgpack:"message" mapstructure:"message"`
}

type Backoff struct {
	Constant    *ConstantBackoff    `json:"constant,omitempty" yaml:"constant,omitempty" msgpack:"constant,omitempty" validate:"required_without=Exponential"`
	Exponential *ExponentialBackoff `json:"exponential,omitempty" yaml:"exponential,omitempty" msgpack:"exponential,omitempty" validate:"required_without=Constant"`
}

type ErrCode int32

const (
	// OK indicates the operation was successful.
	ErrCodeOk ErrCode = 0
	// Canceled indicates the operation was canceled (typically by the caller).
	ErrCodeCanceled ErrCode = 1
	// Unknown error. An example of where this error may be returned is if a Status
	// value received from another address space belongs to an error-space that is not
	// known in this address space. Also errors raised by APIs that do not return
	// enough error information may be converted to this error.
	ErrCodeUnknown ErrCode = 2
	// InvalidArgument indicates client specified an invalid argument. Note that this
	// differs from FailedPrecondition. It indicates arguments that are problematic
	// regardless of the state of the system (e.g., a malformed file name).
	ErrCodeInvalidArgument ErrCode = 3
	// DeadlineExceeded means operation expired before completion. For operations that
	// change the state of the system, this error may be returned even if the operation
	// has completed successfully. For example, a successful response from a server
	// could have been delayed long enough for the deadline to expire.
	ErrCodeDeadlineExceeded ErrCode = 4
	// NotFound means some requested entity (e.g., file or directory) was not found.
	ErrCodeNotFound ErrCode = 5
	// AlreadyExists means an attempt to create an entity failed because one already
	// exists.
	ErrCodeAlreadyExists ErrCode = 6
	// PermissionDenied indicates the caller does not have permission to execute the
	// specified operation. It must not be used for rejections caused by exhausting
	// some resource (use ResourceExhausted instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated instead for those
	// errors).
	ErrCodePermissionDenied ErrCode = 7
	// ResourceExhausted indicates some resource has been exhausted, perhaps a per-user
	// quota, or perhaps the entire file system is out of space.
	ErrCodeResourceExhausted ErrCode = 8
	// FailedPrecondition indicates operation was rejected because the system is not in
	// a state required for the operation's execution. For example, directory to be
	// deleted may be non-empty, an rmdir operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding between
	// FailedPrecondition, Aborted, and Unavailable:  (a) Use Unavailable if the client
	// can retry just the failing call.  (b) Use Aborted if the client should retry at
	// a higher-level      (e.g., restarting a read-modify-write sequence).  (c) Use
	// FailedPrecondition if the client should not retry until      the system state
	// has been explicitly fixed. E.g., if an "rmdir"      fails because the directory
	// is non-empty, FailedPrecondition      should be returned since the client should
	// not retry unless      they have first fixed up the directory by deleting files
	// from it.  (d) Use FailedPrecondition if the client performs conditional
	// Get/Update/Delete on a resource and the resource on the      server does not
	// match the condition. E.g., conflicting      read-modify-write on the same
	// resource.
	ErrCodeFailedPrecondition ErrCode = 9
	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	//
	// See litmus test above for deciding between FailedPrecondition, Aborted, and
	// Unavailable.
	ErrCodeAborted ErrCode = 10
	// OutOfRange means operation was attempted past the valid range. E.g., seeking or
	// reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may be fixed if the
	// system state changes. For example, a 32-bit file system will generate
	// InvalidArgument if asked to read at an offset that is not in the range
	// [0,2^32-1], but it will generate OutOfRange if asked to read from an offset past
	// the current file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and OutOfRange. We
	// recommend using OutOfRange (the more specific error) when it applies so that
	// callers who are iterating through a space can easily look for an OutOfRange
	// error to detect when they are done.
	ErrCodeOutOfRange ErrCode = 11
	// Unimplemented indicates operation is not implemented or not supported/enabled in
	// this service.
	ErrCodeUnimplemented ErrCode = 12
	// Internal errors. Means some invariants expected by underlying system has been
	// broken. If you see one of these errors, something is very broken.
	ErrCodeInternal ErrCode = 13
	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a backoff.
	// Note that it is not always safe to retry non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition, Aborted, and
	// Unavailable.
	ErrCodeUnavailable ErrCode = 14
	// DataLoss indicates unrecoverable data loss or corruption.
	ErrCodeDataLoss ErrCode = 15
	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	ErrCodeUnauthenticated ErrCode = 16
)

var toStringErrCode = map[ErrCode]string{
	ErrCodeOk:                 "ok",
	ErrCodeCanceled:           "canceled",
	ErrCodeUnknown:            "unknown",
	ErrCodeInvalidArgument:    "invalid_argument",
	ErrCodeDeadlineExceeded:   "deadline_exceeded",
	ErrCodeNotFound:           "not_found",
	ErrCodeAlreadyExists:      "already_exists",
	ErrCodePermissionDenied:   "permission_denied",
	ErrCodeResourceExhausted:  "resource_exhausted",
	ErrCodeFailedPrecondition: "failed_precondition",
	ErrCodeAborted:            "aborted",
	ErrCodeOutOfRange:         "out_of_range",
	ErrCodeUnimplemented:      "unimplemented",
	ErrCodeInternal:           "internal",
	ErrCodeUnavailable:        "unavailable",
	ErrCodeDataLoss:           "data_loss",
	ErrCodeUnauthenticated:    "unauthenticated",
}

var toIDErrCode = map[string]ErrCode{
	"ok":                  ErrCodeOk,
	"canceled":            ErrCodeCanceled,
	"unknown":             ErrCodeUnknown,
	"invalid_argument":    ErrCodeInvalidArgument,
	"deadline_exceeded":   ErrCodeDeadlineExceeded,
	"not_found":           ErrCodeNotFound,
	"already_exists":      ErrCodeAlreadyExists,
	"permission_denied":   ErrCodePermissionDenied,
	"resource_exhausted":  ErrCodeResourceExhausted,
	"failed_precondition": ErrCodeFailedPrecondition,
	"aborted":             ErrCodeAborted,
	"out_of_range":        ErrCodeOutOfRange,
	"unimplemented":       ErrCodeUnimplemented,
	"internal":            ErrCodeInternal,
	"unavailable":         ErrCodeUnavailable,
	"data_loss":           ErrCodeDataLoss,
	"unauthenticated":     ErrCodeUnauthenticated,
}

func (e ErrCode) String() string {
	str, ok := toStringErrCode[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *ErrCode) FromString(str string) error {
	var ok bool
	*e, ok = toIDErrCode[str]
	if !ok {
		return fmt.Errorf("unknown value %q for ErrCode", str)
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e ErrCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *ErrCode) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}
