package util

import (
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFetchHTMLItemInnerText(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			result := FetchHTMLItemInnerText("https://www.baidu.com/", "#hotsearch-refresh-btn > span")
			logrus.Infoln("result", result)
			convey.So(result, convey.ShouldEqual, "换一换")
		})
	})
}

func TestFetchDynamicHTMLItemInnerText(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			result, err := FetchDynamicHTMLItemInnerText("https://mvnrepository.com/artifact/org.springframework/spring-core",
				"#maincontent > table > tbody > tr:nth-child(1) > td > span")
			logrus.Infoln("result", result, "error", err)
			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldEqual, "Apache 2.0")
		})
	})
}
