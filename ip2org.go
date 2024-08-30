package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

const version = "v0.0.1"
const batchSize = 50 // Number of IP addresses to process in each batch

func main() {
	var ip string
	var ipListFile string
	var outputFile string
	var verbose bool
	var reverseLookup bool
	var timeout int
	var showVersion bool

	flag.StringVar(&ip, "ip", "", "IP address to lookup")
	flag.StringVar(&ipListFile, "list", "", "File containing IP addresses")
	flag.StringVar(&outputFile, "o", "", "File to save output")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&reverseLookup, "reverse", false, "Perform reverse DNS lookup")
	flag.IntVar(&timeout, "timeout", 10, "Timeout for whois lookup in seconds")
	flag.BoolVar(&showVersion, "version", false, "print version information and exit")
	flag.Parse()

	if showVersion {
		fmt.Println("ip2org version:", version)
		return
	}

	var outFile *os.File
	var err error

	if outputFile != "" {
		outFile, err = os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
	}

	if ip == "" && ipListFile == "" {
		// Check if input is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			processIPsFromScanner(scanner, outFile, verbose, reverseLookup, timeout)
		} else {
			fmt.Println("IP address or file not provided")
		}
	} else if ip != "" {
		lookupIP(sanitizeIP(ip), outFile, verbose, reverseLookup, timeout)
	} else if ipListFile != "" {
		file, err := os.Open(ipListFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		processIPsFromScanner(scanner, outFile, verbose, reverseLookup, timeout)
	}
}

func processIPsFromScanner(scanner *bufio.Scanner, outFile *os.File, verbose, reverseLookup bool, timeout int) {
	ipChannel := make(chan string) // Channel to receive IP addresses
	var wg sync.WaitGroup

	// Start Goroutines to perform lookups
	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChannel {
				lookupIP(sanitizeIP(ip), outFile, verbose, reverseLookup, timeout)
			}
		}()
	}

	// Send IP addresses to the channel
	for scanner.Scan() {
		ip := scanner.Text()
		ipChannel <- ip
	}

	close(ipChannel)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait() // Wait for all Goroutines to finish
}

func lookupIP(ip string, outFile *os.File, verbose, reverseLookup bool, timeout int) {
	// Set up timeout for the whois command
	cmd := exec.Command("whois", "-h", "whois.arin.net", ip)
	cmd.Stdout = new(strings.Builder)
	cmd.Stderr = new(strings.Builder)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		cmd.Process.Kill()
		log.Printf("Lookup for %s timed out\n", ip)
	case err := <-done:
		if err != nil {
			log.Fatal(err)
		}
	}

	out := cmd.Stdout.(*strings.Builder).String()
	if verbose {
		fmt.Printf("Raw output for %s:\n%s\n", ip, out)
	}

	orgName := ""
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OrgName:") {
			orgName = strings.TrimSpace(strings.TrimPrefix(line, "OrgName:"))
			break
		}
	}

	if reverseLookup {
		if domain, err := reverseDNS(ip); err == nil {
			orgName = fmt.Sprintf("%s (Reverse DNS: %s)", orgName, domain)
		}
	}

	result := fmt.Sprintf("%s [%s]\n", ip, orgName)

	// print output to terminal
    fmt.Print(result)

    // write output a file
	if outFile != nil {
		_, err := outFile.WriteString(result)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// sanitizeIP removes http://, https://, and :port from the IP address
func sanitizeIP(ip string) string {
	// Remove http://, https://, and :port
	re := regexp.MustCompile(`^(?:http://|https://|)([\d\.]+)(?::\d+)?$`)
	matches := re.FindStringSubmatch(ip)
	if len(matches) > 1 {
		return matches[1]
	}
	return ip
}

// reverseDNS performs a reverse DNS lookup to get the domain name associated with the IP address
func reverseDNS(ip string) (string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}
	if len(names) > 0 {
		return names[0], nil
	}
	return "", nil
}
