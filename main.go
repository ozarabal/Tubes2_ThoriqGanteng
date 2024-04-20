package main 

import (
	"fmt"
	"time"
	"Tubes2_ThoriqGanteng/query"
)

func main() {
	startTime := time.Now()
	fmt.Println("Hello World")

	query.FindTree("Indonesia", "Bekasi")
	// query.GetLinks("https://en.wikipedia.org/wiki/Indonesia")
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("Execution Time:", duration)
}