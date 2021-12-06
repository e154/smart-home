// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package validation

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/fx"

	m "github.com/e154/smart-home/models"
)

// Validate ...
type Validate struct {
	validate *validator.Validate
	trans    ut.Translator
	config   *m.AppConfig
	lang     string
}

// NewValidate ...
func NewValidate(lc fx.Lifecycle,
	config *m.AppConfig) (v *Validate) {
	v = &Validate{
		lang: config.Lang,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return v.start(ctx)
		},
	})

	return
}

func (v *Validate) start(_ context.Context) (err error) {

	_en := en.New()
	uni := ut.New(_en, _en, ru.New(), es.New())

	var ok bool
	if v.trans, ok = uni.GetTranslator(v.lang); !ok {
		v.trans, _ = uni.GetTranslator("en")
	}

	v.validate = validator.New()
	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if tag == "-" {
			return ""
		}
		return tag
	})

	err = en_translations.RegisterDefaultTranslations(v.validate, v.trans)

	return
}

// ValidVar ...
func (v *Validate) ValidVar(s interface{}, key, tag string) (ok bool, errs validator.ValidationErrorsTranslations) {
	err := v.validate.Var(s, tag)
	if ok = err == nil; !ok {
		if validationErrors, valid := err.(validator.ValidationErrors); valid {
			errs = validationErrors.Translate(v.trans)
			errs[key] = errs[""]
			delete(errs, "")
		}
	}
	return
}

// Valid ...
func (v *Validate) Valid(s interface{}) (ok bool, errs validator.ValidationErrorsTranslations) {
	err := v.validate.Struct(s)
	if err != nil {
		if validationErrors, valid := err.(validator.ValidationErrors); valid {
			errs = validationErrors.Translate(v.trans)
		}
		return
	}
	ok = true
	return
}
