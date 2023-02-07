package core

import (
	"fmt"
	"github.com/antlabs/strsim"
	"github.com/zhshch2002/goreq"
	"os"
	"strconv"
	"strings"
	"time"
	"tool/Global"
	"tool/utils"
	"tool/utils/Slice"
)

type WebStatus struct {
	isFail     bool //填充随机url(一定不存在) 查看返回指纹和新指纹有什么区别
	url        string
	dir        string
	statusCode int
	body       []byte
}

type Task struct {
	url string // url要后边有/ 的
	dir string // FileExplodeDict里的迭代
}
type TaskList struct {
	tasks      []Task
	maxthreads int
}

func judgFingerPrintIsSame(fp1 WebStatus, fp2 WebStatus) float64 {
	//判断两个网页返回的指纹是否相等 10分满分为完全相等 0分为完全不同
	//println(fp1.statusCode)
	//println(fp2.statusCode)
	score := 0.0
	if fp1.statusCode == fp2.statusCode {
		score += 5.0
	}
	score += 5.0 * strsim.Compare(string(fp1.body), string(fp2.body))
	return score
}

func MakeTasks(isSuc *bool) TaskList {
	// 制造任务列表 下一步交给任务执行器多线程执行
	file, readFileErr := os.ReadFile(Global.INPUTFILE)
	taskListObject := TaskList{maxthreads: 512}
	if readFileErr != nil {
		utils.Printerr("文件输入错误")
		os.Exit(0)
		return TaskList{}
	}
	for _, url := range strings.Split(string(file), "\n") {
		if len(url) == 0 {
			continue
		}
		url = utils.UrlFixer(url)
		for _, dir := range Global.Top2000FilesList {
			if len(dir) == 0 {
				continue
			}
			//println(url)

			utils.StringsFilter(&dir)

			taskListObject.tasks = append(taskListObject.tasks, Task{url: url, dir: dir})
		}
	}
	if len(taskListObject.tasks) > 0 {
		*isSuc = true
	}
	return taskListObject
}

func DoTasks(tasklist TaskList, isSuc *bool) {
	defer func() {
		if r := recover(); r != nil {
			if Global.DBG {
				println(r) // DBG模式下报个错 这实际上意味着网站无法访问
			}
		}
	}()

	bannedList := []string{} // 访问失败的，以后禁止访问的
	doneMap := make(map[string]WebStatus)
	maxThreads := tasklist.maxthreads
	threads := 0
	for _, task := range tasklist.tasks {
		//utils.DBGLOG(fmt.Sprintf("%v", task.url+task.dir))
		//检测此次任务是否被标记黑名单 (无法连通)
		if Slice.CheckIsStringInSlice(task.url, bannedList) {
			//utils.DBGLOG("在黑名单中")
			continue
		}
		//随机访问url 获取失败的返回指纹
		if _, k := doneMap[task.url]; !k { // 如果没扫过这个url
			//扫描是否能访问 不能 这个url拉黑
			if !utils.CheckUrlIsAlive(task.url) {
				utils.Printerr(fmt.Sprintf("%v不能访问 拉黑", task.url))
				utils.DBGLOG(fmt.Sprintf("拉黑%v", task.url))
				bannedList = append(bannedList, task.url)
				continue
			}

			//扫描404指纹
			utils.DBGLOG(fmt.Sprintf("测试404指纹%v", task.url))
			FailFingerPrint := WebStatus{isFail: true}
			getFailreq := goreq.Get(task.url + utils.RandStringRunes(8)).SetClient(goreq.NewClient()).SetTimeout(time.Duration(time.Second * 20))
			respFail := getFailreq.Do()
			if respFail.Error() != nil {
				bannedList = append(bannedList, task.url)
			}
			FailFingerPrint.url = task.url
			FailFingerPrint.dir = task.dir
			FailFingerPrint.statusCode = respFail.StatusCode
			doneMap[task.url] = FailFingerPrint
		}

		//执行真正的扫描
		if task.url == "" || task.url == "/" || task.dir == "" {
			//去除空资产
			continue
		}
		//重试机制
		retryNMap := make(map[string]int)
		retryNMap[task.url] = 0
		//执行
	retry:
		if threads <= maxThreads {
			time.Sleep(time.Duration(time.Millisecond * 10))
			go func(task Task) {
			goroutineRetry:
				defer func() {
					if r := recover(); r != nil {
						if Global.DBG {
							println(r) // DBG模式下报个错 这实际上意味着网站无法访问
						}
					}
					threads--
				}()
				threads++
				//开始执行内容
				TmpFingerPrint := WebStatus{isFail: false}
				taskResolved := task.url + task.dir
				getreq := goreq.Get(taskResolved).SetClient(goreq.NewClient())
				resp := getreq.Do()
				if resp.Error() != nil {
					if retryNMap[task.url] <= 5 {
						retryNMap[task.url]++
						goto goroutineRetry
					}
					return
				}
				//获取指纹
				TmpFingerPrint.dir = task.dir
				TmpFingerPrint.url = task.url
				TmpFingerPrint.statusCode = resp.StatusCode
				TmpFingerPrint.body = resp.Body

				//simScore := judgFingerPrintIsSame(FailFingerPrint, TmpFingerPrint)
				//println(doneMap[task.url].statusCode)
				if doneMap[task.url].statusCode != TmpFingerPrint.statusCode { //如果存在漏洞 (与随机字符串的地址相差很大) || doneMap[task.url].body != TmpFingerPrint.body
					if !Slice.CheckIs404Content(string(TmpFingerPrint.body)) { //没有特征迹象
						utils.Printsuc(fmt.Sprintf("URL{%v} RESP_LEN{%v} RESP_CODE{%v}", task.url+task.dir, len(resp.Body), resp.StatusCode))
					}

				} else if Global.DBG {
					println("threads:" + strconv.Itoa(threads))
					utils.Printminfo(fmt.Sprintf("URL{%v} RESP_CODE{%v} 扫描完成 无敏感泄露", task.url+task.dir, resp.StatusCode))

				}
			}(task)
		} else {
			//如果线程超过最大就重试
			time.Sleep(time.Duration(time.Millisecond))
			goto retry
		}
	}
	*isSuc = true
}

func ScanDirsDo(isSuc *bool) {
	// 扫描敏感目录
	isMakeTasksSuc := false
	isDoTasksSuc := false
	tasklist := MakeTasks(&isMakeTasksSuc)
	if !isMakeTasksSuc {
		*isSuc = false
		return
	}
	// 执行任务
	DoTasks(tasklist, &isDoTasksSuc)
	if !isDoTasksSuc {
		*isSuc = false
		return
	}
	*isSuc = true

}
