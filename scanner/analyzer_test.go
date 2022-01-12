package scanner

import (
	"github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("TestCheck", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			convey.So(Check("LGPL-2.1-only", "LGPL-2.1-or-later"), convey.ShouldEqual, "相容")
			convey.So(Check("LGPL-2.1-only", "GPL-2.0-only"), convey.ShouldEqual, "相容, 组合遵循GPL-2.0-only")
			convey.So(Check("LGPL-2.1-only", "Sleepycat"), convey.ShouldEqual, "相容, 并入的代码遵循LGPL-2.1-only")
			convey.So(Check("LGPL-2.1-only", "MS-RL"), convey.ShouldEqual, "冲突")
			convey.So(Check("fwhLicense", "MS-RL"), convey.ShouldEqual, "分析失败，许可证库不包含对应信息")
		})
	})
}

func TestRecommand(t *testing.T) {
	t.Run("TestRecommand", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			var testArr1 = []string{"GPL-2.0-only"}
			var result1 = Recommand(testArr1)
			for _, v1 := range result1 {
				for _, v2 := range testArr1 {
					convey.So(strings.HasPrefix(Check(v1, v2), pass), convey.ShouldBeTrue)
				}
			}
			var testArr2 = []string{"GPL-2.0-only", "EUPL-1.1", "MS-RL"}
			var result2 = Recommand(testArr2)
			for _, v1 := range result2 {
				for _, v2 := range testArr2 {
					convey.So(strings.HasPrefix(Check(v1, v2), pass), convey.ShouldBeTrue)
				}
			}
			var testArr3 = []string{"Apache-2.0", "GPL-2.0-only"}
			var result3 = Recommand(testArr3)
			for _, v1 := range result3 {
				for _, v2 := range testArr3 {
					convey.So(strings.HasPrefix(Check(v1, v2), pass), convey.ShouldBeTrue)
				}
			}
		})
	})
}
