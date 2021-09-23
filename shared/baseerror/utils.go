package baseerror

type CodeError struct {
	code int64
	desc string
	data interface{}
}

func (err *CodeError) Error() string {
	return err.desc
}

func (err *CodeError) Code() int64 {
	return err.code
}

func (err *CodeError) Desc() string {
	return err.desc
}

func (err *CodeError) Data() interface{} {
	return err.data
}

func NewDefaultError(desc string) *CodeError {
	return NewCodeError(401, desc)
}

func NewCodeError(code int64, desc string) *CodeError {
	return &CodeError{
		code: code,
		desc: desc,
	}
}

func FromError(err error) (codeErr *CodeError, ok bool) {
	if se, ok := err.(*CodeError); ok {
		return se, true
	}
	return nil, false
}
