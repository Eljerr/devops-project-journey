package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	Name   string
	Status string
	Detail string
}

type Target struct {
	Name string
	URL  string
}

var targetService = []Target{
	{Name: "Nginx Server", URL: "https://httpbin.io/status/200"},
	{Name: "Database", URL: "https://httpbin.io/status/500"},
	{Name: "Payment-gateway", URL: "https://httpbin.io/status/404"},
}

func checkService(ctx context.Context, target Target) Result {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		target.URL,
		nil,
	)
	if err != nil {
		return Result{
			Name:   target.Name,
			Status: "DOWN",
			Detail: err.Error(),
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{
			Name:   target.Name,
			Status: "DOWN",
			Detail: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return Result{
			Name:   target.Name,
			Status: "UP",
			Detail: resp.Status,
		}
	}

	return Result{
		Name:   target.Name,
		Status: "DOWN",
		Detail: fmt.Sprintf("HTTP Status: %d", resp.StatusCode),
	}
}

func checkAllService(ctx context.Context, targets []Target) []Result {
	var final []Result
	var wg sync.WaitGroup

	results := make(chan Result, len(targets))

	for _, target := range targets {
		wg.Add(1)

		go func() {
			defer wg.Done()

			result := checkService(ctx, target)
			results <- result
		}()
	}

	wg.Wait()
	close(results)

	for result := range results {
		final = append(final, result)
	}

	return final
}

func saveResult(outputFile string, hasil []Result) error {
	data, err := json.MarshalIndent(hasil, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}

func main() {
	outputFile := flag.String("output", "", "path file output")
	flag.Parse()

	if *outputFile == "" {
		log.Fatal("flag -output harus diisi (contoh: go run main.go -output hasil.json)")
	}

	fmt.Println("=== Berikut Audit Status Server ===")

	log.Println("Memulai pengecekan...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := checkAllService(ctx, targetService)

	err := saveResult(*outputFile, results)
	if err != nil {
		fmt.Println("Gagal menyimpan", err)
	}

	log.Println("File berhasil dibuat...")
}
