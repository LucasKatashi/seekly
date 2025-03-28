package core

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func OutputFile(domains []string, filename string) {
	red := color.New(color.FgRed)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("[%s] Failed to create output file '%s': %v\n", red.Sprint("ERR"), filename, err)
		os.Exit(1)
	}

	writer := bufio.NewWriter(file)

	for _, domain := range domains {
		writer.WriteString(domain + "\n")
	}

	writer.Flush()
}
