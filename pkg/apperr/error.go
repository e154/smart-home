// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package apperr

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	code string
	err  error
	root error
	errs validator.ValidationErrorsTranslations
}

func ErrorWithCode(code, err string, root error) error {
	return &Error{
		code: code,
		err:  errors.New(err),
		root: root,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.root.Error(), e.err.Error())
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Message() string {
	return e.err.Error()
}

func (e *Error) Root() string {
	return e.root.Error()
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(e.root, err)
}

func (e *Error) SetValidationErrors(errs validator.ValidationErrorsTranslations) {
	e.errs = errs
}

func (e *Error) ValidationErrors() validator.ValidationErrorsTranslations {
	return e.errs
}

// Unwrap do not use this method directly. Use errors.Unwrap()
func (e *Error) Unwrap() error {
	return e.root
}

func Message(err error) string {
	e := GetError(err)
	if e != nil && e.err != nil {
		return e.err.Error()
	}
	return ""
}

func Code(err error) string {
	e := GetError(err)
	if e != nil {
		return e.code
	}
	return ""
}

func Root(err error) string {
	e := GetError(err)
	if e != nil && e.root != nil {
		return e.root.Error()
	}
	return ""
}

func SetValidationErrors(err error, errs validator.ValidationErrorsTranslations) {
	e := GetError(err)
	if e == nil {
		return
	}
	e.errs = errs
}

func GetError(err error) *Error {
	for {
		if err == nil {
			return nil
		}
		if e, ok := err.(*Error); ok {
			return e
		}
		err = errors.Unwrap(err)
	}
}
