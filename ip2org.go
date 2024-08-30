package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const batchSize = 50 // Number of IP addresses to process in each batch

func main() {
	var ip string
	var ipListFile string

	flag.StringVar(&ip, "ip", "", "IP address to lookup")
	flag.StringVar(&ipListFile, "list", "", "File containing IP addresses")

	flag.Parse()

	if ip == "" && ipListFile == "" {
		// Check if input is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			processIPsFromScanner(scanner)
		} else {
			log.Fatal("IP address or file not provided")
		}
	} else if ip != "" {
		lookupIP(ip)
	} else if ipListFile != "" {
		file, err := os.Open(ipListFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		processIPsFromScanner(scanner)
	}
}

func processIPsFromScanner(scanner *bufio.Scanner) {
	ipChannel := make(chan string) // Channel to receive IP addresses
	var wg sync.WaitGroup

	// Start Goroutines to perform lookups
	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChannel {
				lookupIP(ip)
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

func lookupIP(ip string) {
	out, err := exec.Command("whois", "-h", "whois.arin.net", ip).Output()
	if err != nil {
		log.Fatal(err)
	}

	orgName := ""
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OrgName:") {
			orgName = strings.TrimSpace(strings.TrimPrefix(line, "OrgName:"))
			break
		}
	}

	fmt.Printf("%s [%s]\n", ip, orgName)
}
