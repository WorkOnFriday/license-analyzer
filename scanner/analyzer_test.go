package scanner

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	fmt.Println(Check("LGPL-2.1-only", "LGPL-2.1-or-later"))
	fmt.Println(Check("LGPL-2.1-only", "GPL-2.0-only"))
	fmt.Println(Check("LGPL-2.1-only", "Sleepycat"))
	fmt.Println(Check("LGPL-2.1-only", "MS-RL"))
	fmt.Println(Check("fwhLicense", "MS-RL"))
}

func TestRecommand(t *testing.T) {
	var testArr1 = []string{"GPL-2.0-only"}
	fmt.Println(Recommand(testArr1))
	var testArr2 = []string{"GPL-2.0-only", "EUPL-1.1", "MS-RL"}
	fmt.Println(Recommand(testArr2))
	var testArr3 = []string{"Apache-2.0", "GPL-2.0-only"}
	fmt.Println(Recommand(testArr3))
}
