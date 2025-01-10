package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func checkMxRecord(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("MX lookup failed for %s: %w", domain, err)
	}
	return mxRecords, nil
}

func checkCnameRecord(domain string) (string, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return "", fmt.Errorf("CNAME lookup failed for %s: %w", domain, err)
	}
	return cname, nil
}

func checkTxtRecord(domain string) ([]string, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return nil, fmt.Errorf("TXT lookup failed for %s: %w", domain, err)
	}
	return txtRecords, nil
}

func checkNsRecord(domain string) ([]*net.NS, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		return nil, fmt.Errorf("NS lookup failed for %s: %w", domain, err)
	}
	return nsRecords, nil
}

func resolveDomain(domain string) ([]string, error) {
	addresses, err := net.LookupHost(domain)
	if err != nil {
		return nil, fmt.Errorf("host resolution failed for %s: %w", domain, err)
	}
	return addresses, nil
}

func printResults(domain string, mx []*net.MX, cname string, txt []string, nsrecord []*net.NS, addresses []string, err error) {
	if err != nil {
		fmt.Printf("Error for %s: %v\n", domain, err)
		return
	}

	fmt.Println("Results for", domain)
	if mx != nil {
		fmt.Println("\nMX Records:")
		for _, record := range mx {
			fmt.Printf("  Host: %s, Priority: %d\n", record.Host, record.Pref)
		}
	}
	if cname != "" {
		fmt.Println("\nCNAME Record:")
		fmt.Printf("  %s\n", cname)
	}
	if txt != nil {
		fmt.Println("\nTXT Records:")
		for _, record := range txt {
			fmt.Printf("  %s\n", record)
		}
	}
	if nsrecord != nil {
		fmt.Println("\nNS Records:")
		for _, record := range nsrecord {
			fmt.Printf("  %s\n", record.Host)
		}
	}
	if addresses != nil {
		fmt.Println("\nIP Addresses:")
		for _, addr := range addresses {
			fmt.Printf("  %s\n", addr)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a domain or host name: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	domain := strings.TrimSpace(input) // handles both \n and \r\n

	mx, mxErr := checkMxRecord(domain)
	cname, cnameErr := checkCnameRecord(domain)
	txt, txtErr := checkTxtRecord(domain)
	addresses, resolveErr := resolveDomain(domain)
	ns, nsErr := checkNsRecord(domain)

	// Consolidated printing of results and potential errors
	printResults(domain, mx, cname, txt, ns, addresses, consolidateErrors(mxErr, cnameErr, txtErr, nsErr, resolveErr))

	for {
		fmt.Print("\n\nPress Enter to continue...")
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		fmt.Println()
		break
	}
}

func consolidateErrors(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err // returns the first error found for simplicity
		}
	}
	return nil // if no errors are found, it will return nil
}
