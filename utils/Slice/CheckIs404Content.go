package Slice

import (
	"strings"
	"tool/Global"
)

func CheckIs404Content(content string) bool {
	for _, k := range Global.WebNoContentLib {
		if strings.Contains(content, k) {
			return true
		}
	}
	return false
}
