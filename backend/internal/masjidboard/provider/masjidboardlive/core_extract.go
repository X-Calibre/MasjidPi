package masjidboardlive

import (
	"bytes"
	"fmt"
)

var coreDataMarker = []byte("let data")

// ExtractCoreData locates the JavaScript object assigned to `let data` in a
// public MasjidBoard Live Core board page and returns only the object bytes.
//
// A small scanner is used instead of a regular expression so braces contained
// in quoted JavaScript strings do not terminate the object early.
func ExtractCoreData(html []byte) ([]byte, error) {
	marker := bytes.Index(html, coreDataMarker)
	if marker < 0 {
		return nil, fmt.Errorf("masjidboardlive: Core page has no data assignment")
	}

	rest := html[marker+len(coreDataMarker):]
	equals := bytes.IndexByte(rest, '=')
	if equals < 0 {
		return nil, fmt.Errorf("masjidboardlive: Core data assignment has no equals sign")
	}

	rest = rest[equals+1:]
	start := bytes.IndexByte(rest, '{')
	if start < 0 {
		return nil, fmt.Errorf("masjidboardlive: Core data assignment has no object")
	}

	object := rest[start:]
	depth := 0
	inString := false
	escaped := false

	for i, b := range object {
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
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				result := make([]byte, i+1)
				copy(result, object[:i+1])
				return result, nil
			}
			if depth < 0 {
				return nil, fmt.Errorf("masjidboardlive: malformed Core data object")
			}
		}
	}

	if inString {
		return nil, fmt.Errorf("masjidboardlive: unterminated string in Core data object")
	}
	return nil, fmt.Errorf("masjidboardlive: unterminated Core data object")
}
