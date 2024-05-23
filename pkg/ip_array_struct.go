package pkg

import "strings"

type IpArray []string

// Set 实现 flag.Value 接口的 Set 方法
func (i *IpArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// 实现 flag.Value 接口的 String 方法
func (i *IpArray) String() string {
	return strings.Join(*i, " | ")
}
