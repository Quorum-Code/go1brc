package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	city      string
	totaltemp int64
	count     int
	minTemp   float32
	maxTemp   float32
	avg       float32
}

func main() {
	simple_get_average()
}

func simple_get_average() {
	start := time.Now()
	fmt.Println("Starting")

	// map
	m := make(map[string]*Entry)
	addEntry := buildLogFunc(m)

	// open file
	f, err := os.Open("./measurements.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// read through all rows and print data
	buf := make([]byte, 1024)
	end := ""
	defer f.Close()

	// iterate file contents through buf
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		} else if n > 0 {
			// content to process
			rows := strings.Split(string(buf[:n]), "\n")
			i := 0

			// Process first row manually if end has some text
			if end != "" && len(rows) > 0 {
				i = 1

				// process
				split := strings.Split(end+rows[0], ":")

				if len(split) == 2 {
					city := split[0]

					tempFloat, err := strconv.ParseFloat(split[1], 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					addEntry(city, tempFloat)
				}
				end = ""
			}

			// Iterate through all rows except last
			for ; i < len(rows)-1; i++ {
				row := strings.Split(rows[i], ":")

				city := row[0]
				if len(row) > 1 {
					tempFloat, err := strconv.ParseFloat(row[1], 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					// Check if key already exists
					addEntry(city, tempFloat)
				}
			}

			if len(rows) > 1 {
				end = rows[len(rows)-1]
			}
		}
	}

	// print all entries in dict
	totalCount := 0
	for k := range m {
		entry := m[k]
		entry.avg = float32(entry.totaltemp) / float32(entry.count) / 10.0
		fmt.Printf("City: %s \tAvg: %.1f \tMin: %.1f \tMax: %.1f \tCount:%d\n", entry.city, entry.avg, entry.minTemp, entry.maxTemp, entry.count)
		totalCount += entry.count
	}
	fmt.Printf("Total Count: %d\n", totalCount)
	fmt.Printf("Total time: %s", time.Since(start))
}

func buildLogFunc(m map[string]*Entry) func(string, float64) {
	return func(city string, temperature float64) {
		entry, ok := m[city]

		if !ok {
			fmt.Println("Generating entry")
			m[city] = &Entry{
				city:      city,
				totaltemp: 0,
				count:     0,
				minTemp:   float32(temperature),
				maxTemp:   float32(temperature),
				avg:       0.0,
			}
			entry = m[city]
		}
		entry.totaltemp = entry.totaltemp + int64(temperature*10)
		entry.count++

		if entry.minTemp > float32(temperature) {
			entry.minTemp = float32(temperature)
		}

		if entry.maxTemp < float32(temperature) {
			entry.maxTemp = float32(temperature)
		}
	}
}
