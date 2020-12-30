package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type taobaoConfig struct {
	TaobaoURL string `json:"Taobao URL"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
}

// TaobaoDataStructure = Taobao data structure
type TaobaoDataStructure struct {
	Title string
	Image string
	Link  string
	Price string
	Sales string
	Shop  string
}

var (
	taobaoMsg                       string
	nodeTaobaoLogin, nodeTaobaoMain []*cdp.Node
	taobaoLoginTooltip              string
	taobaoconfig                    *taobaoConfig
	totalTaoBao                     int
	taobaoResultListing             []*TaobaoDataStructure
)

// GetTaoBaoInfo = Get TaoBao items info
func GetTaoBaoInfo(search string) (resultListing []*TaobaoDataStructure, total int, msg string) {
	config := getTaobaoConfig()
	searchContent := strings.ReplaceAll(search, " ", "+")
	taobaoURL := config.TaobaoURL + searchContent
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("ignore-certificate-errors", true),
		// chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
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

	if err := chromedp.Run(ctx,
		chromedp.Navigate(taobaoURL),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`.login-form`, &nodeTaobaoLogin, chromedp.ByQuery, chromedp.AtLeast(0)),
	); err != nil {
		taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
		return taobaoResultListing, totalTaoBao, taobaoMsg
	}
	if len(nodeTaobaoLogin) == 1 {
		if err := chromedp.Run(ctx,
			taobaoLogin(),
		); err != nil {
			taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
			return taobaoResultListing, totalTaoBao, taobaoMsg
		}
	} else {
		return taobaoResultListing, totalTaoBao, taobaoMsg
	}
	for {
		if err := chromedp.Run(ctx,
			taobaoGetMainSection(),
		); err != nil {
			taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
			return taobaoResultListing, totalTaoBao, taobaoMsg
		}
		totalTaoBao += len(nodeTaobaoMain)
		if len(nodeTaobaoMain) > 0 {
			for i := 0; i < len(nodeTaobaoMain); i++ {
				if err := chromedp.Run(ctx,
					taobaoGetallInfo(chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(nodeTaobaoMain[i])),
				); err != nil {
					taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
					return taobaoResultListing, totalTaoBao, taobaoMsg
				}
				taobaoResultListing = append(taobaoResultListing, &TaobaoDataStructure{
					Title: mapTaobaoImg["alt"],
					Image: mapTaobaoImg["data-src"],
					Link:  mapTaobaoLink["href"],
					Price: txtTaobaoPrice,
					Sales: txtTaobaoSales,
					Shop:  mapTaobaoShop["href"],
				})
			}
		} else {
			taobaoMsg = fmt.Sprintf("Done, Check Results Below...")
			return taobaoResultListing, totalTaoBao, taobaoMsg
		}
		if err := chromedp.Run(ctx,
			taobaoCheckNextPage(),
		); err != nil {
			taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
			return taobaoResultListing, totalTaoBao, taobaoMsg
		}
		if totalTaoBao < 200 {
			if len(nodeTaobaoNextPage) == 0 {
				if err := chromedp.Run(ctx,
					chromedp.WaitVisible(`ul[class="items"] > li[class="item next"]`),
					chromedp.Click(`ul[class="items"] > li[class="item next"]`, chromedp.ByQuery),
					chromedp.Sleep(2*time.Second),
				); err != nil {
					taobaoMsg = fmt.Sprintf("Error Occured : %v", err)
					return taobaoResultListing, totalTaoBao, taobaoMsg
				}
			} else {
				taobaoMsg = fmt.Sprintf("Done, Check Results Below...")
				return taobaoResultListing, totalTaoBao, taobaoMsg
			}
		}
	}
}

var (
	mapTaobaoLink     map[string]string
	mapTaobaoImg      map[string]string
	mapTaobaoShop     map[string]string
	txtTaobaoShopName string
	txtTaobaoPrice    string
	txtTaobaoSales    string
)

// price g_price g_price-highlight
// deal-cnt
var nodeTaobaoNextPage []*cdp.Node

func taobaoCheckNextPage() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Nodes(`ul[class="items"] > li[class="item next next-disabled"]`, &nodeTaobaoNextPage, chromedp.ByQuery, chromedp.AtLeast(0)),
	}
}

// shopname
func taobaoGetallInfo(opts ...func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Attributes(`a.pic-link`, &mapTaobaoLink, opts...),
		chromedp.Attributes(`img[class="J_ItemPic img"]`, &mapTaobaoImg, opts...),
		chromedp.Text(`div[class="price g_price g_price-highlight"]`, &txtTaobaoPrice, opts...),
		chromedp.Text(`div[class="deal-cnt"]`, &txtTaobaoSales, opts...),
		chromedp.Attributes(`a.shopname`, &mapTaobaoShop, opts...),
		chromedp.Text(`a.shopname`, &txtTaobaoShopName, opts...),
	}
}

func taobaoGetMainSection() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`div[data-category="auctions"]`),
		chromedp.Nodes(`div[data-category="auctions"]`, &nodeTaobaoMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
	}
}

func taobaoLogin() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#fm-login-id`, chromedp.ByID),
		chromedp.SendKeys(`#fm-login-id`, "kokwei0502", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#fm-login-password`, chromedp.ByID),
		chromedp.SendKeys(`#fm-login-password`, "kwei4188", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`button[class="fm-button fm-submit password-login"]`, chromedp.ByQuery),
		chromedp.Click(`button[class="fm-button fm-submit password-login"]`, chromedp.ByQuery),
	}
}

func getTaobaoConfig() *taobaoConfig {
	config, err := ioutil.ReadFile("./static/config/taobao.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(config, &taobaoconfig)
	return taobaoconfig
}
