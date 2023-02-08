package Slice

import (
	"strings"
	"tool/Global"
)

func CheckIs404Content(content string) bool {
	for _, k := range Global.WebNoContentLib {
		if strings.Contains(strings.ToUpper(content), strings.ToUpper(k)) {
			return true
		}
	}
	return false
}
