package routers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kokwei0502/gohttp-chromedp-marketplace/controllers"
)

var (
	globalTemplate *template.Template
	workingDir, _  = os.Getwd()
)

type PageDataStructure struct {
	TaoBaoResults     []*controllers.TaobaoDataStructure
	ShopeeResults     []*controllers.ShopeeDataStructure
	AmazonResults     []*controllers.AmazonDataStructure
	CarousellResults  []*controllers.CarousellDataStructure
	TotalResultsFound int
	MessageRender     string
}

var (
	taobaoListing    []*controllers.TaobaoDataStructure
	shopeeListing    []*controllers.ShopeeDataStructure
	amazonListing    []*controllers.AmazonDataStructure
	carousellListing []*controllers.CarousellDataStructure
	msgData          string
	totalFound       int
)

func MarketplaceSearchIndexPage(w http.ResponseWriter, r *http.Request) {
	// a, b, c := controllers.GetTaoBaoInfo("a")
	// x := PageDataStructure{
	// 	TaoBaoResults:     a,
	// 	TotalResultsFound: b,
	// 	MessageRender:     c,
	// }
	if r.Method == "POST" {
		submitVal := r.FormValue("submit")
		if submitVal == "submit-search" {
			searchContent := r.FormValue("search-content")
			selectMrtPlace := r.FormValue("search-platform")
			if searchContent != "" {
				switch selectMrtPlace {
				case "taobao-mrkplc":
					taobaoListing, totalFound, msgData = controllers.GetTaoBaoInfo(searchContent)
					shopeeListing = nil
					amazonListing = nil
					carousellListing = nil
				case "shopee-mrkplc":
					shopeeListing, totalFound, msgData = controllers.GetShopeeInfo(searchContent)
					taobaoListing = nil
					amazonListing = nil
					carousellListing = nil
				case "amazon-mrkplc":
					amazonListing, totalFound, msgData = controllers.GetAmazonInfo(searchContent)
					taobaoListing = nil
					shopeeListing = nil
					carousellListing = nil
				case "carousell-mrkplc":
					carousellListing, totalFound, msgData = controllers.GetCarousellInfo(searchContent)
					taobaoListing = nil
					shopeeListing = nil
					amazonListing = nil
				}
			} else {
				msgData = "Please Key Some Keywords instead of Blank..."
			}
		}
	}
	// for _, x := range shopeeListing {
	// 	fmt.Println(x.Title)
	// }
	pageData := &PageDataStructure{
		TaoBaoResults:     taobaoListing,
		ShopeeResults:     shopeeListing,
		AmazonResults:     amazonListing,
		CarousellResults:  carousellListing,
		TotalResultsFound: totalFound,
		MessageRender:     msgData,
	}
	globalTemplate.ExecuteTemplate(w, "search-index.html", pageData)
}

// RetrieveAllTemplate = Get all .html templates
func RetrieveAllTemplate() *template.Template {
	var htmlListing []string
	err := filepath.Walk(workingDir+"/templates/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch checkDir := info.Mode(); {
		case checkDir.IsRegular():
			htmlListing = append(htmlListing, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	globalTemplate = template.Must(template.ParseFiles(htmlListing...))
	return globalTemplate
}
