package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// CarousellDataStructure = Amazon data structure
type CarousellDataStructure struct {
	Title string
	Image string
	Link  string
	Price string
	Sales string
	Shop  string
}

var (
	carousellURL      = "https://www.carousell.sg/search/"
	carousellMsg      string
	totalCarousell    int
	nodeCarousellMain []*cdp.Node
)

func GetCarousellInfo(search string) (resultListing []*CarousellDataStructure, total int, msg string) {
	var carousellListItems []*CarousellDataStructure
	searchContent := strings.ReplaceAll(search, " ", "%20")
	url := carousellURL + searchContent
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
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		carousellMsg = fmt.Sprintf("Error Occured : %v", err)
		return nil, 0, carousellMsg
	}
	for {
		totalCarousell = len(nodeCarousellMain)
		if err := chromedp.Run(ctx,
			scrollCarousellAllElements(),
			chromedp.Nodes(`div[class="TpQXuJG_eo"]`, &nodeCarousellMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		); err != nil {
			carousellMsg = fmt.Sprintf("Error Occured : %v", err)
			return nil, 0, carousellMsg
		}
		fmt.Println(len(nodeCarousellMain), totalCarousell)
		if len(nodeCarousellMain) == 0 {
			break
		} else if totalCarousell == len(nodeCarousellMain) || totalCarousell > 300 {
			for i := 0; i < len(nodeCarousellMain); i++ {
				if err := chromedp.Run(ctx,
					chromedp.ScrollIntoView(`a[class="_2ezhmqseeJ"]`, chromedp.ByQuery, chromedp.FromNode(nodeCarousellMain[i])),
					chromedp.Sleep(200*time.Millisecond),
					carousellGetallInfo(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeCarousellMain[i])),
					// chromedp.Attributes(`img`, &mapCarousellImage, chromedp.ByQuery, chromedp.FromNode(nodeCarousellImage[0])),
				); err != nil {
					carousellMsg = fmt.Sprintf("Error Occured : %v", err)
					fmt.Println(err)
					return carousellListItems, totalCarousell, carousellMsg
				}
				carousellListItems = append(carousellListItems, &CarousellDataStructure{
					Title: txtCarousellTitle,
					Image: mapCarousellImage["src"],
					Link:  mapCarousellLink["href"],
					Price: txtCarousellPrice,
				})

			}
			fmt.Println(len(carousellListItems))
			carousellMsg = fmt.Sprintf("Completed Scrape Data from Carousell")
			return carousellListItems, totalCarousell, carousellMsg
		}
		// if len(nodeCarousellMain) > 0 {
		// testnum = totalCarousell
		// totalCarousell += len(nodeCarousellMain)
		// fmt.Println(testnum, totalCarousell)

		var nodeCarousellNext []*cdp.Node
		if totalCarousell < 300 {
			if err := chromedp.Run(ctx,
				chromedp.Nodes(`button[class="_3dxOPpKVs8 _2Hl0nzGgOH _3KEDnFP0dp _3AGrhxH5DS _2UF39lBLOv yYAF4gRW1m"]`, &nodeCarousellNext, chromedp.ByQuery, chromedp.AtLeast(0)),
			); err != nil {
				carousellMsg = fmt.Sprintf("Error Occured : %v", err)
				log.Fatal(err)
				return carousellListItems, totalCarousell, carousellMsg
			}
			fmt.Println(len(nodeCarousellNext))
			if len(nodeCarousellNext) == 1 {
				if err := chromedp.Run(ctx,
					chromedp.Click(`button[class="_3dxOPpKVs8 _2Hl0nzGgOH _3KEDnFP0dp _3AGrhxH5DS _2UF39lBLOv yYAF4gRW1m"]`),
					chromedp.Sleep(2*time.Second),
				); err != nil {
					carousellMsg = fmt.Sprintf("Error Occured : %v", err)
					return carousellListItems, totalCarousell, carousellMsg
				}
			}
		}
	}
	// } else {
	// 	break
	// }
	// }
	fmt.Println(len(carousellListItems))
	fmt.Println(totalCarousell)
	carousellMsg = fmt.Sprintf("Completed Scrape Data from Carousell")
	return carousellListItems, totalCarousell, carousellMsg
}

var (
	nodeCarousellImage []*cdp.Node
	mapCarousellLink   map[string]string
	mapCarousellImage  map[string]string
	txtCarousellTitle  string
	txtCarousellPrice  string
)

func carousellGetallInfo(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Text(`p[class="_1gJzwc_bJS _2rwkILN6KA Rmplp6XJNu mT74Grr7MA nCFolhPlNA lqg5eVwdBz uxIDPd3H13 _30RANjWDIv"]`, &txtCarousellTitle, opts...),
		chromedp.Text(`p[class="_1gJzwc_bJS _2rwkILN6KA Rmplp6XJNu mT74Grr7MA nCFolhPlNA lqg5eVwdBz _19l6iUes6V _3k5LISAlf6"]`, &txtCarousellPrice, opts...),
		chromedp.Attributes(`a[class="_2ezhmqseeJ"]`, &mapCarousellLink, opts...),
		chromedp.Attributes(`a[class="_2ezhmqseeJ"] img[class="P2llUzsDMi"]`, &mapCarousellImage, opts...),
	}
}

func scrollCarousellAllElements() chromedp.Tasks {
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
		chromedp.Sleep(500 * time.Millisecond),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo({top: 0,behavior: 'smooth',});`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
		// chromedp.Sleep(500 * time.Millisecond),

		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo(0,window.innerHeight*2.5);`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
		// chromedp.Sleep(500 * time.Millisecond),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo(0,window.innerHeight*3.5);`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
		// chromedp.Sleep(500 * time.Millisecond),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo(0,window.innerHeight*4.5);`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
		// chromedp.Sleep(500 * time.Millisecond),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
	}
}
