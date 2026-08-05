package main

import (
  "context"
  "fmt"
  "strings"
  "sync"
  "net/url"
  
  "github.com/chromedp/chromedp"
  "github.com/chromedp/cdproto/network"
)

type Snapshot struct {
  URL   string
  HTML  string
  Title string
  Files map[string][]byte
}

type Crawler struct {
  headers map[string]string
}

func NewCrawler() *Crawler {
  return &Crawler {
    headers: make(map[string]string),
  }
}

func (c *Crawler) AddHeader(name, value string) {
	c.headers[name] = value
}

func (c *Crawler) AddHeaderString(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid header: %q", s)
	}
	c.headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	return nil
}

func (c *Crawler) Crawl(ctx context.Context, target string) (*Snapshot, error) {
  u, err := url.ParseRequestURI(target)
  if err != nil {
    return nil, err
  }
  snap := &Snapshot {
    URL: u.String(),
    Files: make(map[string][]byte),
  }
  headers := make(network.Headers)
  for k, v := range c.headers {
    headers[k] = v
  }
  if err := chromedp.Run(ctx,
    network.Enable(),
    network.SetExtraHTTPHeaders(headers),
  ); err != nil {
    return nil, err
  }
  var (
    mu sync.Mutex
    reqs = make(map[network.RequestID]string)
    wg sync.WaitGroup
  )
  chromedp.ListenTarget(ctx, func(ev interface{}) {
    switch e := ev.(type) {
    case *network.EventResponseReceived:
      if e.Type != network.ResourceTypeScript {
        return
      }
      mu.Lock()
      reqs[e.RequestID] = e.Response.URL
      mu.Unlock()
    case *network.EventLoadingFinished:
      mu.Lock()
      url, ok := reqs[e.RequestID]
      if ok {
        delete(reqs, e.RequestID)
      }
      mu.Unlock()
      if !ok {
        return
      }
      wg.Add(1)
      go func(id network.RequestID, url string) {
        defer wg.Done()
        var body []byte
        err := chromedp.Run(ctx,
          chromedp.ActionFunc(func(ctx context.Context) error {
            var err error
            body, err = network.GetResponseBody(id).Do(ctx)
            return err
          }),
        )
        if err != nil {
          return
        }
        mu.Lock()
        snap.Files[url] = body
        mu.Unlock()
      }(e.RequestID, url)
    }
  })
  var html string
  var title string
  if err := chromedp.Run(ctx,
    chromedp.Navigate(snap.URL),
    chromedp.OuterHTML("html", &html),
    chromedp.Title(&title),
  ); err != nil {
    return nil, err
  }
  wg.Wait()
  snap.HTML = html
  snap.Title = title
  return snap, nil
}

