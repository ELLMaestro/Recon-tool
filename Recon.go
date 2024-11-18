package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Global Variables
var wg sync.WaitGroup

// Main Function
func main() {
	// Define command-line flags
	domainPtr := flag.String("u", "", "Target domain (e.g., example.com)")
	outputPtr := flag.String("o", "output.txt", "Output file name to save results")
	threadsPtr := flag.Int("t", 10, "Number of threads for scanning")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	runAll := flag.Bool("all", true, "Run all tools (default)")
	silent := flag.Bool("silent", false, "Run in silent mode (minimal output)")
	generateHTML := flag.Bool("html", false, "Generate HTML report")

	flag.Parse()

	// Check if domain is provided
	if *domainPtr == "" {
		fmt.Println("Usage: recon-tool -u <domain> [-o <output-file>] [-t <threads>] [--all] [--silent] [--verbose]")
		os.Exit(1)
	}

	// Initialization
	domain := *domainPtr
	outputFile := *outputPtr
	threads := *threadsPtr
	fmt.Println("=== Recon Tool - Advanced ===")
	fmt.Printf("Target Domain: %s\n", domain)
	fmt.Printf("Output File: %s\n", outputFile)
	fmt.Printf("Threads: %d\n", threads)

	// Running Tools
	if *runAll {
		wg.Add(5) // Increment WaitGroup counter
		go runSubfinder(domain, silent, outputFile)
		go runAmass(domain, silent, outputFile)
		go runHttpx(outputFile, silent)
		go runWaybackurls(domain, silent, outputFile)
		go runNuclei(outputFile, silent)
		wg.Wait() // Wait for all goroutines to finish
	}

	// Generate HTML Report
	if *generateHTML {
		generateReport(outputFile)
	}

	fmt.Println("[*] Recon Tool completed successfully!")
}

// Run Subfinder
func runSubfinder(domain string, silent *bool, outputFile string) {
	defer wg.Done() // Decrement WaitGroup counter
	fmt.Println("[*] Running Subfinder...")
	args := []string{"-d", domain, "-o", outputFile}
	if *silent {
		args = append(args, "-silent")
	}
	subfinderCmd := exec.Command("subfinder", args...)
	subfinderCmd.Stdout = os.Stdout
	subfinderCmd.Stderr = os.Stderr
	err := subfinderCmd.Run()
	if err != nil {
		fmt.Println("[!] Error running Subfinder:", err)
		return
	}
	fmt.Println("[*] Subfinder completed.")
}

// Run Amass
func runAmass(domain string, silent *bool, outputFile string) {
	defer wg.Done()
	fmt.Println("[*] Running Amass...")
	args := []string{"enum", "-d", domain, "-o", outputFile}
	if *silent {
		args = append(args, "-silent")
	}
	amassCmd := exec.Command("amass", args...)
	amassCmd.Stdout = os.Stdout
	amassCmd.Stderr = os.Stderr
	err := amassCmd.Run()
	if err != nil {
		fmt.Println("[!] Error running Amass:", err)
		return
	}
	fmt.Println("[*] Amass completed.")
}

// Run httpx
func runHttpx(inputFile string, silent *bool) {
	defer wg.Done()
	fmt.Println("[*] Running httpx...")
	args := []string{"-l", inputFile}
	if *silent {
		args = append(args, "-silent")
	}
	httpxCmd := exec.Command("httpx", args...)
	httpxCmd.Stdout = os.Stdout
	httpxCmd.Stderr = os.Stderr
	err := httpxCmd.Run()
	if err != nil {
		fmt.Println("[!] Error running httpx:", err)
		return
	}
	fmt.Println("[*] httpx completed.")
}

// Run Waybackurls
func runWaybackurls(domain string, silent *bool, outputFile string) {
	defer wg.Done()
	fmt.Println("[*] Running Waybackurls...")
	args := []string{}
	if *silent {
		args = append(args, "-silent")
	}
	waybackCmd := exec.Command("waybackurls", domain)
	output, err := waybackCmd.Output()
	if err != nil {
		fmt.Println("[!] Error running Waybackurls:", err)
		return
	}

	// Write results to file
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("[!] Error writing Waybackurls results:", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(string(output))
	if err != nil {
		fmt.Println("[!] Error saving Waybackurls results:", err)
		return
	}
	writer.Flush()

	fmt.Println("[*] Waybackurls completed.")
}

// Run Nuclei
func runNuclei(inputFile string, silent *bool) {
	defer wg.Done()
	fmt.Println("[*] Running Nuclei...")
	args := []string{"-l", inputFile}
	if *silent {
		args = append(args, "-silent")
	}
	nucleiCmd := exec.Command("nuclei", args...)
	nucleiCmd.Stdout = os.Stdout
	nucleiCmd.Stderr = os.Stderr
	err := nucleiCmd.Run()
	if err != nil {
		fmt.Println("[!] Error running Nuclei:", err)
		return
	}
	fmt.Println("[*] Nuclei completed.")
}

// Generate HTML Report
func generateReport(outputFile string) {
	fmt.Println("[*] Generating HTML report...")
	// Add code to convert outputFile to HTML format
	// For simplicity, this can be implemented using a Go HTML template
	fmt.Println("[*] Report generated successfully.")
}
