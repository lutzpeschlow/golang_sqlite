package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	// "io/ioutil"
)

// readControlFile 
func readControlFile(filename string) (string, error) {
	// file and error object after reading file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("ERROR reading control file: %v", err)
	}
	defer file.Close()
    // read file content
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// check possible keywords
		if line == "FEED" || line == "CONTENT" {
			return line, nil
		} else {
			return "", fmt.Errorf("missing valid control keywords   FEED oder CONTENT", line)
		}
	}
	return "", fmt.Errorf("ERROR: control file is empty")
}



