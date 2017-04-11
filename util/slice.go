package util

import "strings"

func SlaceContains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func SliceErase(a []string, t string) []string {
	i := SliceIndex(len(a), func(i int) bool {
		return strings.EqualFold(a[i], t)
	})
	return SliceEraseByIndex(a, i)
}

func SliceEraseByIndex(a []string, index int) []string {
	if index > -1 {
		copy(a[index:], a[index+1:])
		return a[:len(a)-1]
	}
	return a
}
