package workers

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	totaltemp int64
	count     int
	mintemp   float32
	maxtemp   float32
}

func chunk_worker(chunks <-chan []string, m map[string](*Entry), wg *sync.WaitGroup) {
	defer wg.Done()
	for chunk := range chunks {
		for i := 0; i < len(chunk)-1; i++ {
			split := strings.Split(chunk[i], ":")

			if len(split) == 2 {
				city := split[0]
				temp, err := strconv.ParseFloat(split[1], 64)
				if err != nil {
					fmt.Printf("Error: %s, %f", err.Error(), temp)
					return
				}

				entry, ok := m[city]
				if !ok {
					new_entry := Entry{
						totaltemp: 0,
						count:     0,
						mintemp:   1000,
						maxtemp:   -1000,
					}
					m[city] = &new_entry
					entry = m[city]
				}

				entry.totaltemp += int64(temp * 10.0)

				entry.count++

				if entry.mintemp > float32(temp) {
					entry.mintemp = float32(temp)
				}

				if entry.maxtemp < float32(temp) {
					entry.maxtemp = float32(temp)
				}
			}
		}
	}
}

// Elapsed time: 35.8378618s (4096 byte buffer, 16 workers)
func Worker_map_approach() {
	start := time.Now()

	// open file
	f, err := os.Open("./measurements.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// Create string buffer for rows of data
	buf := make([]byte, 4096)
	end := ""

	// Best return on investment seems to be about 16 workers
	const numworkers = 16

	chunks := make(chan []string)
	maps := make([]map[string]*Entry, numworkers)
	var wg sync.WaitGroup

	// startup workers
	for i := 0; i < numworkers; i++ {
		wg.Add(1)
		maps[i] = make(map[string]*Entry)
		go chunk_worker(chunks, maps[i], &wg)
	}

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

			if len(rows)-1 > 0 {
				end = rows[len(rows)-1]
			} else {
				end = ""
			}

			// Push rows to chunk_workers
			chunks <- rows
		}
	}
	// Close channel so workers know there cannot be anymore data sent
	close(chunks)
	fmt.Println("Waiting for workers...")
	wg.Wait()

	// Aggregate worker information
	for i := 1; i < numworkers; i++ {
		for k, v := range maps[i] {
			entry, ok := maps[0][k]
			if !ok {
				*entry = *v
			} else {
				entry.count += v.count
				entry.totaltemp += v.totaltemp

				if entry.mintemp > v.mintemp {
					entry.mintemp = v.mintemp
				}

				if entry.maxtemp < v.maxtemp {
					entry.maxtemp = v.maxtemp
				}
			}
		}
	}

	// Print data
	if len(maps) <= 0 {
		return
	}
	totalCount := int64(0)
	for k, v := range maps[0] {
		avg := float64(v.totaltemp / int64(v.count) / 10)
		fmt.Printf("%s Avg: %.1f Count: %d Min: %.1f, Max: %.1f\n", k, avg, v.count, v.mintemp, v.maxtemp)
		totalCount += int64(v.count)
	}

	fmt.Printf("Total Count: %d\n", totalCount)
	fmt.Printf("Elapsed time: %s\n", time.Since(start))
}
