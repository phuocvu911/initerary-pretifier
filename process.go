package main

import (
	m "pretifier/model"
	"regexp"
	"strings"
)

var (
	icaoPattern = regexp.MustCompile(`(\*?)##([A-Z]{4})`)
	iataPattern = regexp.MustCompile(`(\*?)##([A-Z]{3})`)
)

func ProcessIntinerary(input string, airports map[string]m.AirportRecord) (string, string) {
	input = strings.ReplaceAll("\v", "\n")
	input = strings.ReplaceAll("\f", "\n")
	input = strings.ReplaceAll("\r", "\n")
	input = strings.ReplaceAll("\r\n", "\n")

	res := input
	resColor := input

	//replace ICAO code first
	res = icaoPattern.ReplaceAllStringFunc(res, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		isCity := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if isCity {
				return rec.Municipality
			}
			return rec.Name
		}
		return match
	})

	resColor = icaoPattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		isCity := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if isCity {
				return colorMagenta + colorBold + rec.Municipality + colorReset
			}
			return colorGreen + colorBold + rec.Name + colorReset
		}
		return match
	})

	//replace IATA code
	res = iataPattern.ReplaceAllStringFunc(res, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		isCity := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if isCity {
				return rec.Municipality
			}
			return rec.Name
		}
		return match
	})

	resColor = iataPattern.ReplaceAllStringFunc(resColor, func(match string) string {
		m := iataPattern.FindStringSubmatch(match)
		isCity := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if isCity {
				return colorMagenta + colorBold + rec.Municipality + colorReset
			}
			return colorGreen + colorBold + rec.Name + colorReset
		}
		return match
	})
}
