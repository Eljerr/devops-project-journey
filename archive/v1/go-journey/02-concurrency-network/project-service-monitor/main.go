package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type Target struct {
	Name string
	URL  string
}

type jobResult struct {
	Index int
	Data  Result
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
		slog.Error("Gagal menghubungi Service", "name", target.Name, "error", err)
		return Result{
			Name:   target.Name,
			Status: "ERROR",
			Detail: err.Error(),
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Gagal menghubungi Service", "name", target.Name, "error", err)
		return Result{
			Name:   target.Name,
			Status: "ERROR",
			Detail: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		slog.Info("Service checked", "name", target.Name, "status", "UP", "detail", "200 OK")
		return Result{
			Name:   target.Name,
			Status: "UP",
			Detail: resp.Status,
		}
	} else {
		slog.Warn("Service returned non-200 status", "name", target.Name, "status", resp.StatusCode)
		return Result{
			Name:   target.Name,
			Status: "DOWN",
			Detail: fmt.Sprintf("HTTP Status: %d", resp.StatusCode),
		}
	}
}

func checkAllService(ctx context.Context, targets []Target) []Result {
	resultsChan := make(chan jobResult, len(targets))
	var wg sync.WaitGroup

	for i, target := range targets {
		wg.Add(1)

		go func(i int, t Target) {
			defer wg.Done()

			hasil := checkService(ctx, t)
			resultsChan <- jobResult{Index: i, Data: hasil}
		}(i, target)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	final := make([]Result, len(targets))
	for res := range resultsChan {
		final[res.Index] = res.Data
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	outputFile := flag.String("output", "", "path file output")
	flag.Parse()

	if *outputFile == "" {
		slog.Error("flag -output harus diisi (contoh: go run main.go -output hasil.json)")
		os.Exit(1)
	}

	fmt.Println("=== Berikut Audit Status Server ===")

	slog.Info("Memulai pengecekan...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := checkAllService(ctx, targetService)

	err := saveResult(*outputFile, results)
	if err != nil {
		fmt.Println("Gagal menyimpan", err)
	}

	slog.Info("File berhasil dibuat...")
}
