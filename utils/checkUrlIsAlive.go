package utils

import (
	"fmt"
	"github.com/zhshch2002/goreq"
	"time"
)

func CheckUrlIsAlive(url string) bool {
	defer func() {
		DBGLOG("check done")
	}()
	DBGLOG(fmt.Sprintf("checking %v", url))
	//ch := make(chan string)
	getreq := goreq.Get(url).SetClient(goreq.NewClient()).SetTimeout(time.Duration(time.Second * 10))
	resp := getreq.Do()
	if resp.Error() == nil {
		return true
	}
	return false
}
