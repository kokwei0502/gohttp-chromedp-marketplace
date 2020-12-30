package controllers

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

type ShopeeDetailsInfo struct {
	ImageList []string
}

var (
	mapListImage     []map[string]string
	nodeShopeeDImage []*cdp.Node
)

func GetShopeeDetailsInfo(url string) {
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
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`div[class="_2JMB9h V1Fpl5"]`, chromedp.ByQuery),
		chromedp.AttributesAll(`div[class="_2Fw7Qu V1Fpl5"]`, &mapListImage, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.SendKeys(`div[class="_2Fw7Qu V1Fpl5"]`, kb.Escape, chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}
}
