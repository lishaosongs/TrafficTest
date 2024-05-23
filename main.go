package main

import (
	"TrafficTest/common"
	"TrafficTest/pkg"
	"fmt"
	"net"
	URL "net/url"
	"os"
	"sync"
	"time"

	"github.com/apoorvam/goterminal"
	CliV2 "github.com/urfave/cli/v2"
)

// 版本和编译时间
var version = ""
var date = ""

var (
	thisTime       int64
	terminalWriter *goterminal.Writer
	customIP       pkg.IpArray
	headers        pkg.HeadersList
)

func init() {
	thisTime = time.Now().Unix()
	terminalWriter = goterminal.New(os.Stdout)
}

// 生成一个文字logo TrafficTest
var cmdLogo = `
  _____           __  __ _      
 |_   _| __ __ _ / _|/ _(_) ___ 
   | || '__/ _| | |_| |_| |/ __|
   | || | | (_| |  _|  _| | (__ 
   |_||_|  \__,_|_| |_| |_|\___|
  _____         _               
 |_   _|__  ___| |_             
   | |/ _ \/ __| __|            
   | |  __/\__ \ |_             
   |_|\___||___/\__|
`

func main() {
	app := &CliV2.App{
		Name:        cmdLogo + "\n" + "TrafficTest",
		Usage:       "流量测试工具",
		UsageText:   "TrafficTest [命令选项] [参数...]",
		Version:     version,
		Description: "编译时间" + date + "\n" + "版本号" + version,
		HideVersion: true,
		Flags: []CliV2.Flag{
			&CliV2.IntFlag{
				Name:    "thread",
				Aliases: []string{"t"},
				Value:   16,
				Usage:   "下载时的并发线程数",
			},
			&CliV2.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Value:   "https://speed.cloudflare.com/__down?bytes=500000000",
				Usage:   "下载的目标URL",
			},
			&CliV2.StringFlag{
				Name:    "post",
				Aliases: []string{"p"},
				Value:   "",
				Usage:   "POST请求的内容",
			},
			&CliV2.BoolFlag{
				Name:    "get",
				Aliases: []string{"g"},
				Value:   false,
				Usage:   "是否在URL中添加随机参数",
			},
			&CliV2.StringFlag{
				Name:    "referer",
				Aliases: []string{"r"},
				Value:   "",
				Usage:   "HTTP请求头中的Referer字段",
			},
			&CliV2.BoolFlag{
				Name:    "forwarded",
				Aliases: []string{"f"},
				Value:   true,
				Usage:   "是否随机生成X-Forwarded-For和X-Real-IP字段",
			},
			&CliV2.StringSliceFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				Usage:   "自定义域名的IP地址,多个地址将自动随机分配",
			},
			&CliV2.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "自定义header请求头",
			},
			&CliV2.StringFlag{
				Name:  "ua",
				Value: "",
				Usage: "自定义User-Agent,填写完整ua或pc或mobile生成对应平台的UA,留空则随机生成",
			},
			&CliV2.IntFlag{
				Name:  "flow",
				Value: 0,
				Usage: "总流量限制(GB),默认不限制,超过流量后程序将自动退出.",
			},
			&CliV2.BoolFlag{
				Name:  "6",
				Value: false,
				Usage: "是否使用IPv6地址",
			},
			&CliV2.StringFlag{
				Name:  "dns",
				Value: "8.8.8.8",
				Usage: "指定DNS服务器地址",
			},
			&CliV2.IntFlag{
				Name:  "time",
				Value: 0,
				Usage: "执行时间限制(秒),默认不限制,超过时间后程序将自动退出.",
			},
			&CliV2.StringFlag{
				Name:  "socks5",
				Value: "",
				Usage: "使用socks5代理,格式:127.0.0.1:1080",
			},
		},
		Action: func(c *CliV2.Context) error {
			if c.Bool("help") {
				_ = CliV2.ShowAppHelp(c)
				return nil
			}
			// 检查URL是否是IP地址
			isIp, _, ipAddress := common.IsIP(c.String("url"))
			if isIp {
				customIP = pkg.IpArray{ipAddress}
			}

			//解析IP地址
			if len(customIP) <= 0 {
				u, _ := URL.Parse(c.String("url"))
				ipArr := common.Nslookup(u.Hostname(), c.String("dns"), c.Bool("6"))
				if len(ipArr) == 0 {
					fmt.Println("您提交的URL,未能从DNS中获取到 IP")
					os.Exit(0)
				}
				customIP = ipArr
			} else {
				for _, ip := range customIP {
					ParseIP := net.ParseIP(ip)
					if ParseIP == nil {
						fmt.Println("您提交的IP,格式错误.", ip)
						os.Exit(0)
					}
				}
			}

			routines := c.Int("thread")
			// 显示统计信息
			go pkg.ShowStat(customIP, c.String("url"), terminalWriter, c.Int("flow"), thisTime, int64(c.Int("time")))
			pkg.HandleSignals(thisTime, terminalWriter)

			var waitGroup sync.WaitGroup

			if routines <= 0 {
				routines = 16
			}

			for i := 0; i < routines; i++ {
				waitGroup.Add(1)
				go pkg.GoFun(
					c.String("url"),
					c.String("post"),
					c.Bool("get"),
					c.String("referer"),
					c.Bool("forwarded"),
					c.String("ua"),
					c.String("socks5"),
					customIP,
					headers,
					&waitGroup,
				)
			}
			waitGroup.Wait()
			return nil
		},
	}
	app.Before = func(c *CliV2.Context) error {
		customIP = c.StringSlice("ip")
		// 将[]string转换为pkg.HeadersList
		for _, headerStr := range c.StringSlice("header") {
			_ = headers.Set(headerStr)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
