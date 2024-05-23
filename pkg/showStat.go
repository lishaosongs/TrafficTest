package pkg

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	netstat "github.com/shirou/gopsutil/net"
	"time"

	"TrafficTest/common"
)

// ShowStat 在控制台中显示统计信息
func ShowStat(customIP IpArray, url string, TerminalWriter *goterminal.Writer, flow int, thisTime int64, endTime int64) {
	// 获取初始网络计数器
	initialNetCounter, _ := netstat.IOCounters(true)
	//总接收和发送
	TotalRecv := make([]float64, len(initialNetCounter))
	TotalSent := make([]float64, len(initialNetCounter))
	// 获取 IP 列表
	var iplist string
	if len(customIP) > 0 {
		iplist = customIP.String()
	} else {
		iplist = "未知的IP"
	}
	for {
		// 获取 CPU 占用率
		percent, _ := cpu.Percent(time.Second, false)
		percentTotal := common.SumFloat64Array(percent)

		// 获取内存使用情况
		memStat, _ := mem.VirtualMemory()

		// 获取网络计数器
		netCounter, _ := netstat.IOCounters(true)

		// 获取系统负载情况
		loadStat, _ := load.Avg()
		//var memStats runtime.MemStats
		//runtime.ReadMemStats(&memStats)
		//// 获取当前程序占用的内存大小
		//memUsage := memStats.Alloc
		// 输出统计信息
		_, _ = fmt.Fprintf(TerminalWriter, "URL: %s\n", url)
		_, _ = fmt.Fprintf(TerminalWriter, "IP: %s\n", iplist)
		_, _ = fmt.Fprintf(TerminalWriter, "CPU: %.3f%% \n", percentTotal)
		_, _ = fmt.Fprintf(TerminalWriter, "内存: %.3f%% \n", memStat.UsedPercent)
		//_, _ = fmt.Fprintf(TerminalWriter, "程序占用内存: %s \n", common.ReadableBytes(float64(memUsage)))
		_, _ = fmt.Fprintf(TerminalWriter, "负载: %.3f %.3f %.3f\n", loadStat.Load1, loadStat.Load5, loadStat.Load15)

		// 遍历网络计数器
		for i := 0; i < len(netCounter); i++ {
			// 如果接收和发送的字节数都为 0，则跳过
			if netCounter[i].BytesRecv == 0 && netCounter[i].BytesSent == 0 {
				continue
			}
			// 计算接收和发送的字节数
			RecvBytes := float64(netCounter[i].BytesRecv - initialNetCounter[i].BytesRecv)
			TotalRecv[i] = TotalRecv[i] + RecvBytes
			SendBytes := float64(netCounter[i].BytesSent - initialNetCounter[i].BytesSent)
			TotalSent[i] = TotalSent[i] + SendBytes

			// 输出网络统计信息
			_, _ = fmt.Fprintf(TerminalWriter, "网卡: %v | 接收: %s(%s/s) | 发送: %s(%s/s)\n", netCounter[i].Name,
				common.ReadableBytes(TotalRecv[i]),
				common.ReadableBytes(RecvBytes),
				common.ReadableBytes(TotalSent[i]),
				common.ReadableBytes(SendBytes))
		}
		_, _ = fmt.Fprintf(TerminalWriter, "请求次数: %d\n", TotalRequests)
		for i := 0; i < len(TotalRecv); i++ {
			TotalRecvs = TotalRecvs + float64(netCounter[i].BytesRecv-initialNetCounter[i].BytesRecv)
		}
		for i := 0; i < len(TotalSent); i++ {
			TotalSents = TotalSents + float64(netCounter[i].BytesSent-initialNetCounter[i].BytesSent)
		}
		// 更新初始网络计数器
		initialNetCounter = netCounter
		// 清空终端输出
		TerminalWriter.Clear()

		// 输出统计信息
		_ = TerminalWriter.Print()
		//判断是否达到流量限制
		if flow > 0 {
			if TotalRecvs > float64(flow*1024*1024*1024) {
				fmt.Println("已达到流量限制...")
				statistics(thisTime, TerminalWriter)
			}
		}
		if endTime != 0 {
			if time.Now().Unix() >= thisTime+endTime {
				fmt.Println("已达到时间限制...")
				statistics(thisTime, TerminalWriter)
			}
		}
		// 等待 50 毫秒
		time.Sleep(50 * time.Millisecond)
	}
}
