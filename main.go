package main

import (
  "context"
  "fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func run() error {
	// creating browser window
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", cfg.Headless),
		chromedp.Flag("disable-gpu", true),
	)
	if cfg.Profile != "" {
	  opts = append(opts,
	    chromedp.Flag("user-data-dir", cfg.Profile),
	   )
	}
	// setting up the browser window
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	// crawler part
	crawler := NewCrawler()
	for _, h := range cfg.Headers {
	  if err := crawler.AddHeaderString(h); err != nil {
	    return err
	  }
	}
	snapshot, err := crawler.Crawl(ctx, cfg.TargetURL)
	if err != nil {
	  return err
	}
	// scanner part
  scanner := NewScanner()
  if !cfg.SkipDefaultRegexp {
    if err := scanner.LoadDefaultRegexes(); err != nil {
      return err
    }
  }
  for _, val := range cfg.CustomRegexps {
	  parts := strings.SplitN(val, ":", 2)
	  if len(parts) != 2 {
		  return fmt.Errorf("invalid regexp: %q", val)
	  }
	  name := strings.TrimSpace(parts[0])
	  value := strings.TrimSpace(parts[1])
	  if err := scanner.AddRegex(name, value); err != nil {
	    return err
	  }
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
