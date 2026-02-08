package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	start, err := time.Parse(dateLayout, strings.TrimSpace(dstart))
	if err != nil {
		return "", err
	}

	now = dateOnly(now)
	start = dateOnly(start)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat format")
		}

		// минимум один перенос от dstart
		d := start.AddDate(1, 0, 0)
		for !afterNow(d, now) {
			d = d.AddDate(1, 0, 0)
		}
		return d.Format(dateLayout), nil

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat format")
		}

		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid day interval")
		}
		if n < 1 || n > 400 {
			return "", errors.New("day interval out of range")
		}

		// минимум один перенос от dstart
		d := start.AddDate(0, 0, n)
		for !afterNow(d, now) {
			d = d.AddDate(0, 0, n)
		}
		return d.Format(dateLayout), nil

	default:
		return "", errors.New("unsupported repeat format")
	}
}


func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := strings.TrimSpace(r.FormValue("now"))
	dateStr := strings.TrimSpace(r.FormValue("date"))
	repeatStr := r.FormValue("repeat")

	if dateStr == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(repeatStr) == "" {
		http.Error(w, "repeat is required", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, next)
}


func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
