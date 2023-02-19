/*
 * Copyright © 2019 Hedzr Yeh.
 */

package util_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hedzr/consul-tags/util"
)

func TestExtract2(t *testing.T) {
	s1, s2 := util.Extract2("xcvb:dfgh:456", ":")
	if !strings.EqualFold(s1, "xcvb") {
		t.Logf("Extract2() return bad `s1`: '%s'", s1)
		t.Fail()
	}
	if !strings.EqualFold(s2, "dfgh:456") {
		t.Logf("Extract2() return bad `s2`: '%s'", s2)
		t.Fail()
	}
}

func ExampleExtract2() {
	const str, delim = "xcvb:dfgh456", ":"
	s1, s2 := util.Extract2("xcvb:dfgh:456", ":")
	fmt.Printf("s1 = %s, s2 = %s\n", s1, s2)
	// Output: s1 = xcvb, s2 = dfgh:456
}
