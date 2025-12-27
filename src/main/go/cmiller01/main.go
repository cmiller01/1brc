package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
)

type measurements struct {
	min   int64
	max   int64
	count int
	sum   int64
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
	results := make(map[string]*measurements, 10_000)

	// what if we just read the whole file into memory?!
	contents, err := os.ReadFile("measurements.txt")
	if err != nil {
		log.Fatal("could not read file\n", err)
	}
	processChunk(contents, results)
	formatResults(results)

}

func processChunk(chunk []byte, results map[string]*measurements) {
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
func parseTemp(temp []byte) int64 {
	tempS := bytes.SplitN(temp, []byte("."), 2)
	whole, _ := strconv.Atoi(string(tempS[0]))
	whole *= 10
	frac, _ := strconv.Atoi(string(tempS[1]))
	return int64(whole + frac)
}

func processLine(line []byte, results map[string]*measurements) {
	city, temp, _ := bytes.Cut(line, separator)
	cityS := string(city)
	// turn the temp into a number
	tempVal := parseTemp(temp)
	m, ok := results[cityS]
	if !ok {
		results[cityS] = &measurements{
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

func formatResults(results map[string]*measurements) {
	// just iterate and print, will need to format correctly
	// we need to sort the cities
	cities := make([]string, len(results))
	idx := 0
	for city := range results {
		cities[idx] = city
		idx++
	}
	sort.Strings(cities)
	fmt.Print("{")
	for idx, city := range cities {
		if idx == len(cities)-1 {
			fmt.Printf("%s=%.1f/%.1f/%.1f", city, float64(results[city].min)/10.0, round(float64(results[city].sum)/float64(results[city].count)), float64(results[city].max)/10.0)
		} else {
			fmt.Printf(outputFormat, city, float64(results[city].min)/10.0, round(float64(results[city].sum)/float64(results[city].count)), float64(results[city].max)/10.0)
		}
	}
	fmt.Println("}")

}

func round(x float64) float64 {
	// TODO: this is a bit yikes, it's surely going to be a bit slow
	intermediate := math.Round(x*scale2) / scale2
	return math.Round(intermediate*scale1) / scale1
}
