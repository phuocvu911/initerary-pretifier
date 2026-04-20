# Itinerary Prettifier

A command-line tool that converts administrator-formatted flight itineraries into customer-friendly documents.

## Overview

Anywhere Holidays back-office administrators receive raw flight itineraries containing airport codes, ISO 8601 timestamps, and other technical formatting. This tool automates the conversion into readable, customer-friendly text.

## Requirements

- Go 1.25.0
- airport-lookup.csv file
- input file (in the form of flight itinerary)

## Installation

```bash
git clone https://gitea.kood.tech/hoangphuocvu/prettifier.git
cd prettifier
go build -o prettifier .
```

Or run directly without building:

```bash
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

## Usage

```
go run . <input> <output> <airport-lookup> [--color]
```

**Arguments:**

| Argument | Description |
|---|---|
| `input` | Path to the raw itinerary text file |
| `output` | Path where the prettified itinerary will be written |
| `airport-lookup` | Path to the airport codes CSV file |
| `--color` | Bonus Feature. Print colour-highlighted output to stdout in addition to writing the file |

**Help flag:**

```bash
go run . -h
```

## Features

### Airport Code Conversion

The tool replaces encoded airport codes with human-readable names:

| Input | Output |
|---|---|
| `##EGLL` | `London Heathrow Airport` (ICAO 4-letter code) |
| `#LAX` | `Los Angeles International Airport` (IATA 3-letter code) |


If a code is not found in the lookup CSV, it is left unchanged.

### Date & Time Formatting

ISO 8601 timestamps are converted to friendly formats:

| Input | Output |
|---|---|
| `D(2007-04-05T12:30-02:00)` | `05 Apr 2007` |
| `T12(2007-04-05T12:30-02:00)` | `12:30PM (-02:00)` |
| `T24(2007-04-05T12:30-02:00)` | `12:30 (-02:00)` |
| `T24(2007-04-05T12:30Z)` | `12:30 (+00:00)` |

Malformed date/time tokens are left unchanged. For example:
- `T13(...)` — invalid clock type, unchanged
- `T12(2007-04-05T12:30-2:00)` — malformed offset (must be `±HH:MM`), unchanged

### Whitespace Trimming

- Vertical whitespace characters (`\v`, `\f`, `\r`) are converted to newlines (`\n`)
- No more than one consecutive blank line appears in the output

### Error Handling

| Condition | Message |
|---|---|
| Wrong number of arguments | Displays usage |
| Input file not found | `Input not found` |
| Airport lookup file not found | `Airport lookup not found` |
| Malformed airport CSV data | `Airport lookup malformed` |

In all error cases, no output file is created or overwritten.

## Airport Lookup CSV Format

The CSV file must include the following columns (in any order):

- `name` — full airport name
- `iso_country` — ISO country code
- `municipality` — city/town name
- `icao_code` — 4-letter ICAO code
- `iata_code` — 3-letter IATA code
- `coordinates` — lat/lon coordinates

The first row must be a header row. Any missing or blank required column will cause a malformed error.

## Bonus Features

### Coloured stdout output

Pass the optional `--color` flag to also print the prettified itinerary to stdout with ANSI colour highlighting:

```bash
go run . ./input.txt ./output.txt ./airport-lookup.csv --color
```

- **Airport names** — bold green
- **City names** — bold magenta
- **Dates** — cyan
- **Times** — yellow
- **Timezone offsets** — blue

By default the program produces no stdout output, in line with the spec. The `--color` flag is strictly opt-in.

### City names

Prefix any airport code with `*` to get the city/municipality name instead of the airport name:

| Input | Output |
|---|---|
| `*##EGLL` | `London` (city name instead of airport name) |
| `*#LAX` | `Los Angeles` (city name instead of airport name) |

### Dynamic column order

The airport lookup CSV columns can appear in any order. The tool reads the header row at runtime to determine column positions, so reordered CSVs work without any code changes.

### Extra columns tolerance

The lookup CSV may contain more columns than the six required ones. Blank cells in those extra columns are silently ignored — only the required columns (`name`, `iso_country`, `municipality`, `icao_code`, `iata_code`, `coordinates`) are checked for blank values. This means real-world CSV exports with additional fields like `type`, `elevation_ft`, etc. work out of the box, even if some of those extra fields are empty (does not applicable if extra column fall in between required columns).

## Extra tool to shuffle columns
Based on provided lookup CSV, the tool locating in ./tool can be used and configured to generate dynamic column order or add extra columns to mimic the real world data.

## Examples

**Input:**
```
Departure: ##KLAX - T12(2007-04-05T09:30-07:00)
Arrival: ##EGLL - T12(2007-04-06T06:15+01:00)
Departing from *##KLAX
```

**Output:**
```
Departure: Los Angeles International Airport - 9:30AM (-07:00)
Arrival: London Heathrow Airport - 6:15AM (+01:00)
Departing from Los Angeles
```

## Design Notes

- ICAO codes (`##`) are processed before IATA codes (`#`) to avoid the double-hash being partially matched as a single hash.
- Airport codes use uppercase only, matching standard airline industry conventions.
- The tool reads the entire input into memory, processes it, then writes output — no partial writes on error.
