package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	vegeta "github.com/tsenart/vegeta/lib"
)

type Endpoint struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
}
type ResultResponse struct {
	Result *vegeta.Result `json:"result"`
	// Response string         `json:"response"`
}
type ApiArgs struct {
	Conc       int
	URL        string
	Endpoint   string
	Method     string
	Headers    string
	Body       string
	OutputFile string
	Duration   time.Duration
	Report     bool
}

var apiArgs ApiArgs

// var conc int
// var url string
// var endpoint string
// var method string
// var headers string
// var body string

// // var resultFileName string
// // var responseBodyFileName string
// var outputFileName string

var apiStressCmd = &cobra.Command{

	Use:   "api-stress",
	Short: "A API stress testing tool for applications.",
	Run: func(cmd *cobra.Command, args []string) {
		headerMap := parseHeaders(apiArgs.Headers)

		target := Endpoint{
			Method:  apiArgs.Method,
			URL:     apiArgs.URL + apiArgs.Endpoint,
			Headers: headerMap,
			Body:    parseJSONBody(apiArgs.Body),
		}

		rate := vegeta.Rate{Freq: apiArgs.Conc, Per: time.Second}
		duration := apiArgs.Duration

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: target.Method,
			URL:    target.URL,
			Header: createHeader(target.Headers),
			Body:   createJSONBody(target.Body),
		})

		attacker := vegeta.NewAttacker()

		results := make(chan ResultResponse)
		done := make(chan struct{})
		outputFileName := apiArgs.OutputFile
		go func() {
			defer close(results)
			for res := range attacker.Attack(targeter, rate, duration, "Vegeta!") {
				resultResponse := ResultResponse{
					Result: res,
				}
				results <- resultResponse
			}
		}()
		go func() {
			successCount := 0
			errorCount := 0
			for result := range results {
				if result.Result.Code == 200 {
					successCount++
				} else {
					errorCount++
				}
				//fmt.Printf("Result: %s\n", result)
			}

			totalRequests := successCount + errorCount
			successPercentage := float64(successCount) / float64(totalRequests) * 100
			errorPercentage := float64(errorCount) / float64(totalRequests) * 100
			fmt.Printf("Success Count: %d (%.2f%%)\n", successCount, successPercentage)
			fmt.Printf("Error Count: %d (%.2f%%)\n", errorCount, errorPercentage)
			close(done)
		}()

		resultFile, err := os.Create(outputFileName)
		if err != nil {
			//fmt.Printf("Error opening result file: %v", err)

			return
		}
		defer resultFile.Close()

		responseBodyFile, err := os.Create("response_bodies.json")
		if err != nil {
			fmt.Printf("Error opening response body file: %v", err)
			return
		}
		defer responseBodyFile.Close()
		for result := range results {
			resultJSON, err := json.Marshal(result)
			if err != nil {
				fmt.Printf("Error encoding result to JSON: %v", err)
			}

			//outputFile.WriteString(string(resultJSON) + "\n")
			resultFile.WriteString(string(resultJSON) + "\n")
			fmt.Printf("Result: %s", string(resultJSON))

		}

		<-done

	},
}

func parseHeaders(headerString string) map[string]string {
	headerMap := make(map[string]string)
	if headerString != "" {
		headers := strings.Split(headerString, ",")
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return headerMap
}

func createHeader(headers map[string]string) map[string][]string {
	header := make(map[string][]string)
	for key, value := range headers {
		header[key] = []string{value}
	}
	return header
}

func parseJSONBody(body string) interface{} {
	if body == "" {
		return nil
	}

	var bodyData interface{}
	err := json.Unmarshal([]byte(body), &bodyData)
	if err != nil {
		fmt.Printf("Error parsing JSON body: %v\n", err)
		return nil
	}

	return bodyData
}

func createJSONBody(body interface{}) []byte {
	if body == nil {
		return nil
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error encoding JSON body: %v\n", err)
		return nil
	}

	return bodyJSON
}

func init() {
	//rootCmd.AddCommand(apiStressCmd)
	apiStressCmd.Flags().IntVarP(&apiArgs.Conc, "concurrency", "c", 2, "Concurrency level")
	apiStressCmd.Flags().StringVarP(&apiArgs.URL, "url", "u", "http://your-api-endpoint.com", "API URL")
	apiStressCmd.Flags().StringVarP(&apiArgs.Endpoint, "endpoint", "e", "/", "API Endpoint")
	apiStressCmd.Flags().StringVarP(&apiArgs.Method, "method", "m", "GET", "HTTP Method (PUT, GET, POST, DELETE, etc.)")
	apiStressCmd.Flags().StringVarP(&apiArgs.Headers, "headers", "H", "", "Headers in key:value format (comma-separated)")
	apiStressCmd.Flags().StringVarP(&apiArgs.Body, "body", "b", "", "JSON Body")
	apiStressCmd.Flags().BoolVarP(&apiArgs.Report, "report", "R", false, "Report summary of results including success and error percentages")
	apiStressCmd.Flags().DurationVarP(&apiArgs.Duration, "duration", "d", 10*time.Second, "Duration of the attack (e.g., 10s)") // Corrected this line
	apiStressCmd.Flags().StringVarP(&apiArgs.OutputFile, "output-file", "O", "output.json", "File to write both results and response bodies in JSON format")
	// apiStressCmd.Flags().StringVarP(&resultFileName, "result-file", "R", "results.json", "File to write results in JSON format")
	// apiStressCmd.Flags().StringVarP(&responseBodyFileName, "response-body-file", "B", "response_bodies.txt", "File to write response bodies in text format")
	RootCmd.PersistentFlags().Bool("verbose", true, "Print more information during execution")
	RootCmd.AddCommand(apiStressCmd)

}
