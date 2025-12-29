package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
)

type measurements struct {
	city  string
	min   float64
	max   float64
	count int
	sum   float64
}

const (
	scale2 float64 = 100
	scale1 float64 = 10

	// format string for _most_ of the output
	outputFormat string = "%s=%.1f/%.1f/%.1f, "
)

var separator []byte = []byte(";")

func main() {
	// TOOD: use pointer?
	if os.Getenv("PROFILE") != "" {
		f, err := os.Create("cpu.profile")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	results := make(map[uint64]*measurements, 10_000)

	// what if we just read the whole file into memory?!
	contents, err := os.ReadFile("measurements.txt")
	if err != nil {
		log.Fatal("could not read file\n", err)
	}
	processChunk(contents, results)
	formatResults(results)

}

func processChunk(chunk []byte, results map[uint64]*measurements) {
	start := 0
	for i := range len(chunk) {
		if chunk[i] == '\n' {
			// Process line from start to i (exclusive)
			if i > start {
				processLine(chunk[start:i], results)
			}
			start = i + 1
		}
	}
	// Process last line if file doesn't end with newline
	if start < len(chunk) {
		processLine(chunk[start:], results)
	}
}

// parseTemp takes the value like 12.4 and converts it into 124 so we work with integers, not floats
func parseTemp(tempBytes []byte) float64 {
	// from: https://github.com/benhoyt/go-1brc/blob/bb0641a5086474d0640996e8a2aefe9721e6d814/r3.go#L38-L54
	// note the "float64" is a clever way to go from bytes to int because of ASCII encoding
	// we check if the first character is a negative symbol
	negative := false
	index := 0
	if tempBytes[index] == '-' {
		index++
		negative = true
	}
	temp := float64(tempBytes[index] - '0') // parse first digit
	index++
	if tempBytes[index] != '.' {
		temp = temp*10 + float64(tempBytes[index]-'0') // parse optional second digit
		index++
	}
	index++                                    // skip '.'
	temp += float64(tempBytes[index]-'0') / 10 // parse decimal digit
	if negative {
		temp = -temp
	}
	return float64(temp)
}

func hashString(s []byte) uint64 {
	var h uint64 = 0
	for _, b := range s {
		h = h*31 + uint64(b)
	}
	return h
}

func processLine(line []byte, results map[uint64]*measurements) {
	city, temp, _ := bytes.Cut(line, separator)
	// turn the temp into a number
	tempVal := parseTemp(temp)
	hash := hashString(city)
	m, ok := results[hash]
	if !ok {
		results[hash] = &measurements{
			city:  string(city),
			min:   tempVal,
			max:   tempVal,
			sum:   tempVal,
			count: 1,
		}
	} else {
		if tempVal < m.min {
			m.min = tempVal
		}
		if tempVal > m.max {
			m.max = tempVal
		}
		m.sum += tempVal
		m.count++
	}
}

func formatResults(results map[uint64]*measurements) {
	// just iterate and print, will need to format correctly
	// we need to sort the cities
	cities := make([]string, 0, len(results))
	for _, m := range results {
		cities = append(cities, m.city)
	}
	sort.Strings(cities)

	// Create a reverse lookup from city name to measurement
	cityToMeasurement := make(map[string]*measurements, len(results))
	for _, m := range results {
		cityToMeasurement[m.city] = m
	}

	fmt.Print("{")
	for idx, city := range cities {
		m := cityToMeasurement[city]
		if idx == len(cities)-1 {
			fmt.Printf("%s=%.1f/%.1f/%.1f", city, m.min, round(float64(m.sum)/float64(m.count)), m.max)
		} else {
			fmt.Printf(outputFormat, city, m.min, round(float64(m.sum)/float64(m.count)), m.max)
		}
	}
	fmt.Println("}")

}

func round(x float64) float64 {
	// TODO: this is a bit yikes, it's surely going to be a bit slow
	intermediate := math.Round(x*scale2) / scale2
	return math.Round(intermediate*scale1) / scale1
}
