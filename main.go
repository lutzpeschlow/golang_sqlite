package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
	
)
// "math/rand"    "time"


// oject CONTROL
type Control_Object struct {
	Ctrl string
}

// object DATA


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





func getDataFiles(dir string) ([]string, error) {
	fmt.Println(" get data file list ... ")
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, fmt.Errorf("Fehler beim Lesen des Verzeichnisses: %v", err)
    }

    var dataFiles []string
    
    for _, file := range files {
        if !file.IsDir() {
            filename := file.Name()
            // Prüfe ob Datei mit "data" beginnt und .txt Endung hat
            if strings.HasPrefix(filename, "data") && strings.HasSuffix(filename, ".txt") {
                dataFiles = append(dataFiles, filename)
            }
        }
    }
    
    return dataFiles, nil
}


// ======================================================================================
// main
// ======================================================================================
func main() {
	// instance of con
	dir := "/tmp"
	ctrl_obj := Control_Object{}
	
	// Control-File einlesen
	ctrl, err := readControlFile("control.txt")
	if err != nil {
		fmt.Printf(" %v\n", err)
		os.Exit(1)
	}
	
	// save attribute
	ctrl_obj.Ctrl = ctrl
	fmt.Printf("ctrl keyword: %s\n", ctrl_obj.Ctrl)
	
	// 
	switch ctrl_obj.Ctrl {
	case "FEED":
		fmt.Println("FEED is active")

		files, err := getDataFiles(dir)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
        for i, file := range files {
            fmt.Printf("%d: %s\n", i+1, file)
        }


	case "CONTENT":
		fmt.Println("CONTENT is active")
	}
}

