package apperr

import "github.com/pkg/errors"

type Error struct {
	root        error
	err         error
	contextInfo map[string]string
}

func New(err string, root error) error {
	return &Error{
		root: root,
		err:  errors.New(err),
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(e.root, err)
}

func (e *Error) Unwrap() error {
	return e.err
}

//func (e *Error) Cause(err error) bool {
//	return e.root
//}

func AddContext(err error, field, message string) {
	err = errors.Cause(err)
	if err == nil {
		return
	}
	if customErr, ok := err.(*Error); ok {
		if customErr.contextInfo == nil {
			customErr.contextInfo = make(map[string]string)
		}
		customErr.contextInfo[field] = message
	}
}

func GetContext(err error) map[string]string {
	err = errors.Cause(err)
	if err == nil {
		return nil
	}
	if customErr, ok := err.(*Error); ok {
		return customErr.contextInfo
	}
	return nil
}

func SetContext(err error, contextInfo map[string]string) {
	if customErr, ok := err.(*Error); ok {
		customErr.contextInfo = contextInfo
	}
}

func SetRoot(err, root error) {
	err = errors.Cause(err)
	if err == nil {
		return
	}
	if customErr, ok := err.(*Error); ok {
		customErr.root = root
	}
}
