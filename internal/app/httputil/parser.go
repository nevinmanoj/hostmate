package httputil

import (
	"strconv"
	"strings"
	"time"
)

func ParseInt64Slice(v string) ([]int64, error) {
	parts := strings.Split(v, ",")
	out := make([]int64, 0, len(parts))

	for _, p := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	return out, nil
}

const dateLayout = "2006-01-02"

func ParseDatePtr(v string) (*time.Time, error) {
	t, err := time.Parse(dateLayout, v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
