package scanner

/*
测试许可证全名转简写
测试检查许可证冲突
测试推荐的许可证无冲突
*/

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLicenseLongNameToShort(t *testing.T) {
	InitializeAnalyzer()
	t.Run("TestLicenseLongNameToShort", func(t *testing.T) {
		testCases := []struct {
			longName  string
			shortName string
		}{
			// 允许无GNU
			{"GNU LESSER GENERAL PUBLIC LICENSE Version 2.1", "LGPL-2.1-or-later"},
			{"LESSER GENERAL PUBLIC LICENSE Version 2.1", "LGPL-2.1-or-later"},
			// 允许空格增减
			{"EUROPEAN UNION PUBLIC LICENCE V. 1.1", "EUPL-1.1"},
			{"EUROPEAN UNION PUBLIC LICENCE V.1.1", "EUPL-1.1"},

			{"MICROSOFT RECIPROCAL LICENSE", "MS-RL"},
		}
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("Test%d", i), func(t *testing.T) {
				convey.Convey("Test", t, func() {
					convey.So(LicenseLongNameToShort(testCase.longName), convey.ShouldEqual, testCase.shortName)
				})
			})
		}
	})
}

func TestCheckLicenseConflictByShortName(t *testing.T) {
	InitializeAnalyzer()
	t.Run("TestCheckLicenseConflictByShortName", func(t *testing.T) {
		testCases := []struct {
			mainLicense string
			libLicense  string
			result      ConflictResult
		}{
			{"LGPL-2.1-only", "LGPL-2.1-or-later",
				ConflictResult{false, true, pass}},
			{"LGPL-2.1-only", "GPL-2.0-only",
				ConflictResult{false, true, "相容, 组合遵循GPL-2.0-only"}},
			{"LGPL-2.1-only", "Sleepycat",
				ConflictResult{false, true, "相容, 并入的代码遵循LGPL-2.1-only"}},
			{"LGPL-2.1-only", "MS-RL",
				ConflictResult{false, false, fail}},
			{"fwhLicense", "MS-RL",
				ConflictResult{true, false, unknown}},
		}
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("Test%d", i), func(t *testing.T) {
				convey.Convey("Test", t, func() {
					result := CheckLicenseConflictByShortName(testCase.mainLicense, testCase.libLicense)
					convey.So(result.Unknown, convey.ShouldEqual, testCase.result.Unknown)
					convey.So(result.Pass, convey.ShouldEqual, testCase.result.Pass)
					convey.So(result.Message, convey.ShouldEqual, testCase.result.Message)
				})
			})
		}
	})
}

func TestRecommendByLibraryLicenseLongName(t *testing.T) {
	InitializeAnalyzer()
	t.Run("TestRecommend", func(t *testing.T) {
		testCases := []struct {
			libLicenses []string
		}{
			{libLicenses: []string{"GPL-2.0-only"}},
			{libLicenses: []string{"GPL-2.0-only", "EUPL-1.1", "MS-RL"}},
			{libLicenses: []string{"Apache-2.0", "GPL-2.0-only"}},
		}
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("Test%d", i), func(t *testing.T) {
				convey.Convey("Test", t, func() {
					var recommendLicenses = RecommendByLibraryLicenseShortName(testCase.libLicenses)
					for _, recommendLicense := range recommendLicenses {
						for _, libLicense := range testCase.libLicenses {
							conflictResult := CheckLicenseConflictByShortName(recommendLicense, libLicense)
							convey.So(conflictResult.Pass, convey.ShouldBeTrue)
						}
					}
				})
			})
		}
	})
}
