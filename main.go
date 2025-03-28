package main

import "fmt"

func main() {
	ips := []string{
		// "https://www.baidu.com",
		"https://juejin.cn",
		"https://mermaid.live/edit",
	}

	// get extensoin url
	extensionUrl, err := GetExtensionUrl()
	if err != nil {
		fmt.Println("can not find extion url:", err)
		return
	}
	chrome_extension_url = "chrome-extension://" + extensionUrl + "/popup.html"
	fmt.Println("chrome_extension_url:" + chrome_extension_url)

	urls := CheckSiteAvailable(ips)

	fmt.Println("found valid urls len:", len(urls))

	// api 爆破
	results := LoadExtensionWrapper(urls)

	WriteResults(results, "results.json")

}
