package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type HTMLClient struct {
	BaseURL string
	Client  *http.Client
}

func NewHTMLClient(baseURL string) *HTMLClient {
	return &HTMLClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *HTMLClient) Name() string {
	return "html"
}

func (c *HTMLClient) Fetch(urlStr string) (string, string, []string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", "", nil, err
	}
	if parsedURL.Scheme == "" {
		urlStr = c.BaseURL + urlStr
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", "", nil, err
	}

	req.Header.Set("User-Agent", "llm-data-collector/1.0")

	resp, err := c.Client.Do(req)
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
	title := extractTitle(htmlStr)
	text := extractText(htmlStr)
	links := extractLinks(htmlStr)

	return title, text, links, nil
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