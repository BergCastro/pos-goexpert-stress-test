package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	statusCode int
	duration   time.Duration
}

func worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{statusCode: 0, duration: duration}
			continue
		}

		results <- Result{statusCode: resp.StatusCode, duration: duration}
		resp.Body.Close()
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("A URL é obrigatória")
		return
	}

	totalRequests := *requests
	concurrentRequests := *concurrency

	results := make(chan Result, totalRequests)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go worker(*url, totalRequests/concurrentRequests, results, &wg)
	}

	wg.Wait()
	close(results)

	totalDuration := time.Since(start)
	statusCounts := make(map[int]int)
	var totalRequestsMade int

	for result := range results {
		statusCounts[result.statusCode]++
		totalRequestsMade++
	}

	fmt.Printf("Tempo total gasto: %v\n", totalDuration)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequestsMade)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", statusCounts[200])
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range statusCounts {
		if status != 200 {
			fmt.Printf("Status %d: %d\n", status, count)
		}
	}
}