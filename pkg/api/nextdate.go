package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	d, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	rule := parts[0]

	switch rule {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат для правила y")
		}
		d = d.AddDate(1, 0, 0)
		for !d.After(now) {
			d = d.AddDate(1, 0, 0)
		}
		return d.Format(DateFormat), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат для правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимое количество дней: %s", parts[1])
		}
		d = d.AddDate(0, 0, days)
		for !d.After(now) {
			d = d.AddDate(0, 0, days)
		}
		return d.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила")
	}
}
