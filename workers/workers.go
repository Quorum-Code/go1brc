package workers

import (
	"fmt"
	"sync"
	"time"
)

type EntryWorker struct {
	City string
	Job  <-chan float64
	wg   *sync.WaitGroup
}

func StartWorker(ew *EntryWorker) {
	defer ew.wg.Done()

	job := <-ew.Job
	fmt.Println(job)
}

func worker(jobs <-chan float64, counter *int) {
	for range jobs {
		*counter++
	}
}

func Worker_approach() {
	start := time.Now()

	// may need proceeding line for printing values after calcing
	//var wg sync.WaitGroup

	// jobs
	// strings := []string{"sacramento", "arcata", "sacramento"}
	allJobs := make(map[string](chan float64))
	allCounts := make(map[string]*int)

	allJobs["arcata"] = make(chan float64)
	allCounts["arcata"] = new(int)
	go worker(allJobs["arcata"], allCounts["arcata"])

	// create jobs
	allJobs["arcata"] <- 0.0
	allJobs["arcata"] <- 1.0

	time.Sleep(time.Second)

	// read results
	fmt.Printf("Count: %d", *allCounts["arcata"])

	fmt.Printf("Elapsed time: %s", time.Since(start))
}
