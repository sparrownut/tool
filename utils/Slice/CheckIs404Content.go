package Slice

import (
	"tool/Global"
)

func CheckIs404Content(content string) bool {
	for _, k := range Global.WebNoContentLib {
		if content == k {
			return true
		}
	}
	return false
}
