package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Payload struct for JSON request
type Payload struct {
	ToSort [][]int `json:"to_sort"`
}

// Response struct for JSON response
type Response struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortArraySequential(arr []int) {
	sort.Ints(arr)
}

func sortArrayConcurrent(arr []int, wg *sync.WaitGroup) {
	defer wg.Done()
	sort.Ints(arr)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	for _, subArray := range payload.ToSort {
		sortArraySequential(subArray)
	}

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds()

	response := Response{
		SortedArrays: payload.ToSort,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var wg sync.WaitGroup

	for _, subArray := range payload.ToSort {
		wg.Add(1)
		go sortArrayConcurrent(subArray, &wg)
	}

	wg.Wait()

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds()

	response := Response{
		SortedArrays: payload.ToSort,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)


	port := 8000
	fmt.Printf("Server is running on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

