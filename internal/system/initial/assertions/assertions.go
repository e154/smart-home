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

package assertions

import (
	"fmt"
	"log"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey/reporting"
)

const assertionSuccess = ""

var (
	// ShouldEqual ...
	ShouldEqual = assertions.ShouldEqual
	// ShouldBeNil ...
	ShouldBeNil = assertions.ShouldBeNil
	// ShouldBeZeroValue ...
	ShouldBeZeroValue = assertions.ShouldBeZeroValue
)

type assertion func(actual interface{}, expected ...interface{}) string

// So ...
func So(actual interface{}, assert assertion, expected ...interface{}) {
	if result := assert(actual, expected...); result == assertionSuccess {
		fmt.Printf(".")
	} else {
		log.Panicf("%v", reporting.NewFailureReport(result))
	}
}
