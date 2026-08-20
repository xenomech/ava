package api

import "errors"

func asError(err error, target **Error) bool {
	return errors.As(err, target)
}
