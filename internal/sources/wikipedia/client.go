package wikipedia

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type QueryResponse struct {
	Query struct {
		Pages map[string]struct {
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

type LinkResponse struct {
	Query struct {
		Pages map[string]struct {
			Title string `json:"title"`
			Links []struct {
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Name() string {
	return "wikipedia"
}

func (c *Client) Fetch(title string) (string, string, []string, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("prop", "extracts|links")
	params.Set("explaintext", "1")
	params.Set("pllimit", "max")
	params.Set("titles", title)

	req, _ := http.NewRequest("GET", c.BaseURL+"?"+params.Encode(), nil)
	req.Header.Set("User-Agent", "llm-data-collector/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", nil, nil
	}

	var data struct {
		Query struct {
			Pages map[string]struct {
				Title   string `json:"title"`
				Extract string `json:"extract"`
				Links   []struct {
					Title string `json:"title"`
				} `json:"links"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", nil, err
	}

	for _, page := range data.Query.Pages {
		var links []string
		for _, l := range page.Links {
			if l.Title != "" {
				links = append(links, l.Title)
			}
		}
		return page.Title, page.Extract, links, nil
	}
	return "", "", nil, nil
}