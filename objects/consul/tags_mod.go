package consul

import (
	log "github.com/cihub/seelog"
	"github.com/hedzr/consul-tags/util"
	"strings"
)

func ModifyTags(tags, addTags, removeTags []string, delim string, hasClear, isPlainMode, isString bool) []string {
	if hasClear {
		tags = make([]string, 0)
	}

	if !isPlainMode {
		log.Debug("    --- in ext mode")
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
			log.Debugf("    --- slice: erasing %s", t)
			for {
				erased := false
				for i, v := range tags {
					va := strings.Split(v, delim)
					ta := strings.Split(t, delim)
					if len(ta) > 0 && strings.EqualFold(va[0], ta[0]) {
						tags = util.SliceEraseByIndex(tags, i)
						log.Debugf("      - slice: erased '%s%s...'", va[0], delim)
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
			log.Debugf("    --- slice: appending %s", t)
			matched := false
			for i, v := range tags {
				va := strings.Split(v, delim)
				ta := strings.Split(t, delim)
				if len(va) > 0 && strings.EqualFold(strings.TrimSpace(va[0]), strings.TrimSpace(ta[0])) {
					tags[i] = t
					log.Debugf("      - slice: appended '%s%s...'", va[0], delim)
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
