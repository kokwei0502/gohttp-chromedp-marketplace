package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// ShopeeDataStructure = Shopee data structure
type ShopeeDataStructure struct {
	Title string
	Image string
	Link  string
	Price string
	Sales string
	Shop  string
}

var (
	shopeeBaseURL  = "https://shopee.sg"
	shopeeURL      = "https://shopee.sg/search?keyword="
	shopeePageURL  = "&page="
	pageNumber     = 2
	shopeeMsg      string
	nodeShopeeMain []*cdp.Node
	totalShopee    int
)

// GetShopeeInfo = Get shopee data
func GetShopeeInfo(search string) (resultListing []*ShopeeDataStructure, total int, msg string) {
	var shopeeListItems []*ShopeeDataStructure
	searchContent := strings.ReplaceAll(search, " ", "%20")
	shopeeURL = shopeeURL + searchContent
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("start-maximized", true),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	// Create chrome window
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	var langugeSelection []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate(shopeeURL),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`div.language-selection__list`, &langugeSelection, chromedp.ByQuery, chromedp.AtLeast(0)),
	); err != nil {
		shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
		return nil, 0, shopeeMsg
	}
	if len(langugeSelection) > 0 {
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(`div.language-selection__list > div.language-selection__list-item:nth-child(1)`),
			chromedp.Click(`div.language-selection__list > div.language-selection__list-item:nth-child(1)`),
			chromedp.Sleep(500*time.Millisecond),
		); err != nil {
			shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
			return nil, 0, shopeeMsg
		}
	}
	for {
		if err := chromedp.Run(ctx,
			chromedp.Nodes(`div[data-sqe="item"]`, &nodeShopeeMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		); err != nil {
			shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
			return nil, 0, shopeeMsg
		}
		if len(nodeShopeeMain) > 0 {
			totalShopee += len(nodeShopeeMain)
			// if err := chromedp.Run(ctx,
			// 	scrollAllElements(),
			// ); err != nil {
			// 	shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
			// 	return shopeeListItems, totalShopee, shopeeMsg
			// }
			for i := 0; i < len(nodeShopeeMain); i++ {
				if err := chromedp.Run(ctx,
					chromedp.ScrollIntoView(`div[data-sqe="name"]`, chromedp.ByQuery, chromedp.FromNode(nodeShopeeMain[i])),
					shopeeCheckElements(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeShopeeMain[i])),
				); err != nil {
					shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
					return shopeeListItems, totalShopee, shopeeMsg
				}
				if len(nodeShopeeLink) == 1 {
					if err := chromedp.Run(ctx,
						shopeeGetallInfo(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeShopeeMain[i])),
					); err != nil {
						shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
						return shopeeListItems, totalShopee, shopeeMsg
					}
					if txtShopeeSales == "" {
						txtShopeeSales = "No Records"
					}
					txtShopeeLink = shopeeBaseURL + mapShopeeLink["href"]

					shopeeListItems = append(shopeeListItems, &ShopeeDataStructure{
						Title: txtShopeeTitle,
						Image: mapShopeeImg["src"],
						Link:  txtShopeeLink,
						Price: txtShopeePrice,
						Sales: txtShopeeSales,
						// Shop  string,
					})
				} else {
					continue
				}
			}
		} else {
			shopeeMsg = fmt.Sprintf("No More Results Found...")
			return shopeeListItems, totalShopee, shopeeMsg
		}
		fmt.Println(totalShopee)
		if len(nodeShopeeMain) == 50 && totalShopee < 300 {
			if err := chromedp.Run(ctx,
				shopeeNextPage(),
			); err != nil {
				shopeeMsg = fmt.Sprintf("Error Occured : %v", err)
				return shopeeListItems, totalShopee, shopeeMsg
			}
			pageNumber++
		} else {
			break
		}
	}
	fmt.Println(len(shopeeListItems))
	shopeeMsg = fmt.Sprintf("Completed Scrape Data from Shopee")
	return shopeeListItems, totalShopee, shopeeMsg
}

var (
	mapShopeeLink  map[string]string
	mapShopeeImg   map[string]string
	txtShopeeTitle string
	txtShopeePrice string
	txtShopeeSales string
	txtShopeeLink  string
)

// data-sqe="link"
func shopeeGetallInfo(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Attributes(`a[data-sqe="link"]`, &mapShopeeLink, opts...),
		chromedp.Attributes(`a[data-sqe="link"] img`, &mapShopeeImg, opts...),
		chromedp.Text(`div[data-sqe="name"]>div`, &txtShopeeTitle, opts...),
		chromedp.Text(`div[class="_1w9jLI _37ge-4 _2ZYSiu"`, &txtShopeePrice, opts...),
		chromedp.Text(`div[class="_2-i6yP"]>div:nth-child(3)`, &txtShopeeSales, opts...),
	}
}

var (
	nodeShopeeLink []*cdp.Node
)

func shopeeCheckElements(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Nodes(`a[data-sqe="link"]`, &nodeShopeeLink, opts...),
		// chromedp.Sleep(500 * time.Millisecond),
	}
}

func scrollAllElements() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight/3);`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight/2);`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight+600);`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
	}
}

func shopeeNextPage() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(shopeeURL + shopeePageURL + strconv.Itoa(pageNumber)),
		chromedp.Sleep(2 * time.Second),
	}
}
