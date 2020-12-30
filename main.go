package main

import (
	"log"
	"net/http"

	"github.com/kokwei0502/gohttp-chromedp-marketplace/routers"
)

func init() {
	routers.RetrieveAllTemplate()
}

func main() {
	// controllers.GetTaoBaoInfo("被单")
	// controllers.GetShopeeInfo("被单")
	// controllers.GetAmazonInfo("dell laptop")
	// controllers.GetCarousellInfo("asus laptop")
	// controllers.GetAlibabaCNInfo("被单")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", routers.MarketplaceSearchIndexPage)
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// ctx, cancel := chromedp.NewContext(
// 	context.Background(),
// 	chromedp.WithLogf(log.Printf),
// )
// defer cancel()

// // Set the timeout limit
// ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
// defer cancel()

// // Start to navigate to the jobsDB URL
// if err := chromedp.Run(ctx,
// 	chromedp.Navigate(jobsDBURL),
// 	chromedp.Sleep(1*time.Second),
// ); err != nil {
// 	jobsDBErrMsg = fmt.Sprintf("Navigation Error, Please Try Again Later!\n%v", err)
// 	fmt.Println(jobsDBErrMsg)
// 	return jobsdblist, jobsDBErrMsg, 0
// }

// // Looping to get data until timeout limit reached / total data limit reached / pages limit reached
// for {
// 	// Start to get the main section
// 	if err := chromedp.Run(ctx,
// 		chromedp.Nodes(`div[class="job-container result organic-job"]`, &nodeJobsDBMain, chromedp.ByQueryAll, chromedp.AtLeast(0)),
// 	); err != nil {
// 		jobsDBErrMsg = fmt.Sprintf("Error to Get Main Division, Please Try Again Later\n%v", err)
// 		return jobsdblist, jobsDBErrMsg, 0
// 	}
