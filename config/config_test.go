package config

import (
	"testing"

	"github.com/rogeecn/fabfile"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoad(t *testing.T) {
	Convey("Test Load", t, func() {
		file, err := fabfile.Find(".fun_test.yaml")
		So(err, ShouldBeNil)

		err = Load(file)
		So(err, ShouldBeNil)
	})
}
