package fastunsafeurl

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedURL  = errors.New("malformed url")
	ErrMissingScheme = fmt.Errorf("missing scheme %w", ErrMalformedURL)
)

// ParseSchemeHost will return scheme://[credential@]host from the endpoint.
// The endpoint must contain both scheme and host.
// This function does not ensure the validity of the endpoint, thus can produce unexpected results when the endpoint is invalid.
func ParseSchemeHost(endpoint string) (string, int, error) {
	slashCount := 0
	colonPos := -1
	for pos, char := range endpoint {
		switch char {
		case ':':
			// Only process the first colon.
			if colonPos >= 0 {
				break
			}
			colonPos = pos
			// Must have at least 1 char for the scheme and 1 char for the host.
			if colonPos < 1 || colonPos+3 >= len(endpoint) {
				return "", 0, fmt.Errorf("error parse url at %d: %s %w", pos, endpoint, ErrMalformedURL)
			}
			// Make sure the colon is followed by two slashes.
			if endpoint[colonPos+1] != '/' {
				return "", 0, fmt.Errorf("error parse url at %d: %s %w", pos+1, endpoint, ErrMissingScheme)
			}
			if endpoint[colonPos+2] != '/' {
				return "", 0, fmt.Errorf("error parse url at %d: %s %w", pos+2, endpoint, ErrMissingScheme)
			}
		case '/':
			// Colon must come before slash.
			if colonPos < 1 {
				return "", 0, fmt.Errorf("error parse url at %d: %s %w", pos, endpoint, ErrMalformedURL)
			}
			slashCount++
			if slashCount > 2 {
				return endpoint[0:pos], colonPos, nil
			}
		}
	}
	if colonPos < 1 || slashCount < 2 {
		return "", 0, fmt.Errorf("error parse url at: %s %w", endpoint, ErrMissingScheme)
	}
	return endpoint, colonPos, nil
}

// ParseHost will return host and optional credential from the endpoint.
// The endpoint must contain both scheme and host.
// This function does not ensure the validity of the endpoint, thus can produce unexpected results when the endpoint is invalid.
func ParseHost(endpoint string) (string, error) {
	h, colonPos, err := ParseSchemeHost(endpoint)
	if err != nil {
		return "", err
	}
	return h[colonPos+3:], nil
}
