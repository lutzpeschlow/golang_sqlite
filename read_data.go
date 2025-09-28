package main

import (
	"fmt"
	// "bufio"
	// "os"
	"strings"
	"io/ioutil"
	// "math/rand"   
	// "time"
)


func getDataFiles(dir string) ([]string, error) {
    // variables
    var dataFiles []string
	fmt.Println(" get data file list ... ", dir)

    // list of all files in directory dir
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, fmt.Errorf("Fehler beim Lesen des Verzeichnisses: %v", err)
    }

    // go through files and check pre and suffix
    for _, file := range files {
        if !file.IsDir() {
            // get file name itself
            filename := file.Name()
            // starts with "data" and ends with ".txt"
            if strings.HasPrefix(filename, "data") && strings.HasSuffix(filename, ".txt") {
                dataFiles = append(dataFiles, filename)
            }
        }
    }
    fmt.Printf(" number of data files found: %d\n", len(dataFiles))
    return dataFiles, nil
}






func getDataFields(filenames []string){
    for _, filename := range filenames {
        // open and read
        content, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Printf("Fehler beim Lesen von %s: %v", filename, err)
            continue
        }
        fmt.Println(string(content))
    }
}














