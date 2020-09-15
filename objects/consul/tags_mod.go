/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags/util"
	"strings"
)

func ModifyTags(tags, addTags, removeTags []string, delim string, hasClear, isPlainMode, isString bool) []string {
	if hasClear {
		tags = make([]string, 0)
	}

	if !isPlainMode {
		cmdr.Logger.Debugf("    --- in ext mode")
		list := make([]string, 0)
		for _, t := range removeTags {
			if isString {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			cmdr.Logger.Debugf("    --- slice: erasing %s", t)
			for {
				erased := false
				for i, v := range tags {
					va := strings.Split(v, delim)
					ta := strings.Split(t, delim)
					if len(ta) > 0 && strings.EqualFold(va[0], ta[0]) {
						tags = util.SliceEraseByIndex(tags, i)
						cmdr.Logger.Debugf("      - slice: erased '%s%s...'", va[0], delim)
						erased = true
						break
					}
				}
				if !erased {
					break
				}
			}
		}
		list = make([]string, 0)
		for _, t := range addTags {
			if isString {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			cmdr.Logger.Debugf("    --- slice: appending %s", t)
			matched := false
			for i, v := range tags {
				va := strings.Split(v, delim)
				ta := strings.Split(t, delim)
				if len(va) > 0 && strings.EqualFold(strings.TrimSpace(va[0]), strings.TrimSpace(ta[0])) {
					tags[i] = t
					cmdr.Logger.Debugf("      - slice: appended '%s%s...'", va[0], delim)
					matched = true
				}
			}
			if !matched {
				tags = append(tags, t)
			}
		}

	} else {
		for _, t := range removeTags {
			tags = util.SliceErase(tags, t)
		}
		for _, t := range addTags {
			tags = append(tags, t)
		}
	}

	return tags
}

func ModifyNodeMeta(tags map[string]string, addTags, removeTags []string, delim string, hasClear, isPlainMode, isString bool) map[string]string {
	if hasClear {
		tags = make(map[string]string, 0)
		cmdr.Logger.Infof("      - slice: cleared.")
	}

	if !isPlainMode {
		cmdr.Logger.Debugf("    --- in ext mode")
		list := make([]string, 0)
		for _, t := range removeTags {
			if isString {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			cmdr.Logger.Debugf("    --- slice: erasing %s", t)
			for {
				ta := strings.Split(t, delim)
				if ta[0] != "" {
					delete(tags, ta[0])
				}
				cmdr.Logger.Infof("      - slice: erased '%s%s...'", ta[0], delim)
			}
		}
		list = make([]string, 0)
		for _, t := range addTags {
			if isString {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			cmdr.Logger.Debugf("    --- slice: appending %s", t)
			ta := strings.Split(t, delim)
			tak := ta[0]
			tav := strings.Join(ta[1:], delim)
			tags[tak] = tav
			cmdr.Logger.Debugf("      - slice: set/appended '%s%s...'", tak, delim)
		}

	} else {
		for _, t := range removeTags {
			delete(tags, t)
		}
		for _, t := range addTags {
			tags[t] = ""
		}
	}

	return tags
}
