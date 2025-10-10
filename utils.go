package main

import (
	"bufio"
	"os"
	"strings"
	// "io/ioutil"
)

// readControlFile
func ReadControlFile(path string, obj *Control_Object) error {
	// Defaults
	obj.Action = "FEED"
	obj.DataDir = "."
	// get file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// read content from file object and scan
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// trim and split the line
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			switch parts[0] {
			case "ACTION":
				obj.Action = parts[1]
			case "DATALOC":
				obj.DataDir = parts[1]
			}
		}
	}
	return scanner.Err()
}
