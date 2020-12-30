package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// AlibabaCNDataStructure = AlibabaCN data structure
type AlibabaCNDataStructure struct {
	Title string
	Image string
	Link  string
	Price string
	Sales string
	Shop  string
}

var (
	alibabaCNURL                          = "https://www.1688.com/"
	alibabaCNMsg                          string
	nodeAlibabaCNMain, nodeAlibabaCNLogin []*cdp.Node
	totalAlibabaCN                        int
	alibabaListItems                      []*AlibabaCNDataStructure
	sliceAlibabaCNImage                   []string
	txtAlibabaCNImage                     string
)

// GetAlibabaCNInfo = Retrieve Alibaba China product basic data
func GetAlibabaCNInfo(search string) (resultListing []*AlibabaCNDataStructure, total int, msg string) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
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
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	// login-box loading
	if err := chromedp.Run(ctx,
		chromedp.Navigate(alibabaCNURL),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`#home-header-searchbox`, search, chromedp.ByID),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Click(`button[class="single"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`div[class="login-box loading"]`, &nodeAlibabaCNLogin, chromedp.AtLeast(0)),
	); err != nil {
		alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
		return nil, 0, alibabaCNMsg
	}
	fmt.Println(len(nodeAlibabaCNLogin))
	if len(nodeAlibabaCNLogin) == 1 {
		if err := chromedp.Run(ctx,
			chromedp.SendKeys(`#fm-login-id`, "kokwei0502", chromedp.ByID),
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`#fm-login-password`, "kwei4188", chromedp.ByID),
			chromedp.Sleep(500*time.Millisecond),
			chromedp.Click(`button[class="fm-button fm-submit password-login"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`#alisearch-input`, chromedp.ByID),
			chromedp.Sleep(2*time.Second),
		); err != nil {
			alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
			return nil, 0, alibabaCNMsg
		}
	} else {
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(`#alisearch-input`, chromedp.ByID),
			chromedp.Sleep(2*time.Second),
			chromedp.Nodes(`div[class="login-box loading"]`, &nodeAlibabaCNLogin, chromedp.AtLeast(0)),
		); err != nil {
			alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
			return nil, 0, alibabaCNMsg
		}
	}
	for {
		if err := chromedp.Run(ctx,
			scrollAlibabaCNAllElements(),
			chromedp.Nodes(`div[class="card-container"]`, &nodeAlibabaCNMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		); err != nil {
			alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
			return alibabaListItems, totalAlibabaCN, alibabaCNMsg
		}
		fmt.Println(len(nodeAlibabaCNMain))
		if len(nodeAlibabaCNMain) > 0 {
			totalAlibabaCN += len(nodeAlibabaCNMain)
			for i := 0; i < len(nodeAlibabaCNMain); i++ {
				chromedp.Run(ctx,
					alibabaCNGetAllInfo(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeAlibabaCNMain[i])),
				)
				txtAlibabaCNImage = mapAlibabaCNImage["style"]
				sliceAlibabaCNImage = strings.Split(txtAlibabaCNImage, `url("`)
				txtAlibabaCNImage = sliceAlibabaCNImage[len(sliceAlibabaCNImage)-1]
				sliceAlibabaCNImage = strings.Split(txtAlibabaCNImage, `?_`)
				txtAlibabaCNImage = sliceAlibabaCNImage[0]
				txtAlibabaCNPrice = txtAlibabaCNCurrency + txtAlibabaCNPrice
				alibabaListItems = append(alibabaListItems, &AlibabaCNDataStructure{
					Title: txtAlibabaCNTitle,
					Image: txtAlibabaCNImage,
					Link:  mapAlibabaCNLink["href"],
					Price: txtAlibabaCNPrice,
					Sales: txtAlibabaCNSales,
				})
			}
		}
		var nodeAlibabaCNNext []*cdp.Node
		if totalAlibabaCN < 300 {
			if err := chromedp.Run(ctx,
				chromedp.Nodes(`div[class="sm-pagination"] a[class="fui-next fui-next-disabled"]`, &nodeAlibabaCNNext, chromedp.ByQuery, chromedp.AtLeast(0)),
			); err != nil {
				alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
				return alibabaListItems, totalAlibabaCN, alibabaCNMsg
			}
			if len(nodeAlibabaCNNext) == 0 {
				if err := chromedp.Run(ctx,
					chromedp.WaitVisible(`div[class="sm-pagination"] a[class="fui-next"]`, chromedp.ByQuery),
					chromedp.Click(`div[class="sm-pagination"] a[class="fui-next"]`, chromedp.ByQuery),
					chromedp.Sleep(2*time.Second),
				); err != nil {
					alibabaCNMsg = fmt.Sprintf("Error Occured : %v", err)
					return alibabaListItems, totalAlibabaCN, alibabaCNMsg
				}
			}

		} else {
			break
		}
	}
	alibabaCNMsg = fmt.Sprintf("Completed Scrape Data from Alibaba China")
	return alibabaListItems, totalAlibabaCN, alibabaCNMsg
}

var (
	mapAlibabaCNLink     map[string]string
	mapAlibabaCNImage    map[string]string
	txtAlibabaCNTitle    string
	txtAlibabaCNPrice    string
	txtAlibabaCNCurrency string
	txtAlibabaCNSales    string
)

func alibabaCNGetAllInfo(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ScrollIntoView(`div[class="img-container"]`, opts...),
		chromedp.Attributes(`div[class="img-container"] a`, &mapAlibabaCNLink, opts...),
		chromedp.Attributes(`div[class="img-container"] div[class="img"]`, &mapAlibabaCNImage, opts...),
		chromedp.Text(`div[class="desc-container"] div[class="title"]`, &txtAlibabaCNTitle, opts...),
		chromedp.Text(`div[class="price-container"] div[class="price"]`, &txtAlibabaCNPrice, opts...),
		chromedp.Text(`div[class="price-container"] div[class="rmb"]`, &txtAlibabaCNCurrency, opts...),
		chromedp.Text(`div[class="sale"]`, &txtAlibabaCNSales, opts...),
	}
}

func scrollAlibabaCNAllElements() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo({top: 5000,behavior: 'smooth',});`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.Sleep(1500 * time.Millisecond),
	}
}
