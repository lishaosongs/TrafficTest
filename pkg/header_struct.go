package pkg

import (
	"fmt"
	"strings"
)

type Header struct {
	key, value string
}

type HeadersList []Header

// String 方法将 headersList 转换为字符串
func (h *HeadersList) String() string {
	return fmt.Sprint(*h)
}

// IsCumulative 方法返回 true，表示可以累加
func (h *HeadersList) IsCumulative() bool {
	return true
}

// Set 方法将字符串转换为 header 并添加到 headersList 中
func (h *HeadersList) Set(value string) error {
	// 将字符串按冒号分割为两部分
	res := strings.SplitN(value, ":", 2)
	if len(res) != 2 {
		return nil
	}
	// 创建 header 并添加到 headersList 中
	*h = append(*h, Header{
		res[0], strings.Trim(res[1], " "),
	})
	return nil
}
