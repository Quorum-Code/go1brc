package workers

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	totaltemp int64
	count     int
	mintemp   float32
	maxtemp   float32
}

func worker(city string, jobs <-chan string, entry *Entry) {
	for range jobs {
		temp, err := strconv.ParseFloat(<-jobs, 64)

		if err != nil {
			fmt.Print(err.Error())
			return
		}

		// Add total
		entry.totaltemp += int64(temp * 10)

		// Inc count
		entry.count++

		// Check min temp
		if temp < float64(entry.mintemp) {
			entry.mintemp = float32(temp)
		}

		// Check max temp
		if temp > float64(entry.maxtemp) {
			entry.maxtemp = float32(temp)
		}
	}
}

func Worker_map_approach() {
	start := time.Now()

	// jobs
	allJobs := make(map[string](chan string))
	allEntries := make(map[string]*Entry)

	// open file
	f, err := os.Open("./measurements.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Create string buffer for rows of data
	buf := make([]byte, 1024)
	end := ""

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

			if len(rows) > 0 {
				rows[0] = end + rows[0]
			}

			// Iterate through all rows except last
			for i := 0; i < len(rows)-1; i++ {
				split := strings.Split(rows[i], ":")
				job, ok := allJobs[split[0]]

				// No worker/job created yet
				if !ok {
					allJobs[split[0]] = make(chan string)
					allEntries[split[0]] = new(Entry)

					entry := allEntries[split[0]]
					entry.count = 0
					entry.totaltemp = 0
					entry.maxtemp = -1000
					entry.mintemp = 1000

					go worker(split[0], allJobs[split[0]], allEntries[split[0]])

					job = allJobs[split[0]]
				}

				job <- split[1]
			}

			if len(rows) > 1 {
				end = rows[len(rows)-1]
			}
		}
	}

	for k, v := range allEntries {
		fmt.Printf("%s Total Temp: %d Count: %d MinTemp:%.1f MaxTemp: %1.f\n", k, v.totaltemp, v.count, v.mintemp, v.maxtemp)
	}

	fmt.Printf("Elapsed time: %s", time.Since(start))
}
