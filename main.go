package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"tool/Global"
	"tool/core"
	"tool/utils"
	"tool/utils/network"
)

func main() {

	app := &cli.App{
		Name:      "tool",
		Usage:     "如题 它就叫工具 \n仅供授权的渗透测试使用 请遵守法律!", // 这里写协议
		UsageText: "一个工具集合",
		Version:   "0.3.3",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "MOD", Aliases: []string{"M"}, Destination: &Global.MOD, Value: "scandirs", Usage: "模式选择 \n 文件泄露目录扫描 - scandirs", Required: false},
			&cli.StringFlag{Name: "OutputFile", Aliases: []string{"O"}, Destination: &Global.OUTPUTFILE, Value: "default format", Usage: "输出文件", Required: false},
			&cli.StringFlag{Name: "InputFile", Aliases: []string{"F"}, Destination: &Global.INPUTFILE, Value: "list", Usage: "扫描输入文件", Required: true},
			&cli.BoolFlag{Name: "DBG", Aliases: []string{"D"}, Destination: &Global.DBG, Value: false, Usage: "DBG MOD", Required: false},
			&cli.BoolFlag{Name: "Proxy", Aliases: []string{"P"}, Destination: &Global.PROXYOPEN, Value: false, Usage: "是否开启自动代理爬虫(防墙)", Required: false},

			//&cli.IntFlag{Name: "checkN", Aliases: []string{"C"}, Destination: &Global.CHECKN, Value: 3, Usage: "同一端口检测次数", Required: false},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			Init()
			do()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		//panic(err)
	}

	//fmt.Printf(os.Args[1])
}

func Init() {
	if Global.PROXYOPEN {
		go network.ProxyCrawlerInit()
	}
	//授权
	if !network.CheckDomainIsAlive("www.baidu.com") && !network.CheckDomainIsAlive("www.google.com") {
		println("网络连接故障")
		os.Exit(0)
	}
	if !network.CheckDomainIsAlive("toolkey.stuid-fish.co") {
		println("此软件未授权!")
		os.Exit(0)
	}
	utils.Printminfo("授权成功")
	//初始化
	utils.Printminfo("初始化成功")
}
func do() {
	scanDirisSuc := false
	core.ScanDirsDo(&scanDirisSuc)
}
