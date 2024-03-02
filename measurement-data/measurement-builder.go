package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	cities := []string{"Sacramento", "Arcata", "Eureka", "Trinidad", "San Fransisco", "Oakland", "Monterey"}

	fmt.Println("Creating file...")
	f, err := os.Create("./measurements.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Creating entries...")
	for i := 0; i < 10000; i++ {
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
}
