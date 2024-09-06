package markdown

import "fmt"

func baiduBaike(target string) string {
	return fmt.Sprintf("https://baike.baidu.com/item/%s", target)
}

func wikipediaCn(target string) string {
	return fmt.Sprintf("https://zh.wikipedia.org/zh-cn/%s", target)
}
