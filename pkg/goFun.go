package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"golang.org/x/net/proxy"
	"io"
	"math/rand"
	"net"
	"net/http"
	URL "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"TrafficTest/common"
)

// GoFun 执行任务
func GoFun(ReqUrl string, postContent string, getContent bool, Referer string, XForwardFor bool, userAgent string, socks5 string, customIP IpArray, headers HeadersList, wg *sync.WaitGroup) {

	// 异常 重新执行
	defer func() {
		if r := recover(); r != nil {
			go GoFun(ReqUrl, postContent, getContent, Referer, XForwardFor, userAgent, socks5, customIP, headers, wg)
		}
	}()
	// 创建一个http客户端
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	if socks5 != "" {
		// 创建一个SOCKS5拨号器
		socks5Dialer, err := proxy.SOCKS5("tcp", socks5, nil, proxy.Direct)
		if err != nil {
			fmt.Println("无法连接到代理")
			os.Exit(0)
		}
		// 尝试通过代理连接到一个地址
		conn, err := socks5Dialer.Dial("tcp", "www.google.com:80")
		if err != nil {
			fmt.Println("无法通过代理拨号:", "请检查代理是否正常工作")
			os.Exit(0)
		}
		_ = conn.Close()
		// 如果有自定义IP，则创建一个自定义的transport
		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return socks5Dialer.Dial(network, addr)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = transport
	} else {
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		ip := ""
		ipPort := "80"
		ipPortSSL := "443"
		//判断URL是否携带了端口
		parsedURL, err := URL.Parse(ReqUrl)
		if err == nil {
			if parsedURL.Port() != "" {
				ipPort = parsedURL.Port()
				ipPortSSL = parsedURL.Port()
			}
		}
		// 如果有自定义IP，则创建一个自定义的transport
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				ip = customIP[rand.Intn(len(customIP))]
				ParseIP := net.ParseIP(ip)
				if ParseIP == nil {
					// IP 地址格式错误
					fmt.Println("您提交的IP,格式错误.", ip)
					os.Exit(0)
				}
				if ParseIP.To4() != nil {
					// IPv4 地址
					addr = ip + ":" + ipPort
				} else {
					// IPv6 地址
					addr = "[" + ip + "]:" + ipPort
				}
				return dialer.DialContext(ctx, network, addr)
			},
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				ip = customIP[rand.Intn(len(customIP))]
				ParseIP := net.ParseIP(ip)
				if ParseIP == nil {
					// IP 地址格式错误
					fmt.Println("您提交的IP,格式错误.", ip)
					os.Exit(0)
				}
				if ParseIP.To4() != nil {
					// IPv4 地址
					addr = ip + ":" + ipPortSSL
				} else {
					// IPv6 地址
					addr = "[" + ip + "]:" + ipPortSSL
				}
				return tls.DialWithDialer(dialer, network, addr, &tls.Config{
					InsecureSkipVerify: true, // 仅在测试环境中使用，忽略证书验证
				})
			},
		}
		client.Transport = transport
	}
	// 循环发送请求
	for {
		var request *http.Request
		var err1 error = nil
		// 根据postContent是否为空，创建不同的请求
		if len(postContent) > 0 {
			request, err1 = http.NewRequest("POST", ReqUrl, strings.NewReader(postContent))
		} else {
			request, err1 = http.NewRequest("GET", ReqUrl, nil)
		}
		if err1 != nil {
			continue
		}
		if getContent {
			//Query参数
			params := URL.Values{}
			params.Add(common.RandStringBytesMaskImpr(6), common.RandStringBytesMaskImpr(6))
			if request.URL.RawQuery == "" {
				request.URL.RawQuery = params.Encode()
			} else {
				request.URL.RawQuery += "&" + params.Encode()
			}
			request.URL.Path = path.Join(request.URL.Path, "/")
		}
		// 设置请求头
		request.Header.Add("Cookie", common.RandStringBytesMaskImpr(6)+":"+common.RandStringBytesMaskImpr(6)) // 添加随机Cookie
		if userAgent == "pc" {
			request.Header.Add("User-Agent", browser.Chrome()) // 添加PC User-Agent
		} else if userAgent == "mobile" {
			request.Header.Add("User-Agent", browser.Mobile()) // 添加Mobile User-Agent
		} else if userAgent == "" {
			request.Header.Add("User-Agent", browser.Random()) // 添加随机User-Agent
		} else {
			request.Header.Add("User-Agent", userAgent) // 添加随机User-Agent
		}

		if len(Referer) == 0 {
			Referer = ReqUrl
		}
		request.Header.Add("Referer", Referer) // 添加Referer
		if XForwardFor {
			randomIp := common.GenerateRandomIPAddress()
			request.Header.Add("X-Forwarded-For", randomIp) // 添加随机X-Forwarded-For
			request.Header.Add("X-Real-IP", randomIp)       // 添加随机X-Real-IP
		}

		// 如果有自定义的headers，则添加到请求头中
		if len(headers) > 0 {
			for _, head := range headers {
				headKey := head.key
				headValue := head.value
				// 如果header的key以"Random"开头，则将value中的"Random"替换为随机字符串
				if strings.HasPrefix(head.key, "Random") {
					count, convErr := strconv.Atoi(strings.ReplaceAll(head.value, "Random", ""))
					if convErr == nil {
						headKey = common.RandStringBytesMaskImpr(count)
					}
				}
				// 如果header的value以"Random"开头，则将value中的"Random"替换为随机字符串
				if strings.HasPrefix(head.value, "Random") {
					count, convErr := strconv.Atoi(strings.ReplaceAll(head.value, "Random", ""))
					if convErr == nil {
						headValue = common.RandStringBytesMaskImpr(count)
					}
				}
				// 删除原有的header，并添加新的header
				request.Header.Set(headKey, headValue)
			}
		}

		// 发送请求
		resp, err2 := client.Do(request)
		if err2 != nil {
			continue
		}
		TotalRequests++
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
