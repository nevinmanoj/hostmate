package httputil

import (
	"fmt"
	"strconv"
)

func StringsToInt64s(vals []string) ([]int64, error) {
	res := make([]int64, 0, len(vals))

	for _, v := range vals {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid int64 value %q: %w", v, err)
		}
		res = append(res, i)
	}

	return res, nil
}
