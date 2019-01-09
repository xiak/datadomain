package parser

import (
	"regexp"
	"fmt"
)

func DDBoostStorageUnitShow(src string) (int, []string) {
	r := regexp.MustCompile(`(?P<product>avamar-)(?P<id>\d+)\s+`)
	matches := r.FindAllStringSubmatch(src, -1)
	length := len(matches)
	storages := make([]string, length)
	if matches == nil {
		fmt.Println("Parser: Not matche items")
		return length, storages
	}
	for i, m := range matches {
		storages[i] = fmt.Sprintf("%s%s", m[1], m[2])
	}
	return length, storages
}
