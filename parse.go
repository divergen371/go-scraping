package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func parseList(resp *http.Response) ([]Item, error) {
	//レスポンスボディを取得
	body := resp.Body

	//レスポンスに含まれているリクエスト情報からリクエスト先のURLを取得
	requestURL := *resp.Request.URL

	var items []Item
	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, fmt.Errorf("get document error: %w", err)
	}

	// Find関数でテーブル行の行要素を取得
	tr := doc.Find("table tr")
	notFoundMessage := "Page not found"
	if strings.Contains(doc.Text(), notFoundMessage) || tr.Size() == 0 {
		return nil, nil
	}

	tr.Each(func(_ int, s *goquery.Selection) {
		item := Item{}
		// Find関数で商品の各要素を取得
		item.Name = s.Find("td:nth-of-type(2)").Text()
		item.Price, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s.
			Find("td:nth-of-type(3)").Text(), ",", ""), "円", ""))
		itemURL, exists := s.Find("td:nth-of-type(2) a").Attr("href")
		refURL, parseErr := url.Parse(itemURL)
		if exists && parseErr == nil {
			//requestURLとrefURLを結合して絶対URLを取得
			item.URL = requestURL.ResolveReference(refURL).String()
		}
		if item.Name != "" {
			items = append(items, item)
		}
	})
	return items, nil
}
