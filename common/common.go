package common

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/url"

	"github.com/miekg/dns"
)

// GenerateRandomIPAddress 生成随机IP地址
func GenerateRandomIPAddress() string {
	// 创建一个长度为4的byte切片
	ip := make([]byte, 4)
	// 从随机源中读取随机数，并将其写入ip切片中
	rand.Read(ip)
	// 将ip切片中的4个byte转换为点分十进制表示的IP地址
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// ReadableBytes 人性化显示字节大小
func ReadableBytes(bytes float64) (expression string) {
	// 如果字节数为0，则返回"0B"
	if bytes == 0 {
		return "0B"
	}
	// 计算字节数的对数，以1024为底
	var i = math.Floor(math.Log(bytes) / math.Log(1024))
	// 定义字节数单位
	var sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	// 格式化输出字节数和单位
	return fmt.Sprintf("%.3f%s", bytes/math.Pow(1024, i), sizes[int(i)])
}

// Nslookup 查询DNS记录
func Nslookup(targetAddress, server string, queryIpv6 bool) (res []string) {
	// 默认使用谷歌DNS
	if server == "" {
		server = "8.8.8.8"
	}
	// 创建一个DNS客户端
	client := dns.Client{}
	// 创建一个DNS请求
	request := dns.Msg{}
	if queryIpv6 {
		request.SetQuestion(targetAddress+".", dns.TypeAAAA)
	} else {
		request.SetQuestion(targetAddress+".", dns.TypeA)
	}
	// 拼接DNS服务器地址
	ns := server + ":53"
	// 发送DNS请求
	response, _, err := client.Exchange(&request, ns)
	if err != nil {
		fmt.Printf("nameserver %s error: %v\n", ns, err)
		return res
	}
	// 遍历DNS响应
	for _, answer := range response.Answer {
		if queryIpv6 && answer.Header().Rrtype == dns.TypeAAAA {
			AAAArecord := answer.(*dns.AAAA)
			res = append(res, fmt.Sprintf("%s", AAAArecord.AAAA.String()))
		} else if !queryIpv6 && answer.Header().Rrtype == dns.TypeA {
			Arecord := answer.(*dns.A)
			res = append(res, fmt.Sprintf("%s", Arecord.A.String()))
		}
	}
	return res
}

// RandStringBytesMaskImpr 生成指定长度的随机字符串
func RandStringBytesMaskImpr(n int) string {
	// 定义随机字符串的字符集
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 创建一个长度为n的byte切片
	b := make([]byte, n)
	// 循环n次，每次从字符集中随机选择一个字符，并将其写入byte切片中
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	// 将byte切片转换为字符串并返回
	return string(b)
}

// SumFloat64Array 把一个 []float64 循环相加, 返回相加后的值
func SumFloat64Array(arr []float64) (sum float64) {
	for _, v := range arr {
		sum += v
	}
	return sum
}

// IsIP 判断是否为IP地址 (IPv4/IPv6)
func IsIP(urlString string) (bool, int, string) {
	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return false, 0, ""
	}
	host := u.Hostname()
	ip := net.ParseIP(host)
	if ip == nil {
		return false, 0, host
	}
	if ip.To4() != nil {
		return true, 4, host
	}
	return true, 6, host
}
