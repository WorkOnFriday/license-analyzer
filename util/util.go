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
	"log"
	"net/http"
	"time"
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
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36498+9`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	//创建一个上下文，超时时间为10s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 10*time.Second)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.Text(selector, &htmlContent, chromedp.NodeVisible),
	)
	if err != nil {
		logrus.Errorln("Run err : %v", err)
		return "", err
	}
	return htmlContent, nil
}
