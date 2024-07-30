package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type RequestData struct {
	Nipp     string `json:"nipp"`
	Password string `json:"password"`
}

func sendPostRequest(url string, data RequestData, wg *sync.WaitGroup, ch chan<- time.Duration) {
	defer wg.Done()

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	start := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	duration := time.Since(start)
	ch <- duration
}

func main() {
	url := "http://asia-southeast2-ordinal-stone-389604.cloudfunctions.net/login-1/"
	data := RequestData{
		Nipp:     "1204044",
		Password: "12345678",
	}

	numRequests := 100
	var wg sync.WaitGroup
	ch := make(chan time.Duration, numRequests)

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendPostRequest(url, data, &wg, ch)
	}

	wg.Wait()
	close(ch)

	totalDuration := time.Since(start)
	var sum time.Duration
	for duration := range ch {
		sum += duration
	}

	averageDuration := sum / time.Duration(numRequests)
	rps := float64(numRequests) / totalDuration.Seconds()

	fmt.Printf("Total Time: %s\n", totalDuration)
	fmt.Printf("Average Time per Request: %s\n", averageDuration)
	fmt.Printf("Requests Per Second (RPS): %.2f\n", rps)
}
