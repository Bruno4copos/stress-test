package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func main() {
	urlPtr := flag.String("url", "", "URL do serviço a ser testado (obrigatório)")
	requestsPtr := flag.Int("requests", 100, "Número total de requests a serem enviados")
	concurrencyPtr := flag.Int("concurrency", 10, "Número de chamadas simultâneas")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Erro: O parâmetro --url é obrigatório.")
		os.Exit(1)
	}

	url := *urlPtr
	totalRequests := *requestsPtr
	concurrency := *concurrencyPtr

	fmt.Printf("Iniciando teste de carga para: %s\n", url)
	fmt.Printf("Total de requests: %d\n", totalRequests)
	fmt.Printf("Concorrência: %d\n", concurrency)

	startTime := time.Now()
	results := make(chan Result, totalRequests)
	var wg sync.WaitGroup

	// Controla o número de goroutines simultâneas
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Adquire um slot do semáforo
		go func() {
			defer func() {
				<-semaphore // Libera o slot do semáforo
				wg.Done()
			}()

			start := time.Now()
			resp, err := http.Get(url)
			duration := time.Since(start)

			if err != nil {
				fmt.Printf("Erro ao fazer request para %s: %v\n", url, err)
				results <- Result{StatusCode: 0, Duration: duration} // Código 0 para erro
				return
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body) // Consome o body para liberar a conexão

			results <- Result{StatusCode: resp.StatusCode, Duration: duration}
		}()
	}

	wg.Wait()
	close(results)
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Println("\n--- Relatório de Teste de Carga ---")
	fmt.Printf("Tempo total de execução: %s\n", totalTime)

	statusCodeCounts := make(map[int]int)
	totalSuccess := 0

	for result := range results {
		statusCodeCounts[result.StatusCode]++
		if result.StatusCode >= 200 && result.StatusCode < 300 {
			totalSuccess++
		}
	}

	fmt.Printf("Total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Requests com status HTTP 200: %d\n", totalSuccess)

	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for code, count := range statusCodeCounts {
		if code != 200 && (code < 200 || code >= 300) {
			fmt.Printf("  %d - %s: %d\n", code, http.StatusText(code), count)
		}
	}
}
