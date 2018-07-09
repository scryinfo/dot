package dot

import "errors"

type ErrorI interface {
	error
	Id() string
}

type sError struct {
	err error
	id  string
}

func (c *sError) Id() string {
	return c.id
}

func (c *sError) Error() string {
	return c.err.Error()
}

func NewSError(id string, info string) ErrorI {
	err := sError{err: errors.New(info), id: id}
	return &err
}

type Err struct {
	ErrNullParameter ErrorI
	ErrExited        ErrorI
	ErrParameter     ErrorI
}

var Error = &Err{}

func init() {
	Error.ErrNullParameter = NewSError("dot_null_parameter", "the parameter is null")
	Error.ErrExited = NewSError("dot_exited", "the value exited")
	Error.ErrParameter = NewSError("dot_error_parameter", "the parameter error")
}
