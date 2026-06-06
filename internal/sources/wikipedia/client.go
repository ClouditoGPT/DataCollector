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
