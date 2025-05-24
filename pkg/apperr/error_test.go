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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {

	var (
		baseErr   = errors.New("base error")
		commonErr = errors.New("common error")

		ErrFirst  = ErrorWithCode("FIRST_ERROR", "first error", baseErr)
		ErrSecond = ErrorWithCode("SECOND_ERROR", "second error", commonErr)
	)

	t.Run("error", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			firstError := ErrFirst
			ctx.So(errors.Is(firstError, baseErr), ShouldBeTrue)
			ctx.So(errors.Is(firstError, commonErr), ShouldBeFalse)

			secondError := ErrSecond
			ctx.So(errors.Is(secondError, commonErr), ShouldBeTrue)
			ctx.So(errors.Is(secondError, baseErr), ShouldBeFalse)

			ctx.So(ErrFirst.Error(), ShouldEqual, "base error: first error")
			ctx.So(ErrSecond.Error(), ShouldEqual, "common error: second error")

			ctx.So(Code(ErrFirst), ShouldEqual, "FIRST_ERROR")
			ctx.So(Message(ErrFirst), ShouldEqual, "first error")
			ctx.So(Root(ErrFirst), ShouldEqual, "base error")

			err := GetError(ErrFirst)
			ctx.So(err.Code(), ShouldEqual, "FIRST_ERROR")
			ctx.So(err.Message(), ShouldEqual, "first error")
			ctx.So(err.Root(), ShouldEqual, "base error")

			ctx.So(errors.Is(secondError, secondError), ShouldBeTrue)
			ctx.So(errors.Is(firstError, firstError), ShouldBeTrue)
			ctx.So(errors.Is(secondError, firstError), ShouldBeFalse)
		})
	})
}
