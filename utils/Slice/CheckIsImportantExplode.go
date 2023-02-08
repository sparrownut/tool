package Slice

import (
	"strings"
	"tool/Global"
)

func CheckIsImportantExplode(s string) bool { //如果输入body里包含了重大发现的列表中的body
	for _, k := range Global.ImportantExplodeDict {
		if strings.Contains(strings.ToUpper(s), strings.ToUpper(k)) {
			return true
		}
	}
	return false
}
