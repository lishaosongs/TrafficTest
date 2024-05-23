package pkg

import (
	"TrafficTest/common"
	"fmt"
	"github.com/apoorvam/goterminal"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// HandleSignals 处理信号
func HandleSignals(thisTime int64, TerminalWriter *goterminal.Writer) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		// 重置终端输出
		statistics(thisTime, TerminalWriter)
	}()
}

// statistics 统计
func statistics(thisTime int64, TerminalWriter *goterminal.Writer) {
	TerminalWriter.Reset()
	fmt.Println("正在退出程序...")
	// 其他必要的处理 预留
	fmt.Println("-------------------------------")
	runTime := time.Now().Unix() - thisTime
	avgTotalRecvs := TotalRecvs / float64(runTime)
	avgTotalSents := TotalSents / float64(runTime)
	fmt.Println("运行时长:" + strconv.FormatInt(runTime, 10) + "秒")
	fmt.Println("程序执行,共接收:"+common.ReadableBytes(TotalRecvs), "平均速度:"+common.ReadableBytes(avgTotalRecvs)+"/s")
	fmt.Println("程序执行,共发送:"+common.ReadableBytes(TotalSents), "平均速度:"+common.ReadableBytes(avgTotalSents)+"/s")
	fmt.Println("-------------------------------")
	os.Exit(0)
}
