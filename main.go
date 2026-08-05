package main

import (
  "context"
  "fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func run() error {
	// creating browser window
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", cbrHeadless),
		chromedp.Flag("disable-gpu", true),
	)
	if cbrProfile != "" {
	  opts = append(opts,
	    chromedp.Flag("user-data-dir", cbrProfile),
	   )
	}
	// setting up the browser window
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	crawler := NewCrawler()
	for _, h := range cbrHeaders {
	  if err := crawler.AddHeaderString(h); err != nil {
	    return err
	  }
	}
	snapshot, err := crawler.Crawl(ctx, cbrTargetURL)
	if err != nil {
	  return err
	}
	// creating scanner
  scanner := NewScanner()
  if err := scanner.LoadDefaultRegexes(); err != nil {
    return err
  }
  if err := scanner.LoadCustomRegexes(); err != nil {
    return err
  }
  for url, contents := range snapshot.Files {
    matches := scanner.Scan(url, contents)
    for _, match := range matches {
      fmt.Printf("[%s]: %s\n", match.Rule, match.Value)
    }
  }
	return nil
}

func main() {
  if err := rootCmd.Execute(); err != nil {
    log.Fatal(err)
  }
}
