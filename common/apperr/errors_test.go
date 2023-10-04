// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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
	"fmt"
	"testing"

	"github.com/pkg/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {

	t.Run("error1", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			baseErr := errors.New("base")
			err := errors.Wrap(baseErr, "first")
			err = errors.Wrap(err, "second")
			err = errors.Wrap(err, ErrDashboardImport.Error())

			//fmt.Println("---+", errors.Cause(err))

			ctx.So(errors.Is(err, baseErr), ShouldBeTrue)
			ctx.So(errors.Is(err, ErrDashboardImport), ShouldBeFalse)

			for {
				err = errors.Unwrap(err)
				if err == nil {
					break
				}
				//fmt.Println("--->", err.Error())
			}
		})
	})

	t.Run("error2", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			err := errors.Wrap(ErrDashboardImport, "first")
			AddContext(err, "name", "not found")
			err = errors.Wrap(err, "second")

			//fmt.Println("---+", errors.Cause(err))
			ctx.So(errors.Is(err, ErrDashboardImport), ShouldBeTrue)
			ctx.So(errors.Is(err, ErrInternal), ShouldBeTrue)

			SetRoot(err, ErrNotFound)

			ctx.So(errors.Is(err, ErrDashboardImport), ShouldBeTrue)
			ctx.So(errors.Is(err, ErrInternal), ShouldBeFalse)
			ctx.So(errors.Is(err, ErrNotFound), ShouldBeTrue)

			c := GetContext(errors.Cause(err))
			fmt.Println(c)

		})
	})

}
