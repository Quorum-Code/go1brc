package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	power := 0
	for power == 0 {
		fmt.Print("Generate rows equal to 10 to the power of: ")
		scanner.Scan()
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("ERROR: %s", err)
		} else if num >= 10 || num <= 0 {
			fmt.Println("Please enter a number between 0 and 10 (Exclusive)")
		} else {
			fmt.Printf("power: %d", power)
			power = num
		}
	}

	start := time.Now()
	cities := []string{"Sacramento", "Arcata", "Eureka", "Trinidad", "San Fransisco", "Oakland", "Monterey"}
	fmt.Println("Creating file...")
	f, err := os.Create("./measurements.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	count := int32(math.Pow10(power))
	fmt.Println("Creating entries...")
	for i := int32(0); i < count; i++ {
		// get random city
		city := cities[rand.IntN(len(cities))]

		// get random val
		temp := rand.Float64()*199.9 - 100.0

		_, err := f.WriteString(fmt.Sprintf("%v:%.1f\n", city, temp))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Finished creating measurements.txt...")
	fmt.Printf("Elapsed time: %s", time.Since(start))
}
