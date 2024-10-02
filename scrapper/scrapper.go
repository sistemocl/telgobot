package scrapper

import (
	"context"
	"strings"
	"tel_gobot/models"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

func FetchDataFromTable(url string) ([]models.PPSQueueSize, error) {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	tasks, err := parseHTML(htmlContent)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func parseHTML(htmlContent string) ([]models.PPSQueueSize, error) {
	var tasks []models.PPSQueueSize

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var parseTable func(*html.Node)
	parseTable = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "tableId" {
					// Found the table with the specified id
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "tbody" {
							for tr := c.FirstChild; tr != nil; tr = tr.NextSibling {
								if tr.Type == html.ElementNode && tr.Data == "tr" {
									var task models.PPSQueueSize
									tdIndex := 0
									for td := tr.FirstChild; td != nil; td = td.NextSibling {
										if td.Type == html.ElementNode && td.Data == "td" {
											switch tdIndex {
											case 0:
												if td.FirstChild != nil && td.FirstChild.Type == html.TextNode {
													task.PPSID = td.FirstChild.Data
												}
											case 16:
												if td.FirstChild != nil && td.FirstChild.Type == html.TextNode {
													task.QueueSize = td.FirstChild.Data
												}
											}
											tdIndex++
										}
									}
									tasks = append(tasks, task)
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseTable(c)
		}
	}

	parseTable(doc)

	return tasks, nil
}
