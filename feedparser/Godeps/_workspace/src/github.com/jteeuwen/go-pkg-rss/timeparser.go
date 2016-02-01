package feeder

import (
	"strings"
	"time"
)

func parseTime(formatted string) (time.Time, error) {
	var layouts = [...]string{
		"Mon, _2 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		"Mon, 2, Jan 2006 15:4",
		"02 Jan 2006 15:04:05 MST",
	}
	var t time.Time
	var err error
	formatted = strings.TrimSpace(formatted)
	for _, layout := range layouts {
		t, err = time.Parse(layout, formatted)
		if !t.IsZero() {
			break
		}
	}
	return t, err
}
