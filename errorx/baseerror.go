package errorx

import (
	"fmt"
)

type BaseError struct {
	Code    int
	Message string
}

func (b BaseError) Error() string {
	return fmt.Sprintf("errCode: %d - errMsg: %s", b.Code, b.Message)
}

type BizError struct {
	BaseError
}

func (b *BizError) WithMessage(msg string) *BizError {
	e := *b
	e.Message = msg
	return &e
}

func NewBizError(code int, msg string) *BizError {
	e := &BizError{}
	e.Code = code
	e.Message = msg
	return e
}

func IsBizError(err error) bool {
	switch err.(type) {
	case *BizError, BizError:
		return true
	default:
		return false
	}
}
