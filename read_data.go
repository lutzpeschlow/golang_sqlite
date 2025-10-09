package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// "math/rand"
	// "time"
)

func getData(dir string, model *Model) error {
	// variables
	var file_count int = 0
	var res_count int = 0
	var dataFiles []string
	fmt.Println(" get data file list ... ", dir)

	// list of all files in directory dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ERROR reading: %v", err)
	}

	// go through files in directory and check pre and suffix
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

	// feed the object content

	for i, file := range dataFiles {
		fmt.Printf("%d: %s\n", i+1, file)
		f := File{ID: uint(i + 1),
			Name: file}
		model.Files = append(model.Files, f)
	}

	for _, fname := range dataFiles {
		// create full path and filename
		fullPath := filepath.Join(dir, fname)
		file, err := os.Open(fullPath)
		if err != nil {
			fmt.Printf("ERROR: opening %s: %v\n", fullPath, err)
			continue
		}
		defer file.Close()
		// file counter for id setting
		file_count = file_count + 1
		// content of file into result objects
		scanner := bufio.NewScanner(file)
		fmt.Printf("... reading %s\n", fname)
		for scanner.Scan() {
			// fmt.Println(scanner.Text())
			res_count = res_count + 1
			res_value, _ := strconv.Atoi(scanner.Text())
			r := Result{ID: uint(res_count),
				Score:  res_value,
				FileID: file_count}
			model.Results = append(model.Results, r)
		}

	}

	return err
}
