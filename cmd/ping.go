package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var host string
var port string

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a host or IP address.",
	Run: func(cmd *cobra.Command, args []string) {
		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-interruptChan:
					// Handle interrupt signal, cleanup if needed
					fmt.Println("\nReceived interrupt signal. Exiting...")
					os.Exit(0)
				default:
					output, err := runCommandAndGetOutput("ping", "-c", "4", host)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error running ping: %v\n", err)
						os.Exit(1)
					}
					fmt.Println(output)

					// Sleep for a while before the next iteration
					time.Sleep(5 * time.Second)
				}
			}
		}()

		// Block until an interrupt signal is received
		<-interruptChan
	},
}

func init() {

	RootCmd.AddCommand(pingCmd)

	// Add flags for host and port
	pingCmd.Flags().StringVarP(&host, "host", "H", "", "Host to ping")
	pingCmd.Flags().StringVarP(&port, "port", "P", "", "Port to ping")

}

func runCommandAndGetOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
