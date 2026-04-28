package datamodel

import (
	"encoding/json"
	"strings"
	"time"
)

// DateTime is a custom time type that supports multiple date formats during
// JSON unmarshaling: ISO 8601 (with optional trailing "Z"), and common human
// date formats like "11 JAN 26", "27 January 2026", "January 27, 2026".
type DateTime struct {
	time.Time
}

var fallbackDateFormats = []string{
	"02 Jan 06",        // "11 JAN 26"
	"02 Jan 2006",      // "17 JAN 2026"
	"02 January 2006",  // "27 January 2026"
	"January 02, 2006", // "January 27, 2026"
}

// UnmarshalJSON parses a JSON string into a DateTime.
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == "" {
		dt.Time = time.Time{}
		return nil
	}

	// ISO 8601 with trailing Z
	s := strings.Replace(raw, "Z", "+00:00", 1)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		dt.Time = t
		return nil
	}
	// ISO 8601 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		dt.Time = t
		return nil
	}

	for _, fmt := range fallbackDateFormats {
		if t, err := time.Parse(fmt, raw); err == nil {
			dt.Time = t
			return nil
		}
	}

	return &time.ParseError{
		Value:   raw,
		Message: "cannot parse datetime",
	}
}

// MarshalJSON writes the DateTime as an ISO 8601 string with milliseconds and Z suffix.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(dt.Time.UTC().Format("2006-01-02T15:04:05.000Z"))
}
