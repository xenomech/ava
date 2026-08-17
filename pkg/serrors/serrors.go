package serrors

import "errors"

type coded struct {
	code    string
	message string
}

func (e *coded) Error() string { return e.message }

func (e *coded) Code() string { return e.code }

func New(text string) error {
	return errors.New(text)
}

func NewCoded(code, text string) error {
	return &coded{code: code, message: text}
}

type Coded interface {
	error
	Code() string
}

func CodeOf(err error) string {
	var c Coded
	if errors.As(err, &c) {
		return c.Code()
	}

	return ""
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
