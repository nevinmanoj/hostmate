package property

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	errMap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
)

func parsePropertyFilter(q url.Values) (property.PropertyFilter, *errMap.BadRequestError) {
	var f property.PropertyFilter

	if v := q.Get("type"); v != "" {
		types, err := parsePropertyTypeSlice(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "type",
				Reason: err.Error(),
			}
		}
		f.Type = types
	}

	if v := q.Get("active"); v != "" {
		active, err := strconv.ParseBool(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "active",
				Reason: err.Error(),
			}
		}
		f.Active = &active
	}

	// Pagination defaults
	f.Limit = 100
	f.Offset = 0

	if v := q.Get("limit"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "limit",
				Reason: err.Error(),
			}
		} else if limit > 0 && limit < 100 {
			f.Limit = limit
		}
	}

	if v := q.Get("offset"); v != "" {
		offset, err := strconv.Atoi(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "offset",
				Reason: err.Error(),
			}
		} else if offset > 0 {
			f.Offset = offset
		}
	}

	return f, nil
}

func ParsePropertyType(v string) (property.PropertyType, error) {
	switch strings.ToLower(v) {
	case "apartment":
		return property.PropertyApartment, nil
	case "villa":
		return property.PropertyVilla, nil
	case "room":
		return property.PropertyRoom, nil
	default:
		return "", fmt.Errorf("Invalid type, must be ['apartment','villa','room'] ")
	}
}

func parsePropertyTypeSlice(v string) ([]property.PropertyType, error) {
	parts := strings.Split(v, ",")
	out := make([]property.PropertyType, 0, len(parts))

	for _, p := range parts {
		t, err := ParsePropertyType(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}
