package dailycontent

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	defaultTimeout   = 20 * time.Second
	maxResponseSize  = 1 << 20
	javascriptPrefix = "let translations ="
)

var (
	breakPattern   = regexp.MustCompile(`(?i)<br\s*/?>`)
	tagPattern     = regexp.MustCompile(`<[^>]*>`)
	dateKeyPattern = regexp.MustCompile(`^(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun) ([A-Z][a-z]{2}) ([0-9]{2}) ([0-9]{4}) `)
)

type Client struct {
	HTTPClient *http.Client
	APIURL     string
	Language   string
	Now        func() time.Time
}

func (c Client) Fetch(ctx context.Context) (Content, error) {
	endpoint := strings.TrimSpace(c.APIURL)
	if endpoint == "" {
		endpoint = SourceURL
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return Content{}, fmt.Errorf("daily Islamic content: invalid API URL: %w", err)
	}
	language := strings.TrimSpace(c.Language)
	if language == "" {
		language = "en"
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Content{}, fmt.Errorf("daily Islamic content: create request: %w", err)
	}
	req.Header.Set("Accept", "text/javascript, application/javascript, text/plain;q=0.9, */*;q=0.1")
	req.Header.Set("User-Agent", "MasjidPi Daily Islamic Content")
	response, err := client.Do(req)
	if err != nil {
		return Content{}, fmt.Errorf("daily Islamic content: fetch: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Content{}, fmt.Errorf("daily Islamic content: unexpected HTTP status %s", response.Status)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return Content{}, fmt.Errorf("daily Islamic content: read response: %w", err)
	}
	if len(body) > maxResponseSize {
		return Content{}, fmt.Errorf("daily Islamic content: response exceeds %d bytes", maxResponseSize)
	}
	now := time.Now
	if c.Now != nil {
		now = c.Now
	}
	return parse(body, language, endpoint, now().UTC())
}

func parse(body []byte, language, endpoint string, fetchedAt time.Time) (Content, error) {
	source := strings.TrimSpace(strings.TrimPrefix(string(body), "\ufeff"))
	if !strings.HasPrefix(source, javascriptPrefix) {
		return Content{}, fmt.Errorf("daily Islamic content: expected JavaScript translations assignment")
	}
	source = strings.TrimSpace(strings.TrimPrefix(source, javascriptPrefix))
	source = strings.TrimSuffix(source, ";")
	var values map[string]json.RawMessage
	decoder := json.NewDecoder(strings.NewReader(source))
	if err := decoder.Decode(&values); err != nil {
		return Content{}, fmt.Errorf("daily Islamic content: decode translations: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return Content{}, fmt.Errorf("daily Islamic content: unexpected data after translations object")
	}
	value := func(field string) (string, error) {
		raw, ok := values[field]
		if !ok {
			return "", fmt.Errorf("daily Islamic content: required field %q not found", field)
		}
		var translations map[string]string
		if err := json.Unmarshal(raw, &translations); err != nil {
			return "", fmt.Errorf("daily Islamic content: decode field %q: %w", field, err)
		}
		return normalizeText(translations[language]), nil
	}
	fields := make(map[string]string, 9)
	for _, field := range []string{"ayahSurah", "AyahNo", "ayah", "hadithHeading", "hadith", "hadithRef", "sunnahHeading", "sunnah", "sunnahRef"} {
		parsed, err := value(field)
		if err != nil {
			return Content{}, err
		}
		fields[field] = parsed
	}
	content := Content{
		Ayah:     Ayah{Surah: fields["ayahSurah"], AyahNumber: fields["AyahNo"], Text: fields["ayah"]},
		Hadith:   Hadith{Heading: fields["hadithHeading"], Text: fields["hadith"], Reference: fields["hadithRef"]},
		Sunnah:   Sunnah{Heading: fields["sunnahHeading"], Text: fields["sunnah"], Reference: fields["sunnahRef"]},
		Language: language, Source: SourceName, SourceURL: endpoint, ContentDate: contentDate(values), FetchedAt: fetchedAt,
	}
	if !content.Valid() {
		return Content{}, fmt.Errorf("daily Islamic content: source returned incomplete %q content", language)
	}
	return content, nil
}

func normalizeText(value string) string {
	value = breakPattern.ReplaceAllString(value, "\n")
	value = tagPattern.ReplaceAllString(value, "")
	value = html.UnescapeString(value)
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	lines := strings.Split(value, "\n")
	normalized := make([]string, 0, len(lines))
	blank := false
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			if len(normalized) > 0 && !blank {
				normalized = append(normalized, "")
			}
			blank = true
			continue
		}
		normalized = append(normalized, line)
		blank = false
	}
	return strings.TrimSpace(strings.Join(normalized, "\n"))
}

func contentDate(values map[string]json.RawMessage) string {
	for key := range values {
		match := dateKeyPattern.FindStringSubmatch(key)
		if len(match) != 4 {
			continue
		}
		parsed, err := time.Parse("Jan 02 2006", strings.Join(match[1:], " "))
		if err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return ""
}
