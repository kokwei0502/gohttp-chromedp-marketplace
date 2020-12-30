package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// AmazonDataStructure = Amazon data structure
type AmazonDataStructure struct {
	Title string
	Image string
	Link  string
	Price string
	Sales string
	Shop  string
}

var (
	amazonBaseURL   = "https://www.amazon.sg"
	amazonURL       = "https://www.amazon.sg/s?k="
	amazonMsg       string
	nodeAmazonMain  []*cdp.Node
	totalAmazon     int
	amazonListitems []*AmazonDataStructure
)

// GetAmazonInfo = Retrieve Amazon basic product info
func GetAmazonInfo(search string) (resultListing []*AmazonDataStructure, total int, msg string) {
	searchContent := strings.ReplaceAll(search, " ", "+")
	amazonURL = amazonURL + searchContent
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	// Create chrome window
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.Navigate(amazonURL),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		amazonMsg = fmt.Sprintf("Error Occured : %v", err)
		return nil, 0, amazonMsg
	}
	for {
		if err := chromedp.Run(ctx,
			chromedp.Nodes(`div[data-component-type="s-search-result"]`, &nodeAmazonMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		); err != nil {
			amazonMsg = fmt.Sprintf("Error Occured : %v", err)
			return amazonListitems, totalAmazon, amazonMsg
		}
		if len(nodeAmazonMain) > 0 {
			totalAmazon += len(nodeAmazonMain)
			for i := 0; i < len(nodeAmazonMain); i++ {
				if err := chromedp.Run(ctx,
					amazonGetallInfo(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeAmazonMain[i])),
				); err != nil {
					amazonMsg = fmt.Sprintf("Error Occured : %v", err)
					return amazonListitems, totalAmazon, amazonMsg
				}
				if len(nodeAmazonPrice) == 1 {
					if err := chromedp.Run(ctx,
						chromedp.Text(`span.a-price > span`, &txtAmazonPrice, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeAmazonMain[i])),
					); err != nil {
						amazonMsg = fmt.Sprintf("Error Occured : %v", err)
						return amazonListitems, totalAmazon, amazonMsg
					}
				} else {
					txtAmazonPrice = "Not Stated"
				}
				txtAmazonLink = amazonBaseURL + mapAmazonLink["href"]
				amazonListitems = append(amazonListitems, &AmazonDataStructure{
					Title: txtAmazonTitle,
					Link:  txtAmazonLink,
					Image: mapAmazonImage["src"],
					Price: txtAmazonPrice,
				})
			}
		}
		var nodeAmazonNext []*cdp.Node
		fmt.Println(totalAmazon)
		if totalAmazon < 300 {
			if err := chromedp.Run(ctx,
				chromedp.Nodes(`li[class="a-disabled a-last"]`, &nodeAmazonNext, chromedp.ByQuery, chromedp.AtLeast(0)),
			); err != nil {
				amazonMsg = fmt.Sprintf("Error Occured : %v", err)
				return amazonListitems, totalAmazon, amazonMsg
			}
			if len(nodeAmazonNext) == 0 {
				if err := chromedp.Run(ctx,
					chromedp.Click(`li[class="a-last"]`, chromedp.ByQuery),
					chromedp.Sleep(2*time.Second),
				); err != nil {
					amazonMsg = fmt.Sprintf("Error Occured : %v", err)
					return amazonListitems, totalAmazon, amazonMsg
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	amazonMsg = fmt.Sprintf("Completed Scrape Data from Amazon")
	return amazonListitems, totalAmazon, amazonMsg
}

var (
	txtAmazonTitle  string
	txtAmazonPrice  string
	txtAmazonLink   string
	mapAmazonLink   map[string]string
	mapAmazonImage  map[string]string
	nodeAmazonPrice []*cdp.Node
)

func amazonGetallInfo(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Text(`h2 > a > span`, &txtAmazonTitle, opts...),
		chromedp.Attributes(`h2 > a`, &mapAmazonLink, opts...),
		chromedp.Attributes(`img.s-image`, &mapAmazonImage, opts...),
		chromedp.Nodes(`span.a-price > span`, &nodeAmazonPrice, opts...),
	}
}
