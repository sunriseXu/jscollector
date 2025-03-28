package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func LoadExtension(url string, ctx context.Context, resCh chan SiteComponent, wg *sync.WaitGroup, sem chan struct{}) (string, context.Context, context.CancelFunc, error) {
	if wg != nil {
		defer wg.Done()
	}

	if sem != nil {
		defer func() { <-sem }()
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	var res string
	var cancel context.CancelFunc
	var err error

	//新建请求头选项
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //不加载图片，提高速度
		chromedp.Flag("load-extension", extension_path),
		chromedp.Flag("disable-extensions", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	if ctx == nil {
		options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
		ctx, _ = chromedp.NewExecAllocator(context.Background(), options...)
		ctx, _ = chromedp.NewContext(ctx)
	}

	ctx, cancel = context.WithTimeout(ctx, 50*time.Second) //超时限制
	defer cancel()

	err = chromedp.Run(ctx, // 创建另一个新的上下文代表第二个新标签页
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(chrome_extension_url),
	)
	if err != nil {
		fmt.Println("load extension failed", err)
		return res, ctx, cancel, err
	}

	newCtx, _ := chromedp.NewContext(ctx)
	// 先加载页面
	err = chromedp.Run(newCtx,
		network.SetBlockedURLs([]string{
			"*.png",
			"*.svg",
			"*.gif",
			"*.jpg",
			"*.mp4",
			"*.otf",
		}),
		chromedp.Navigate("about:blank"),
		chromedp.Sleep(4*time.Second),
		chromedp.Navigate(url),
		chromedp.Sleep(8*time.Second),
	)
	if err != nil {
		fmt.Println("open page time out, continue", err)
		return res, ctx, cancel, err
	}
	test := fmt.Sprintf("document.querySelector('#results').innerText")
	err = chromedp.Run(ctx,
		chromedp.Evaluate(test, &res),
	)
	if err != nil {
		fmt.Println("load extension failed", err)
		return res, ctx, cancel, err
	}

	if res != "" {
		coms := parseVuls(res)
		resCh <- SiteComponent{
			Url:        url,
			Components: coms,
		}
	}

	return res, ctx, cancel, err
}

func WriteResults(results []SiteComponent, dest string) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("JSON 编码失败: %s", err)
	}

	// 打印或保存到文件
	// fmt.Println(string(jsonData))

	// 如果需要将 JSON 写入文件，可以使用以下代码
	file, err := os.Create(dest)
	if err != nil {
		log.Fatalf("无法创建文件: %s", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("写入文件失败: %s", err)
	}

	// fmt.Println("JSON 数据已写入 fruits.json")
	return err
}

func LoadExtensionWrapper(urls []string) []SiteComponent {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 60)
	successCh := make(chan SiteComponent, len(urls))
	var successRes []SiteComponent

	for _, url := range urls {
		wg.Add(1)
		sem <- struct{}{}
		go LoadExtension(url, nil, successCh, &wg, sem)
	}

	wg.Wait()
	close(successCh)
	for result := range successCh {
		successRes = append(successRes, result)
		if len(result.Components) > 0 {
			fmt.Println("#################################")
			fmt.Println("url:", result.Url)
			fmt.Println("component:", len(result.Components), result.Components[0].Component, result.Components[0].Version)
			fmt.Println()
		}
	}
	return successRes
}
