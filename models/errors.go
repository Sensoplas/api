package models

type ClientErrorType int

const (
	BadRequest ClientErrorType = iota
	Unauthorized
	Forbidden
	NotFound
	Conflict
)

type ServerErrorType int

const (
	Base ServerErrorType = iota
	NotImplemented
	ServiceUnavailable
)

type ClientError interface {
	Error() string
	ClientMessage() string
	Type() ClientErrorType
}

type ServerError interface {
	Error() string
	ClientMessage() string
	Type() ServerErrorType
}

var _ error = &GenericClientError{}
var _ ClientError = &GenericClientError{}

func NewGenericClientError(message string, e error) *GenericClientError {
	return &GenericClientError{message, e, BadRequest}
}

type GenericClientError struct {
	message string
	e       error
	t       ClientErrorType
}

func (e *GenericClientError) Error() string {
	return e.e.Error()
}

func (e *GenericClientError) Unwrap() error {
	return e.e
}

func (e *GenericClientError) ClientMessage() string {
	return e.message
}

func (e *GenericClientError) Type() ClientErrorType {
	return e.t
}

var _ error = &GenericServerError{}
var _ ServerError = &GenericServerError{}

type GenericServerError struct {
	message string
	e       error
	t       ServerErrorType
}

func (e *GenericServerError) Error() string {
	return e.e.Error()
}

func (e *GenericServerError) Unwrap() error {
	return e.e
}

func (e *GenericServerError) ClientMessage() string {
	return e.message
}

func (e *GenericServerError) Type() ServerErrorType {
	return e.t
}
