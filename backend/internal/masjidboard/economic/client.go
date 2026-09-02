package economic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultAPIURL  = "https://www.jamiatsa.org/wp-json/wp/v2/posts?categories=45&per_page=1&_fields=date,link,title,content"
	defaultTimeout = 20 * time.Second
)

type Client struct {
	HTTPClient *http.Client
	APIURL     string
	Now        func() time.Time
}

type wordpressPost struct {
	Date  string `json:"date"`
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
}

func (c Client) Fetch(ctx context.Context) (Indicators, error) {
	endpoint := strings.TrimSpace(c.APIURL)
	if endpoint == "" {
		endpoint = DefaultAPIURL
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: invalid API URL: %w", err)
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MasjidPi Islamic Economic Indicators")
	response, err := client.Do(req)
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: fetch: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Indicators{}, fmt.Errorf("economic indicators: unexpected HTTP status %s", response.Status)
	}
	var posts []wordpressPost
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&posts); err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: decode response: %w", err)
	}
	if len(posts) == 0 {
		return Indicators{}, fmt.Errorf("economic indicators: source returned no posts")
	}
	now := time.Now
	if c.Now != nil {
		now = c.Now
	}
	return parsePost(posts[0], now().UTC())
}

func parsePost(post wordpressPost, fetchedAt time.Time) (Indicators, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(post.Content.Rendered))
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: parse table: %w", err)
	}
	table := document.Find("table").First()
	if table.Length() == 0 {
		return Indicators{}, fmt.Errorf("economic indicators: table not found")
	}
	headings := make(map[string]int)
	table.Find("thead th").Each(func(index int, cell *goquery.Selection) {
		headings[normalizeHeading(cell.Text())] = index
	})
	required := []string{"hijri", "date", "rand-dollar", "24 carat", "22 carat", "18 carat", "silver", "nisaab", "min mahr", "mahr faatimi", "krugerrand"}
	for _, heading := range required {
		if _, ok := headings[heading]; !ok {
			return Indicators{}, fmt.Errorf("economic indicators: required column %q not found", heading)
		}
	}
	row := table.Find("tbody tr").First()
	if row.Length() == 0 {
		return Indicators{}, fmt.Errorf("economic indicators: table has no data rows")
	}
	values := make([]string, 0, len(headings))
	row.Find("td").Each(func(_ int, cell *goquery.Selection) { values = append(values, strings.TrimSpace(cell.Text())) })
	value := func(heading string) (string, error) {
		index := headings[heading]
		if index >= len(values) {
			return "", fmt.Errorf("economic indicators: row is missing %q", heading)
		}
		return values[index], nil
	}
	amount := func(heading string) (float64, error) {
		raw, err := value(heading)
		if err != nil {
			return 0, err
		}
		raw = strings.NewReplacer("R", "", ",", "", " ", "").Replace(raw)
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil || parsed <= 0 {
			return 0, fmt.Errorf("economic indicators: invalid %q value %q", heading, raw)
		}
		return parsed, nil
	}
	postDate, err := time.Parse("2006-01-02T15:04:05", post.Date)
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: invalid post date: %w", err)
	}
	dateText, _ := value("date")
	effective, err := parseEffectiveDate(dateText, postDate.Year())
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: invalid effective date %q: %w", dateText, err)
	}
	// A January post can legitimately include late-December rows from the
	// previous Gregorian year.
	if effective.After(postDate.AddDate(0, 6, 0)) {
		effective = effective.AddDate(-1, 0, 0)
	}
	hijriDay, _ := value("hijri")
	monthTitle := strings.TrimSpace(post.Title.Rendered)
	if titleDocument, titleErr := goquery.NewDocumentFromReader(strings.NewReader(monthTitle)); titleErr == nil {
		monthTitle = strings.TrimSpace(titleDocument.Text())
	}
	result := Indicators{Source: SourceName, SourceURL: post.Link, EffectiveDate: effective.Format("2006-01-02"), HijriDate: strings.TrimSpace(hijriDay + " " + monthTitle), FetchedAt: fetchedAt}
	fields := []struct {
		heading string
		target  *float64
	}{
		{"rand-dollar", &result.RandDollar}, {"24 carat", &result.Gold24Carat}, {"22 carat", &result.Gold22Carat},
		{"18 carat", &result.Gold18Carat}, {"14 carat", &result.Gold14Carat}, {"9 carat", &result.Gold9Carat},
		{"silver", &result.Silver}, {"nisaab", &result.Nisaab},
		{"min mahr", &result.MinimumMahr}, {"mahr faatimi", &result.MahrFaatimi}, {"krugerrand", &result.Krugerrand},
	}
	for _, field := range fields {
		parsed, parseErr := amount(field.heading)
		if parseErr != nil {
			return Indicators{}, parseErr
		}
		*field.target = parsed
	}
	return result, nil
}

func parseEffectiveDate(value string, year int) (time.Time, error) {
	normalized := strings.NewReplacer(
		"-", " ",
		"‐", " ",
		"‑", " ",
		"–", " ",
		"—", " ",
	).Replace(strings.TrimSpace(value))
	parts := strings.Fields(normalized)
	if len(parts) < 2 || len(parts) > 3 {
		return time.Time{}, fmt.Errorf("expected day, month and optional year")
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day %q", parts[0])
	}
	month, ok := economicMonth(parts[1])
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month %q", parts[1])
	}
	if len(parts) == 3 {
		yearText := strings.TrimLeft(parts[2], "'’")
		parsedYear, parseErr := strconv.Atoi(yearText)
		if parseErr != nil || (len(yearText) != 2 && len(yearText) != 4) {
			return time.Time{}, fmt.Errorf("invalid year %q", parts[2])
		}
		if len(yearText) == 2 {
			century := year / 100 * 100
			parsedYear += century
			if parsedYear > year+50 {
				parsedYear -= 100
			} else if parsedYear < year-50 {
				parsedYear += 100
			}
		}
		year = parsedYear
	}

	parsed := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if parsed.Year() != year || parsed.Month() != month || parsed.Day() != day {
		return time.Time{}, fmt.Errorf("invalid calendar date %q", value)
	}
	return parsed, nil
}

func economicMonth(value string) (time.Month, bool) {
	months := map[string]time.Month{
		"jan": time.January, "january": time.January,
		"feb": time.February, "february": time.February,
		"mar": time.March, "march": time.March,
		"apr": time.April, "april": time.April,
		"may": time.May,
		"jun": time.June, "june": time.June,
		"jul": time.July, "july": time.July,
		"aug": time.August, "august": time.August,
		"sep": time.September, "sept": time.September, "september": time.September,
		"oct": time.October, "october": time.October,
		"nov": time.November, "november": time.November,
		"dec": time.December, "december": time.December,
	}
	month, ok := months[strings.ToLower(strings.TrimSpace(value))]
	return month, ok
}

func normalizeHeading(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
}
