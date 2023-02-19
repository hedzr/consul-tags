/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul_test

import (
	"fmt"
	"testing"

	// _ set "github.com/deckarep/golang-set"
	"github.com/hedzr/consul-tags/objects/consul"
)

type TC struct {
	Tags                                []string
	Add                                 []string
	Remove                              []string
	Delimitor                           string
	hasClear, isPlainMode, isStringMode bool
	ExpectTags                          []string
}

func tagsModifyTestCases() func() map[string]TC {
	ETESTS := map[string]TC{
		// rm and add
		"0": {
			[]string{"version=1.0", "b"},
			[]string{"a,b,c"},
			[]string{"c,z"},
			"=",
			false, false, false,
			[]string{"version=1.0", "a", "b", "c"},
		},
		// rm and add, spaces arround the delimiter
		"1": {
			[]string{"version = 1.0"},
			[]string{"version=2.0"},
			[]string{},
			"=",
			false, false, false,
			[]string{"version=2.0"},
		},
		// toggle, remove old 'role', 'type'
		"2": {
			[]string{"version = 1.0", "role = KV", "type = disc"},
			[]string{"aaa"},
			[]string{"role", "type"},
			"=",
			false, false, false,
			[]string{"version = 1.0", "aaa"},
		},
		// toggle : add new 'role' and 'type'
		"3": {
			[]string{"version = 1.0", "role = KV", "type = disc", "aaa"},
			[]string{"role=master,type=ram"},
			[]string{"aaa"},
			"=",
			false, false, false,
			[]string{"version = 1.0", "role=master", "type=ram"},
		},
		// rm and add, with clear flag at first
		"4": {
			[]string{"version=1.0", "b"},
			[]string{"a,b,c"},
			[]string{"c,z"},
			"=",
			true, false, false,
			[]string{"a", "b", "c"},
		},
	}
	return func() map[string]TC {
		return ETESTS
	}
}

func diff(X, Y []string) []string {

	diff := []string{}
	vals := map[string]struct{}{}

	for _, x := range X {
		vals[x] = struct{}{}
	}

	for _, x := range Y {
		if _, ok := vals[x]; !ok {
			diff = append(diff, x)
		}
	}

	return diff
}

func TestTags(t *testing.T) {
	for k, tc := range tagsModifyTestCases()() {
		tags := consul.ModifyTags(tc.Tags, tc.Add, tc.Remove, tc.Delimitor, tc.hasClear, tc.isPlainMode, tc.isStringMode)
		result1 := diff(tags, tc.ExpectTags)
		if len(result1) == 0 {
			fmt.Printf(" -> sub-test '%s' done.\n", k)
		} else {
			fmt.Printf(" -> sub-test '%s' failed: expect %v, but got %v. diff = %v\n", k, tc.ExpectTags, tags, result1)
			t.Logf(" -> sub-test '%s' failed: expect %v, but got %v. diff = %v\n", k, tc.ExpectTags, tags, result1)
			t.Fail()
		}
	}
	// t.Logf("TestConsulConnection() return ERROR: %v", nil)
	// t.Fail()
}
