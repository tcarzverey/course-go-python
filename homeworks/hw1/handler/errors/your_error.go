package errors

import "fmt"

type AdditionalMessageError struct {
	msg string
	err error
}

func NewAdditionalMessageError(err error, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return &AdditionalMessageError{msg: msg, err: err}
}

func (e *AdditionalMessageError) Error() string {
	if e == nil {
		return "<nil>"
	}
	switch {
	case e.err == nil && e.msg == "":
		return ""
	case e.err == nil:
		return e.msg
	case e.msg == "":
		return e.err.Error()
	default:
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
}

// Unwrap распаковать внутреннюю ошибку, для ее участия в error-chain.
func (e *AdditionalMessageError) Unwrap() error { return e.err }

// Компиляционная проверка *AdditionalMessageError реализует интерфейс error.
var _ error = (*AdditionalMessageError)(nil)
