package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type measurements struct {
	min   float64 // min, max and sum are all 10x
	max   float64
	count int
	sum   float64
}

func main() {
	// simplest (but hopefully memory efficient) implementation

	// initialize a map, we know a max of 10K stations
	// TOOD: use pointer?
	results := make(map[string]measurements, 10_000)

	// start rolling through the file!
	f, err := os.Open("measurements.txt")
	if err != nil {
		log.Fatal("could not open file to read\n", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	sep := ";"
	for scanner.Scan() {
		// parse the line
		// TODO: are bytes really better
		city, temp, _ := strings.Cut(scanner.Text(), sep)
		m := results[city]
		// turn the temp into a number
		tempVal, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			log.Fatalf("couldn't parse number city %s tempval %s, error %v", city, temp, err)
		}
		if m.count == 0 {
			m.min = tempVal
			m.max = tempVal
		} else {
			// TODO: is this slow?
			if tempVal < m.min {
				m.min = tempVal
			}
			if tempVal > m.max {
				m.max = tempVal
			}
		}
		m.sum += tempVal
		m.count++
		results[city] = m
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("unexpected error on scanning: %v", err)
	}
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
		formatter := "%s=%.1f/%.1f/%.1f, "
		if idx == len(cities)-1 {
			formatter = "%s=%.1f/%.1f/%.1f"
		}
		if os.Getenv("DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "DEBUG %s %#v", city, results[city])
		}
		fmt.Printf(formatter, city, float64(results[city].min), round(results[city].sum/float64(results[city].count)), float64(results[city].max))
	}
	fmt.Println("}")
}

func round(x float64) float64 {
	// TODO: this is a bit yikes, it's surely going to be a bit slow
	scale := math.Pow(10, float64(2))
	intermediate := math.Round(x*scale) / scale
	scale = math.Pow(10, float64(1))
	return math.Round(intermediate*scale) / scale
}
