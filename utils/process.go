package utils

import (
	m "pretifier/model"
	"regexp"
	"strings"
)

var (
	icaoPattern = regexp.MustCompile(`(\*?)##([A-Z]{4})`)
	iataPattern = regexp.MustCompile(`(\*?)#([A-Z]{3})`)
	datePattern = regexp.MustCompile(`D\(([^)]+)\)`)
	t12Pattern  = regexp.MustCompile(`T12\(([^)]+)\)`)
	t24Pattern  = regexp.MustCompile(`T24\(([^)]+)\)`)
)

func ProcessIntinerary(input string, airports map[string]m.AirportRecord) (string, string) {
	input = strings.ReplaceAll(input, "\v", "\n")
	input = strings.ReplaceAll(input, "\f", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	input = strings.ReplaceAll(input, "\r\n", "\n")

	res := input
	resColor := input

	//replace ICAO code first
	res = icaoPattern.ReplaceAllStringFunc(res, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		if len(m) > 0 {
			isCity := m[1] == "*"
			code := m[2]
			if rec, ok := airports[code]; ok {
				if isCity {
					return rec.Municipality
				}
				return rec.Name
			}
		}
		return match
	})

	resColor = icaoPattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		if len(m) > 0 {
			isCity := m[1] == "*"
			code := m[2]
			if rec, ok := airports[code]; ok {
				if isCity {
					return colorMagenta + colorBold + rec.Municipality + colorReset
				}
				return colorGreen + colorBold + rec.Name + colorReset
			}
		}
		return match
	})

	//replace IATA code
	res = iataPattern.ReplaceAllStringFunc(res, func(match string) string {
		m := iataPattern.FindStringSubmatch(match)
		if len(m) > 0 {
			isCity := m[1] == "*"
			code := m[2]
			if rec, ok := airports[code]; ok {
				if isCity {
					return rec.Municipality
				}
				return rec.Name
			}
		}
		return match
	})

	resColor = iataPattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := iataPattern.FindStringSubmatch(match)
		if len(m) > 0 {
			isCity := m[1] == "*"
			code := m[2]
			if rec, ok := airports[code]; ok {
				if isCity {
					return colorMagenta + colorBold + rec.Municipality + colorReset
				}
				return colorGreen + colorBold + rec.Name + colorReset
			}
		}
		return match
	})

	//replace date
	res = datePattern.ReplaceAllStringFunc(res, func(match string) string {
		m := datePattern.FindStringSubmatch(match)
		formatted, ok := formatDate(m[1])
		if ok {
			return formatted
		}
		return match
	})

	resColor = datePattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := datePattern.FindStringSubmatch(match)
		formatted, ok := formatDate(m[1])
		if ok {
			return colorCyan + formatted + colorReset
		}
		return match
	})

	//replace T12
	res = t12Pattern.ReplaceAllStringFunc(res, func(match string) string {
		m := t12Pattern.FindStringSubmatch(match)
		formatted, tz, ok := formatT12(m[1])
		if ok {
			return formatted + "(" + tz + ")"
		}
		return match
	})

	resColor = t12Pattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := t12Pattern.FindStringSubmatch(match)
		formatted, tz, ok := formatT12(m[1])
		if ok {
			return colorYellow + formatted + colorBlue + "(" + tz + ")" + colorReset
		}
		return match
	})

	//replace T24
	res = t24Pattern.ReplaceAllStringFunc(res, func(match string) string {
		m := t24Pattern.FindStringSubmatch(match)
		formatted, tz, ok := formatT24(m[1])
		if ok {
			return formatted + "(" + tz + ")"
		}
		return match
	})

	resColor = t24Pattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := t24Pattern.FindStringSubmatch(match)
		formatted, tz, ok := formatT24(m[1])
		if ok {
			return colorYellow + formatted + colorBlue + "(" + tz + ")" + colorReset
		}
		return match
	})

	return cleanNewlines(res), cleanNewlines(resColor)
}

func formatDate(dateStr string) (string, bool) {
	t, _, ok := ParseISO8601(dateStr)
	if !ok {
		return "", false
	}
	return t.Format("02 Jan 2006"), true
}

func formatT12(dateStr string) (string, string, bool) {
	t, tz, ok := ParseISO8601(dateStr)
	if !ok {
		return "", "", false
	}
	return t.Format("03:04PM "), tz, true
}

func formatT24(dateStr string) (string, string, bool) {
	t, tz, ok := ParseISO8601(dateStr)
	if !ok {
		return "", "", false
	}
	return t.Format("15:04 "), tz, true
}

var toomuchnewlinePattern = regexp.MustCompile(`\n{3,}`)

func cleanNewlines(s string) string {
	return toomuchnewlinePattern.ReplaceAllString(s, "\n")
}
