package errorx

import (
	"errors"
	"fmt"
)

type (
	CodeError interface {
		error
		Status() int
		Code() int
	}

	codeError struct {
		status int
		code   int
		desc   string
	}
)

func (err *codeError) Error() string {
	return err.desc
}

func (err *codeError) Code() int {
	return err.code
}

func (err *codeError) Status() int {
	return err.status
}

func (err *codeError) String() string {
	return fmt.Sprintf("Status: %d, Code: %d, Desc: %s", err.status, err.code, err.desc)
}

func NewCodeError(code int, desc string) CodeError {
	return NewStatCodeError(400, code, desc)
}

func NewDefaultError(desc string) CodeError {
	return NewStatCodeError(400, 406001, desc)
}

func NewStatCodeError(status, code int, desc string) CodeError {
	return &codeError{
		status: status,
		code:   code,
		desc:   desc,
	}
}

func FromError(err error) (CodeError, bool) {
	if err == nil {
		return nil, false
	}
	var ce CodeError
	if ok := errors.As(err, &ce); ok {
		return ce, ok
	}

	return nil, false
}
