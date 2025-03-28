package main

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

func CheckUrl(ipAddress string) bool {
	// 为每个路径发送请求

	url := ipAddress

	// 创建一个HTTP请求
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		// fmt.Printf("Error accessing %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	// 检查状态码是否为200
	if resp.StatusCode == http.StatusOK || (resp.StatusCode >= 300 && resp.StatusCode < 400) {
		return true
	}

	return false
}

func CheckUrlWrapper(ipAddress string, resCh chan string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	isValid := CheckUrl(ipAddress)
	if isValid {
		resCh <- ipAddress
	}
}

func CheckSiteAvailable(ips []string) []string {
	var urls []string
	var wg0 sync.WaitGroup
	sem0 := make(chan struct{}, 100)
	resCh := make(chan string, len(ips))

	// 1. 找到登录页面
	for _, ipAddress := range ips {
		wg0.Add(1)
		sem0 <- struct{}{}
		go CheckUrlWrapper(ipAddress, resCh, &wg0, sem0)
	}
	wg0.Wait()
	close(resCh)
	for result := range resCh {
		urls = append(urls, result)
	}
	return urls
}
