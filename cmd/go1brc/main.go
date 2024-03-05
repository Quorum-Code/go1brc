package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"src/workers"
	"strconv"
	"strings"
	"testing"
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
	// p, err := os.Create("sequential.prof")
	// if err != nil {
	// 	return
	// }
	// pprof.StartCPUProfile(p)
	// defer pprof.StopCPUProfile()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("[1] Sequential Approach")
		fmt.Println("[2] Worker Pool Approach")
		fmt.Println("[Q] Quit")
		fmt.Print("Use: ")
		scanner.Scan()
		input := scanner.Text()

		if input == "1" {
			SimpleAverage()
		} else if input == "2" {
			workers.WorkerMapApproach()
		} else if input == "Q" {
			break
		} else {
			fmt.Printf("Error: '%s' is invalid input...", input)
		}
	}
}

func TestSimpleAverage(b *testing.T) {
	SimpleAverage()
}

func BenchmarkWorkers(b *testing.B) {
	workers.WorkerMapApproach()
}

// Time to process 1B rows in 2m30.2183342s
func SimpleAverage() {
	start := time.Now()
	fmt.Println("Starting")

	// map
	m := make(map[string]*Entry)
	addEntry := buildEntryFunc(m)

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
				addEntry(end + rows[0])
				end = ""
			}

			// Iterate through all rows except last
			for ; i < len(rows)-1; i++ {
				addEntry(rows[i])
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
	fmt.Printf("Total time: %s\n", time.Since(start))
}

func buildEntryFunc(m map[string]*Entry) func(string) {
	return func(row string) {
		split := strings.Split(row, ":")
		city := ""
		temperature, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		entry, ok := m[city]

		if !ok {
			fmt.Println("Generating new city entry")
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
