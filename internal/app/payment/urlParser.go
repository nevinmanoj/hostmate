package payment

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	errMap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	httputil "github.com/nevinmanoj/hostmate/internal/app/httputil"
	payment "github.com/nevinmanoj/hostmate/internal/domain/payment"
)

func parsePaymentFilter(q url.Values) (payment.PaymentFilter, *errMap.BadRequestError) {
	var f payment.PaymentFilter

	if v := q.Get("payment_type"); v != "" {
		paymentType, err := parsePaymentTypeSlice(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "payment_type",
				Reason: err.Error(),
			}
		}
		f.PaymentType = paymentType
	}
	if v := q.Get("from_date"); v != "" {
		fromDate, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "from_date",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.FromDate = fromDate
	}

	if v := q.Get("to_date"); v != "" {
		toDate, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "to_date",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.ToDate = toDate
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

func parsePaymentType(v string) (payment.PaymentType, error) {
	switch strings.ToLower(v) {

	case "UPI":
		return payment.PaymentUPI, nil
	case "cash":
		return payment.PaymentCash, nil
	case "bank-transfer":
		return payment.PaymentBankTransfer, nil
	case "other":
		return payment.PaymentOther, nil
	default:
		return "", fmt.Errorf("Invalid type, must be ['UPI','cash','bank-transfer','other'] ")
	}
}

func parsePaymentTypeSlice(v string) ([]payment.PaymentType, error) {
	parts := strings.Split(v, ",")
	out := make([]payment.PaymentType, 0, len(parts))

	for _, p := range parts {
		t, err := parsePaymentType(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}
