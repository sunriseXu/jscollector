package main

import "fmt"

func main() {
	ips := []string{
		// "https://www.baidu.com",
		"https://juejin.cn",
		"https://mermaid.live/edit",
	}

	urls := CheckSiteAvailable(ips)

	fmt.Println("found valid urls len:", len(urls))
	// fmt.Println(urls)

	// urls = []string{
	// 	"https://juejin.cn/post/7241096652919193658",
	// }

	// api 爆破
	results := LoadExtensionWrapper(urls)

	WriteResults(results, "test.json")

}
