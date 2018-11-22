package assertions

import (
	"fmt"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey/reporting"
)

const assertionSuccess = ""

var (
	ShouldEqual = assertions.ShouldEqual
	ShouldBeNil = assertions.ShouldBeNil
)

type assertion func(actual interface{}, expected ...interface{}) string

func So(actual interface{}, assert assertion, expected ...interface{}) {
	if result := assert(actual, expected...); result == assertionSuccess {
		fmt.Printf(".")
	} else {
		fmt.Println()
		fmt.Println(reporting.NewFailureReport(result))
	}
}
