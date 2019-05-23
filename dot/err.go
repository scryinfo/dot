// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import "errors"

var (
	_ Errorer = (*sError)(nil)
)

//Errorer dot error interface
type Errorer interface {
	error
	Code() string
	AddNewError(info string) Errorer
}

type sError struct {
	err  error
	code string
}

//Code error id
func (c *sError) Code() string {
	return c.code
}

func (c *sError) AddNewError(info string) Errorer {
	return NewError(c.Code(), c.Error()+info)
}

//SError return error info
func (c *sError) Error() string {
	return c.err.Error()
}

//NewError new Errorer
func NewError(code string, info string) Errorer {
	err := sError{err: errors.New(info), code: code}
	return &err
}

//Error dot error
type Error struct {
	NilParameter     Errorer
	Existed          Errorer
	NotExisted       Errorer
	Parameter        Errorer
	RelyTypeNotMatch Errorer
	TypeIdEmpty      Errorer
	Config           Errorer
	NoDotNewer       Errorer
	NotStruct        Errorer
	DotInvalid       Errorer
}

//SError error object frequently used by dot
var SError = &Error{}

func init() {
	SError.NilParameter = NewError("dot_null_parameter", "the parameter is null ")
	SError.Existed = NewError("dot_existed", "the value exited: ")
	SError.NotExisted = NewError("dot_not_existed", "the not exited: ")
	SError.Parameter = NewError("dot_error_parameter", "the parameter error ")
	SError.RelyTypeNotMatch = NewError("dot_rely_type", "rely type not match ")
	SError.TypeIdEmpty = NewError("dot_typeid_null", "typeid null: ")
	SError.Config = NewError("dot_config", "config error: ")
	SError.NoDotNewer = NewError("dot_no_newer", "Do not newer: ")
	SError.NotStruct = NewError("dot_not_struct", "Not struct: ")
	SError.DotInvalid = NewError("dot_invalid", "dot invalid: ")
}
