package utils

import (
	"fmt"
	"strings"
)

func StringsFilter(s *string) {
	*s = strings.ReplaceAll(*s, " ", "")
	*s = strings.ReplaceAll(*s, "\n", "")
	*s = strings.ReplaceAll(*s, "\t", "")
	*s = strings.ReplaceAll(*s, "\r", "")
}

func UrlFixer(u string) string {
	url := u // 缓冲
	StringsFilter(&url)
	//print(hex.EncodeToString([]byte(url)))
	//println(url)
	if url[len(url)-1:] != "/" {
		//如果最后一位不是/
		url = fmt.Sprintf("%v/", url)
	}
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		//如果缺少http头
		url = "https://" + url
	}

	//println(url)
	return url
}
