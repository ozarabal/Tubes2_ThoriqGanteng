package main 

import (
	"fmt"
	"time"
	"Tubes2_ThoriqGanteng/query"
)

func main() {
	startTime := time.Now()
	fmt.Println("Hello World")

	query.FindTree("Indonesia", "Singapura")
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("Execution Time:", duration)
}