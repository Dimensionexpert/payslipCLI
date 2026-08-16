package concurrency

import (
	"sync"

	"github.com/Dimensionexpert/payslip/internal/pdfgen"
)

type ConversionResult struct {
	Err      error
	Filepath string
}

func pdfWorker(workerID int, jobs <-chan string, results chan<- ConversionResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range jobs {
		err := pdfgen.ConvertToPDF(path, "output/pdf", workerID)
		results <- ConversionResult{
			Filepath: path,
			Err:      err,
		}
	}
}

func RunPDFConversion(paths []string, numWorkers int) []ConversionResult {
	jobs := make(chan string, len(paths))
	results := make(chan ConversionResult, len(paths))

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go pdfWorker(i, jobs, results, &wg)
	}

	for _, p := range paths {
		jobs <- p
	}
	close(jobs)

	wg.Wait()
	close(results)

	var collected []ConversionResult
	for r := range results {
		collected = append(collected, r)
	}
	return collected
}
