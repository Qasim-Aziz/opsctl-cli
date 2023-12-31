package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type DnsCheck struct {
	Host  string
	Debug bool
}

var dnsCheck DnsCheck
var dnsCheckCmd = &cobra.Command{
	Use:   "dns-check",
	Short: "A collection of dns check.",
	Run: func(cmd *cobra.Command, args []string) {
		if dnsCheck.Debug {
			runDebugLookup(dnsCheck.Host)
		} else {
			ips, err := net.LookupIP(dnsCheck.Host)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
				os.Exit(1)
			}
			for _, ip := range ips {
				fmt.Printf("%s IN A %s\n", dnsCheck.Host, ip.String())
			}
		}
	},
}

func runDebugLookup(host string) {
	cmd := exec.Command("nslookup", "-debug", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running nslookup: %v\n", err)
		os.Exit(1)
	}

	// Display the output in a table
	displayDebugTable(strings.Split(string(output), "\n"))
}

func displayDebugTable(lines []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Debug Output"})
	table.SetAutoWrapText(false) // Disable automatic text wrapping

	for _, line := range lines {
		// Split each line into words to prevent long words from breaking formatting
		words := strings.Fields(line)
		table.Append(words)
	}

	table.Render()
}

func init() {
	// Add flags or additional setup for the dns-check command, if needed.

	// Register the dns-check command with the root command.
	dnsCheckCmd.Flags().StringVarP(&dnsCheck.Host, "host", "H", "example.com", "Host to lookup")
	dnsCheckCmd.Flags().BoolVarP(&dnsCheck.Debug, "debug", "d", false, "Enable debug mode with nslookup")
	RootCmd.AddCommand(dnsCheckCmd)
}
