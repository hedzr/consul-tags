/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"strings"

	"github.com/hedzr/cmdr/v2/pkg/logz"

	"github.com/hedzr/consul-tags/util"
)

func ModifyTags(tags, addTags, removeTags []string, delim string, hasClear, isPlainMode, isStringMode bool) []string {
	if hasClear {
		tags = make([]string, 0)
	}

	if !isPlainMode {
		logz.Debug("    --- in ext mode")
		list := make([]string, 0)
		for _, t := range removeTags {
			if isStringMode {
				if t != "" {
					list = append(list, t)
				}
			} else {
				for _, t1 := range strings.Split(t, ",") {
					if t1 != "" {
						list = append(list, t1)
					}
				}
			}
		}
		for _, t := range list {
			logz.Debug("    --- slice: erasing %s", t)
			for {
				erased := false
				for i, v := range tags {
					va := strings.Split(v, delim)
					ta := strings.Split(t, delim)
					if len(ta) > 0 && strings.EqualFold(va[0], ta[0]) {
						tags = util.SliceEraseByIndex(tags, i)
						logz.Debug("      - slice: tags erased", "0", va[0], "delim", delim)
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
			if isStringMode {
				if t != "" {
					list = append(list, t)
				}
			} else {
				for _, t1 := range strings.Split(t, ",") {
					if t1 != "" {
						list = append(list, t1)
					}
				}
			}
		}
		for _, t := range list {
			logz.Debug("    --- slice: appending tags", "tag", t)
			matched := false
			for i, v := range tags {
				va := strings.Split(v, delim)
				ta := strings.Split(t, delim)
				if len(va) > 0 && strings.EqualFold(strings.TrimSpace(va[0]), strings.TrimSpace(ta[0])) {
					tags[i] = t
					logz.Debug("      - slice: tags appended", "t", va[0], "delim", delim)
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
		logz.Info("      - slice: cleared.")
	}

	if !isPlainMode {
		logz.Debug("    --- in ext mode")
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
			logz.Debug("    --- slice: erasing tags", "tag", t)
			for {
				ta := strings.Split(t, delim)
				if ta[0] != "" {
					delete(tags, ta[0])
				}
				logz.Debug("      - slice: tags erased '%s%s...'", "0", ta[0], "delim", delim)
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
			logz.Debug("    --- slice: appending tags", "tag", t)
			ta := strings.Split(t, delim)
			tak := ta[0]
			tav := strings.Join(ta[1:], delim)
			tags[tak] = tav
			logz.Debug("      - slice: tags set/appended", "0", tak, "delim", delim)
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
