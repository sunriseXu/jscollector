package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func GetExtensionUrl() (string, error) {
	var ctx context.Context
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

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, _ = chromedp.NewExecAllocator(context.Background(), options...)
	ctx, _ = chromedp.NewContext(ctx)

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second) //超时限制
	defer cancel()
	//#extension-id
	test := fmt.Sprintf("document.querySelector('extensions-manager').shadowRoot.querySelector('#items-list').shadowRoot.querySelector('extensions-item').id")

	err = chromedp.Run(ctx, // 创建另一个新的上下文代表第二个新标签页
		chromedp.Navigate("chrome://extensions"),
		chromedp.Evaluate(test, &res),
	)
	fmt.Println("res:", res)
	if err != nil {
		fmt.Println("load extension failed", err)
		return res, err
	}
	return res, err
}
