package dot

import "errors"

var (
	_ Errorer = (*sError)(nil)
)

//Errorer dot error interface
type Errorer interface {
	error
	Code() string
}

type sError struct {
	err error
	id  string
}

//Code error id
func (c *sError) Code() string {
	return c.id
}

//SError return error info
func (c *sError) Error() string {
	return c.err.Error()
}

//New new Errorer
func New(id string, info string) Errorer {
	err := sError{err: errors.New(info), id: id}
	return &err
}

//Error dot error
type Error struct {
	ErrNullParameter Errorer
	ErrExited        Errorer
	ErrParameter     Errorer
}

//SError dot的全局常用 error对象
var SError = &Error{}

func init() {
	SError.ErrNullParameter = New("dot_null_parameter", "the parameter is null")
	SError.ErrExited = New("dot_exited", "the value exited")
	SError.ErrParameter = New("dot_error_parameter", "the parameter error")
}
