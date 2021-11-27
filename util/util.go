/*
Package util
其它可以由各模块共用的功能
*/
package util

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"net/http"
)

// FetchHTMLItemInnerText 获取HTML中指定元素中的文本，有多个时返回其中一个
func FetchHTMLItemInnerText(url, cssSelector string) (result string) {
	// 避免爬虫被屏蔽
	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatalln("new request error: ", err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36498+9")

	// Request the HTML page.
	response, err := httpClient.Do(request)
	if err != nil {
		logrus.Fatalln("http client do error: ", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		logrus.Fatalf("status code error: %d %s\n", response.StatusCode, response.Status)
	}

	// Load the HTML document
	html, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		logrus.Fatalln("goquery error: ", err)
	}

	// Find the review items
	html.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the innerText
		innerText := s.Text()
		logrus.Debugf("Review %d: %s\n", i, innerText)
		result = innerText
	})
	return result
}

// FetchDynamicHTMLItemInnerText 获取动态HTML中指定元素中的文本，有多个时返回其中一个
func FetchDynamicHTMLItemInnerText(url string, selector string) (string, error) {
	// 浏览器选项
	options := []chromedp.ExecAllocatorOption{
		// 是否不显示浏览器窗口，调试时设为false
		chromedp.Flag("headless", true),
		// 不获取图片
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 设置USer-Agent，对于mvnrepository必须
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36498+9`),
	}
	//初始化浏览器参数
	chromeCtx, _ := chromedp.NewExecAllocator(context.Background(), options...)
	// 设置日志输出方式（没看出什么情况会输出日志）
	//chromeCtx, cancel := chromedp.NewContext(chromeCtx, chromedp.WithLogf(logrus.Infof))
	// 创建浏览器，或在继承的浏览器上创建新选项卡
	chromeCtx, cancel := chromedp.NewContext(chromeCtx)
	defer cancel() // 关闭浏览器
	// 执行一个空task, 用提前创建Chrome实例（可能是为了优化？）
	//chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	//创建一个上下文，超时时间为10s（应该用于网络差但未断开的情况，因为网络断开将直接返回错误）
	//chromeCtx, cancel2 := context.WithTimeout(chromeCtx, 10*time.Second)
	//defer cancel2()
	var innerText string
	// 打开浏览器
	err := chromedp.Run(chromeCtx,
		// 跳转到URL
		chromedp.Navigate(url),
		// 选择文本（浏览器发来并可见）
		chromedp.Text(selector, &innerText, chromedp.NodeVisible),
	)
	if err != nil {
		logrus.Errorln("chrome run error: %v", err)
		return "", err
	}
	return innerText, nil
}
