package masjidboardlive

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var (
	premiumBoardIDMarker = []byte("let boardId")
	premiumInfoMarker    = []byte("let theInfo")
)

// ExtractPremiumPageData resolves the opaque API board ID and the initial
// 29-row payload embedded in a generated MasjidBoard Live Premium page.
// The public mid remains the stable external identifier; boardID is an
// upstream implementation detail that may change when the page is rebuilt.
func ExtractPremiumPageData(html []byte) (string, []json.RawMessage, error) {
	boardID, err := extractAssignedJSONString(html, premiumBoardIDMarker)
	if err != nil {
		return "", nil, fmt.Errorf("masjidboardlive: extract Premium board ID: %w", err)
	}

	rawRows, err := extractAssignedJSONArray(html, premiumInfoMarker)
	if err != nil {
		return "", nil, fmt.Errorf("masjidboardlive: extract Premium theInfo: %w", err)
	}

	var rows []json.RawMessage
	if err := json.Unmarshal(rawRows, &rows); err != nil {
		return "", nil, fmt.Errorf("masjidboardlive: decode Premium theInfo: %w", err)
	}
	if len(rows) != 29 {
		return "", nil, fmt.Errorf("masjidboardlive: expected 29 Premium rows, got %d", len(rows))
	}
	return boardID, rows, nil
}

func extractAssignedJSONString(source, marker []byte) (string, error) {
	rest, err := assignmentValue(source, marker)
	if err != nil {
		return "", err
	}
	start := bytes.IndexByte(rest, '"')
	if start < 0 {
		return "", fmt.Errorf("assignment has no JSON string")
	}

	value := rest[start:]
	escaped := false
	for i := 1; i < len(value); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch value[i] {
		case '\\':
			escaped = true
		case '"':
			var decoded string
			if err := json.Unmarshal(value[:i+1], &decoded); err != nil {
				return "", fmt.Errorf("decode assigned string: %w", err)
			}
			if decoded == "" {
				return "", fmt.Errorf("assigned string is empty")
			}
			return decoded, nil
		}
	}
	return "", fmt.Errorf("unterminated assigned string")
}

func extractAssignedJSONArray(source, marker []byte) ([]byte, error) {
	rest, err := assignmentValue(source, marker)
	if err != nil {
		return nil, err
	}
	start := bytes.IndexByte(rest, '[')
	if start < 0 {
		return nil, fmt.Errorf("assignment has no array")
	}

	array := rest[start:]
	depth := 0
	inString := false
	escaped := false
	for i, b := range array {
		if inString {
			if escaped {
				escaped = false
				continue
			}
			switch b {
			case '\\':
				escaped = true
			case '"':
				inString = false
			}
			continue
		}

		switch b {
		case '"':
			inString = true
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				result := make([]byte, i+1)
				copy(result, array[:i+1])
				return result, nil
			}
			if depth < 0 {
				return nil, fmt.Errorf("malformed assigned array")
			}
		}
	}
	if inString {
		return nil, fmt.Errorf("unterminated string in assigned array")
	}
	return nil, fmt.Errorf("unterminated assigned array")
}

func assignmentValue(source, marker []byte) ([]byte, error) {
	position := bytes.Index(source, marker)
	if position < 0 {
		return nil, fmt.Errorf("page has no %s assignment", marker)
	}
	rest := source[position+len(marker):]
	equals := bytes.IndexByte(rest, '=')
	if equals < 0 {
		return nil, fmt.Errorf("%s assignment has no equals sign", marker)
	}
	return rest[equals+1:], nil
}
