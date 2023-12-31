package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	targetURL   string
	concurrency int
	numRequests int
)

var cpuStressCmd = &cobra.Command{
	Use:   "cpu-stress",
	Short: "A CPU stress testing tool for applications.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running CPU stress test on %s with concurrency %d and %d requests\n", targetURL, concurrency, numRequests)

		var wg sync.WaitGroup

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				client := resty.New()

				for j := 0; j < numRequests; j++ {
					startTime := time.Now()
					_, err := client.R().Get(targetURL)
					if err != nil {
						fmt.Printf("Request error: %v\n", err)
					} else {
						duration := time.Since(startTime)
						fmt.Printf("Request took %s\n", duration)
					}
				}
			}()
		}

		wg.Wait()
	},
}

func init() {
	cpuStressCmd.Flags().StringVarP(&targetURL, "url", "u", "", "URL to test (required)")
	cpuStressCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent users (default 1)")
	cpuStressCmd.Flags().IntVarP(&numRequests, "requests", "r", 100, "Number of requests per user (default 100)")

	cpuStressCmd.MarkFlagRequired("url")

	// Register the cpu-stress command with the root command.
	RootCmd.AddCommand(cpuStressCmd)
}
