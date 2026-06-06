package htmlcrawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"DataCollector/internal/crawler"
	"DataCollector/internal/models"

	"golang.org/x/net/html"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Name() string {
	return "html"
}

func (c *Client) Fetch(urlStr string) (title string, content string, links []string, err error) {
	parsedURL, _ := url.Parse(urlStr)
	if parsedURL.Scheme == "" {
		urlStr = c.BaseURL + urlStr
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", "", nil, err
	}

	req.Header.Set("User-Agent", "llm-data-collector/1.0")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", nil, err
	}

	htmlStr := string(body)
	title = extractTitle(htmlStr)
	content = extractText(htmlStr)
	links = extractLinks(htmlStr)

	return title, content, links, nil
}

func extractTitle(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

func extractText(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	var text strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
			text.WriteString(" ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return text.String()
}

func extractLinks(htmlStr string) []string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

func New(baseURL string, seeds []string) *crawler.Collector {
	return crawler.NewCollector(
		NewClient(baseURL),
		crawler.WithSeeds(seeds),
		crawler.WithLanguage("en"),
		crawler.WithDocType(models.ArticleDocument),
	)
}