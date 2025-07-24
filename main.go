package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var measurementReadyLogRe = regexp.MustCompile(`^[^ ]+ +Measurement ready \((.*)\)$`)

type datapoint struct {
	Seconds int
	T       float32
	RH      float32
	P       float32
	CO2     float32
}

func processLog(logline string) (*datapoint, error) {
	var d datapoint

	matches := measurementReadyLogRe.FindStringSubmatch(logline)
	if len(matches) == 0 {
		return nil, nil
	}
	data := matches[1]

	fields := strings.Split(data, ", ")
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed field in log: %q", field)
		}

		name := parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		switch name {
		case "seconds":
			d.Seconds = value
		case "t":
			d.T = float32(value) / 10
		case "rh":
			d.RH = float32(value) / 2
		case "p":
			d.P = float32(value) / 10
		case "co2":
			d.CO2 = float32(value)
		default:
			continue // ignore unknown fields
		}
	}

	return &d, nil
}

func run() error {
	enc := json.NewEncoder(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		logline := strings.TrimSpace(scanner.Text())
		fmt.Fprintln(os.Stderr, logline)
		d, err := processLog(logline)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse datapoint: %v\n", err.Error())
			continue
		}
		if d == nil {
			continue
		}
		enc.Encode(d)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("bufio.Scanner.Err: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ag-data-logs: %s\n", err.Error())
		os.Exit(1)
	}
}
