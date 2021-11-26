/*
Package util
其它可以由各模块共用的功能
*/
package util

import (
	"github.com/PuerkitoBio/goquery"
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
