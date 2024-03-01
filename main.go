package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	city      string
	totaltemp int64
	count     int
	minTemp   int
	maxTemp   int
	avg       float32
}

func main() {
	get_average()
}

func get_average() {
	// open file
	f, err := os.Open("./measurements.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// map
	m := make(map[string]int64)

	// read through all rows and print data
	buf := make([]byte, 1024)
	end := ""
	defer f.Close()
	for {
		n, err := f.Read(buf)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		} else if n > 0 {
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
					temp := int64(tempFloat * 10)

					_, ok := m[city]
					if ok {
						m[city] += temp
					} else {
						m[city] = temp
					}
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
					_, ok := m[city]
					if ok {
						m[city] += int64(tempFloat * 10)
					} else {
						m[city] = int64(tempFloat * 10)
					}
				}
			}
			end = rows[len(rows)-1]
		}
	}

	// print all entries in dict
	for k, v := range m {
		fmt.Printf("%s: %d\n", k, v)
	}
}
