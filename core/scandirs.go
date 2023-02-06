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
	text       string
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
	score += 5.0 * strsim.Compare(fp1.text, fp2.text)
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
	doneMap := make(map[string]WebStatus)
	maxThreads := tasklist.maxthreads
	threads := 0
	for _, task := range tasklist.tasks {
		//随机访问url 获取失败的返回指纹

		if _, k := doneMap[task.url]; !k { // 如果没扫过这个url

			FailFingerPrint := WebStatus{isFail: true}
			getFailreq := goreq.Get(task.url + utils.RandStringRunes(8)).SetClient(goreq.NewClient())
			respFail := getFailreq.Do()
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
	retry:
		if threads <= maxThreads {
			time.Sleep(time.Duration(time.Millisecond * 10))
			go func(task Task) {
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
				if getreq.Err != nil && resp.Err != nil {
					//错误处理机制
					*isSuc = false
					return
				}
				//获取指纹
				TmpFingerPrint.dir = task.dir
				TmpFingerPrint.url = task.url
				TmpFingerPrint.statusCode = resp.StatusCode
				TmpFingerPrint.text = resp.Text

				//simScore := judgFingerPrintIsSame(FailFingerPrint, TmpFingerPrint)
				//println(doneMap[task.url].statusCode)
				if doneMap[task.url].statusCode != TmpFingerPrint.statusCode { //如果存在漏洞 (与随机字符串的地址相差很大) || doneMap[task.url].text != TmpFingerPrint.text
					if Slice.CheckIs404Content(resp.Text) { //没有特征迹象
						utils.Printsuc(fmt.Sprintf("URL{%v} RESP_LEN{%v}", task.url+task.dir, len(resp.Text)))
					}

				} else if Global.DBG {
					println("threads:" + strconv.Itoa(threads))
					utils.Printminfo(fmt.Sprintf("URL{%v} CODE{%v} BODY{%v} 扫描完成 无敏感泄露", task.url+task.dir, resp.StatusCode, resp.Text))

				}
				//return
				//定时30s结束
				for {
					select {
					case <-time.After(time.Duration(time.Second * 20)):
						threads--
						return
					}
				}
				//return
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
