package wikipedia

import (
	"encoding/json"
	"fmt"
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

func FetchArticle(
	title string,
) (string, string, error) {

	base := "https://fa.wikipedia.org/w/api.php"

	params := url.Values{}

	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("prop", "extracts")
	params.Set("explaintext", "1")
	params.Set("titles", title)

	url := fmt.Sprintf(
		"%s?%s",
		base,
		params.Encode(),
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set(
		"User-Agent",
		"llm-data-collector/1.0 (contact: your-email@example.com)",
	)

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf(
			"wiki error: %d",
			resp.StatusCode,
		)
	}

	defer resp.Body.Close()

	var data QueryResponse

	err = json.NewDecoder(resp.Body).
		Decode(&data)

	if err != nil {
		return "", "", err
	}

	for _, page := range data.Query.Pages {
		return page.Title, page.Extract, nil
	}

	return "", "", nil
}

func FetchLinks(title string) ([]string, error) {

	base := "https://fa.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("prop", "links")
	params.Set("pllimit", "max")
	params.Set("titles", title)

	req, err := http.NewRequest("GET", base+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"llm-data-collector/1.0",
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("wiki error: %d", resp.StatusCode)
	}

	var data LinkResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	links := []string{}

	for _, page := range data.Query.Pages {
		for _, l := range page.Links {
			if l.Title != "" {
				links = append(links, l.Title)
			}
		}
	}

	return links, nil
}
