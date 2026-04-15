package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	// Handle -h flag
	if len(args) == 1 && args[0] == "-h" {
		fmt.Println(usage)
		return
	}

	// Check for optional --color flag
	colorOutput := false
	filtered := args[:0]
	for _, a := range args {
		if a == "--color" {
			colorOutput = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	// Validate argument count
	if len(args) != 3 {
		fmt.Println(usage)
		return
	}

	inputPath := args[0]
	outputPath := args[1]
	lookupPath := args[2]

	// Check input exists
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Input not found")
		return
	}

	// Load airport lookup
	lookupData, err := os.ReadFile(lookupPath)
	if err != nil {
		fmt.Println("Airport lookup not found")
		return
	}

	airports, err := parseAirportLookup(string(lookupData))
	if err != nil {
		fmt.Println("Airport lookup malformed")
		return
	}

	// Process the itinerary
	result, colorResult := processItinerary(string(inputData), airports)

	// Write output
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		fmt.Println("Output path invalid")
		return
	}

	// Bonus: print colored output to stdout (only if --color flag is set)
	if colorOutput {
		fmt.Print(colorResult)
	}
}

// AirportRecord holds parsed airport data
type AirportRecord struct {
	Name         string
	Municipality string
	ICAOCode     string
	IATACode     string
}

func parseAirportLookup(data string) (map[string]*AirportRecord, error) {
	lines := strings.Split(strings.ReplaceAll(data, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	// Parse header row to determine column order (bonus: dynamic column order)
	headerLine := lines[0]
	headers := parseCSVRow(headerLine)

	// Find required column indices
	colIdx := make(map[string]int)
	requiredCols := []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"}
	for i, h := range headers {
		colIdx[strings.TrimSpace(strings.ToLower(h))] = i
	}

	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			return nil, fmt.Errorf("missing column: %s", col)
		}
	}

	airports := make(map[string]*AirportRecord)

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := parseCSVRow(line)
		if len(fields) < 6 { // We expect at least 6 columns based on requiredCols
			return nil, fmt.Errorf("row has too few columns")
		}

		// Check for blank cells in required columns
		for _, col := range requiredCols {
			idx := colIdx[col]
			if idx >= len(fields) || strings.TrimSpace(fields[idx]) == "" {
				return nil, fmt.Errorf("blank cell in column %s", col)
			}
		}

		rec := &AirportRecord{
			Name:         strings.TrimSpace(fields[colIdx["name"]]),
			Municipality: strings.TrimSpace(fields[colIdx["municipality"]]),
			ICAOCode:     strings.TrimSpace(fields[colIdx["icao_code"]]),
			IATACode:     strings.TrimSpace(fields[colIdx["iata_code"]]),
		}

		if rec.ICAOCode != "" {
			airports[rec.ICAOCode] = rec
		}
		if rec.IATACode != "" {
			airports[rec.IATACode] = rec
		}
	}

	return airports, nil
}

// parseCSVRow handles CSV parsing including quoted fields
// redo it using encoding/csv package, but we need to handle the header row separately to determine column indices, so we will use csv.Reader for the whole file and then process the header row to build the column index map. This way we can handle quoted fields and commas properly without reinventing CSV parsing logic.
func parseCSVRow(line string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '"' {
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				current.WriteByte('"')
				i++
			} else {
				inQuotes = !inQuotes
			}
		} else if ch == ',' && !inQuotes {
			fields = append(fields, current.String())
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}
	fields = append(fields, current.String())
	return fields
}

// Regex patterns
var (
	icaoPattern = regexp.MustCompile(`(\*?)##([A-Z]{4})`)
	iataPattern = regexp.MustCompile(`(\*?)#([A-Z]{3})`)
	datePattern = regexp.MustCompile(`D\(([^)]+)\)`)
	t12Pattern  = regexp.MustCompile(`T12\(([^)]+)\)`)
	t24Pattern  = regexp.MustCompile(`T24\(([^)]+)\)`)
)

func processItinerary(input string, airports map[string]*AirportRecord) (string, string) {
	// Normalize line endings: \v \f \r -> \n
	input = strings.ReplaceAll(input, "\v", "\n")
	input = strings.ReplaceAll(input, "\f", "\n")
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	// Process airport codes and dates line by line to preserve structure
	// but we need to replace inline, so process the whole string
	result := input
	colorResult := input

	// Replace ICAO codes first (## before #)
	result = icaoPattern.ReplaceAllStringFunc(result, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		city := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if city {
				return rec.Municipality
			}
			return rec.Name
		}
		return match
	})

	colorResult = icaoPattern.ReplaceAllStringFunc(colorResult, func(match string) string {
		m := icaoPattern.FindStringSubmatch(match)
		city := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if city {
				return colorMagenta + colorBold + rec.Municipality + colorReset
			}
			return colorGreen + colorBold + rec.Name + colorReset
		}
		return match
	})

	// Replace IATA codes
	result = iataPattern.ReplaceAllStringFunc(result, func(match string) string {
		m := iataPattern.FindStringSubmatch(match)
		city := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if city {
				return rec.Municipality
			}
			return rec.Name
		}
		return match
	})

	colorResult = iataPattern.ReplaceAllStringFunc(colorResult, func(match string) string {
		m := iataPattern.FindStringSubmatch(match)
		city := m[1] == "*"
		code := m[2]
		if rec, ok := airports[code]; ok {
			if city {
				return colorMagenta + colorBold + rec.Municipality + colorReset
			}
			return colorGreen + colorBold + rec.Name + colorReset
		}
		return match
	})

	// Replace dates D(...)
	result = datePattern.ReplaceAllStringFunc(result, func(match string) string {
		m := datePattern.FindStringSubmatch(match)
		formatted, ok := formatDate(m[1])
		if !ok {
			return match
		}
		return formatted
	})

	colorResult = datePattern.ReplaceAllStringFunc(colorResult, func(match string) string {
		m := datePattern.FindStringSubmatch(match)
		formatted, ok := formatDate(m[1])
		if !ok {
			return match
		}
		return colorCyan + formatted + colorReset
	})

	// Replace T12(...)
	result = t12Pattern.ReplaceAllStringFunc(result, func(match string) string {
		m := t12Pattern.FindStringSubmatch(match)
		formatted, ok := formatTime12(m[1])
		if !ok {
			return match
		}
		return formatted
	})

	colorResult = t12Pattern.ReplaceAllStringFunc(colorResult, func(match string) string {
		m := t12Pattern.FindStringSubmatch(match)
		formatted, ok := formatTime12(m[1])
		if !ok {
			return match
		}
		// Highlight time in yellow, offset in blue
		parts := strings.SplitN(formatted, " ", 2)
		if len(parts) == 2 {
			return colorYellow + parts[0] + colorReset + " " + colorBlue + parts[1] + colorReset
		}
		return colorYellow + formatted + colorReset
	})

	// Replace T24(...)
	result = t24Pattern.ReplaceAllStringFunc(result, func(match string) string {
		m := t24Pattern.FindStringSubmatch(match)
		formatted, ok := formatTime24(m[1])
		if !ok {
			return match
		}
		return formatted
	})

	colorResult = t24Pattern.ReplaceAllStringFunc(colorResult, func(match string) string {
		m := t24Pattern.FindStringSubmatch(match)
		formatted, ok := formatTime24(m[1])
		if !ok {
			return match
		}
		parts := strings.SplitN(formatted, " ", 2)
		if len(parts) == 2 {
			return colorYellow + parts[0] + colorReset + " " + colorBlue + parts[1] + colorReset
		}
		return colorYellow + formatted + colorReset
	})

	// Trim vertical whitespace: no more than one consecutive blank line
	result = trimExcessiveBlankLines(result)
	colorResult = trimExcessiveBlankLines(colorResult)

	return result, colorResult
}

// ISO 8601 datetime with offset: 2007-04-05T12:30−02:00
// The offset can use ASCII minus (-) or Unicode minus (−)
var iso8601Pattern = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})([+\-−](\d{2}):(\d{2})|Z)$`)

// change it using layout and parse func in time, because time in flight never need second and we can leverage that.
func parseISO8601(s string) (t time.Time, offsetStr string, ok bool) {
	// Normalize Unicode minus sign to ASCII minus
	s = strings.ReplaceAll(s, "\u2212", "-")

	m := iso8601Pattern.FindStringSubmatch(s)
	if m == nil {
		return t, "", false
	}

	year := m[1]
	month := m[2]
	day := m[3]
	hour := m[4]
	min := m[5]
	offsetFull := m[6]

	if offsetFull == "Z" {
		offsetStr = "(+00:00)"
	} else {
		// Validate offset format: ±HH:MM — must be exactly ±02:00 style
		// The sign + 2-digit hour + colon + 2-digit min
		sign := string(offsetFull[0])
		rest := offsetFull[1:]
		parts := strings.Split(rest, ":")
		if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
			return t, "", false
		}
		offsetStr = "(" + sign + parts[0] + ":" + parts[1] + ")"
	}

	// Parse the time using RFC3339-like approach
	// Build a proper RFC3339 string
	var rfc string
	if offsetFull == "Z" {
		rfc = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", year, month, day, hour, min)
	} else {
		rfc = fmt.Sprintf("%s-%s-%sT%s:%s:00%s", year, month, day, hour, min, offsetFull)
	}
	// Normalize back (minus was already normalized)
	t, err := time.Parse(time.RFC3339, rfc)
	if err != nil {
		return t, "", false
	}

	return t, offsetStr, true
}

func formatDate(s string) (string, bool) {
	t, _, ok := parseISO8601(s)
	if !ok {
		return "", false
	}
	return t.Format("02 Jan 2006"), true
}

func formatTime12(s string) (string, bool) {
	t, offsetStr, ok := parseISO8601(s)
	if !ok {
		return "", false
	}
	timeStr := t.Format("3:04PM")
	// Ensure two-digit minutes if needed - time.Format handles it
	return timeStr + " " + offsetStr, true
}

func formatTime24(s string) (string, bool) {
	t, offsetStr, ok := parseISO8601(s)
	if !ok {
		return "", false
	}
	timeStr := t.Format("15:04")
	return timeStr + " " + offsetStr, true
}

func trimExcessiveBlankLines(s string) string {
	// Replace more than 2 consecutive newlines with exactly 2
	re := regexp.MustCompile(`\n{3,}`)
	return re.ReplaceAllString(s, "\n\n")
}
