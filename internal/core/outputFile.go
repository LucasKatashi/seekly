package core

import (
	"bufio"
	"fmt"
	"os"
)

func OutputFile(domains []string, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, domain := range domains {
		if _, err := writer.WriteString(domain + "\n"); err != nil {
			return fmt.Errorf("cannot write output file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("cannot flush output file: %w", err)
	}
	return nil
}
