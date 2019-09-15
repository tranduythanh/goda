package goda

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInterfaceList(t *testing.T) {
	Convey("List", t, func() {
		Convey("UniqueString", func() {
			l := List{"a", "a", "b", "c", "d", "c", "b"}
			r := l.UniqueString()
			So(r, ShouldHaveLength, 4)
		})

		Convey("SortString", func() {
			l := List{"a", "a", "b", "c", "d", "c", "b"}
			r := l.SortString()
			So(r, ShouldHaveLength, 7)
			So(r[4], ShouldEqual, "c")
		})
	})
}
