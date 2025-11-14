package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/golang_sqlite/ctrl"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// objects:

// common model object containing:
// - Files
// - Results
// as two lists of objects
type Model struct {
	Files   []File
	Results []Result
}

// result object contains the numbers itself
// referencing to the data file
type Result struct {
	ID     uint `gorm:"primaryKey"`
	Number int
	Score  int
	FileID int
}

// file object saves the file name which
// was source of the data stored later in result object
type File struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

// dbcontent as short storage report object to see what
// is already saved in the sqlite database
type DbContent struct {
	FileNames    []string
	LastFileID   uint
	LastResultID uint
}

// ============================================================================
// functions

// debugPrintoutDbontent function for debug printout of dbcontent object
func debugPrintoutDbontent(obj *DbContent) {
	fmt.Print("debug printout of dbcontent: \n")
	fmt.Print(" ", obj.FileNames, " ", obj.LastFileID, " ", obj.LastResultID, "\n")
}

// findMaxID function to find a max id in an integer array of uint
func findMaxID(ids []uint) uint {
	var return_value uint
	if len(ids) == 0 {
		return 0
	}
	return_value = ids[0]
	for _, id := range ids {
		if id > return_value {
			return_value = id
		}
	}
	return return_value
}

// ReadDbInfo function is reading a possible existing sqlite database
// and stores some significant information for later usage in the
// dbcontent object
//
// input:
//   - db_name: string
//   - dbcontent: as pointer to object dbcontent
//
// output:
//   - error
func ReadDbInfo(db_name string, dbcontent *DbContent) error {
	// variables
	var files []File
	var results []Result

	// open sqlite database
	fmt.Println("check db content ...")
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		return err
	}

	// save close database - defer mode
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// read file information from sqlite using automatic gorm data assignment
	// executes sql command:  SELECT * FROM files;
	if err := db.Find(&files).Error; err != nil {
		return err
	}
	fmt.Print("number of filenames found: ", len(files), "\n")
	// Create slices with exact capacity
	file_ids := make([]uint, len(files))
	filenames := make([]string, len(files))
	// Loop and assign directly by index
	for i, file := range files {
		file_ids[i] = file.ID
		filenames[i] = file.Name
	}
	fmt.Print(" ", file_ids, " ", filenames, "\n")
	// write data into dbcontent object
	dbcontent.FileNames = filenames
	dbcontent.LastFileID = findMaxID(file_ids)

	// read result information from sqlite
	if err := db.Find(&results).Error; err != nil {
		return err
	}
	fmt.Printf("number of results found: %d\n", len(results))
	// slice with exact length and find max value
	res_ids := make([]uint, len(results))
	for i, result := range results {
		res_ids[i] = result.ID
	}
	dbcontent.LastResultID = findMaxID(res_ids)
	fmt.Print(" ", dbcontent.FileNames, " ", dbcontent.LastFileID, " ", dbcontent.LastResultID, "\n")

	// return error value
	return err
}

// GetData function to read data from pre-defined directory from files
// including  *data*  in file name
// the directory is delivered in control object
//
// before reading the data files a base of already existing data
// in possible existing sqlite database is checked
//
// if there is duplicate data, the new read data is reduced accordingly
//
// input:
//   - dir: pre-defined directoy
//   - model: pointer to model_object
//
// output:
//   - err
func GetData(ctrl *ctrl.Control_Object, model *Model) error {
	// variables
	var file_count int = 0
	var res_count int = 0
	var line_count int = 0
	var dataFiles []string
	dbcontent := DbContent{}

	// check db content for possible later data reduction
	err := ReadDbInfo(ctrl.DbName, &dbcontent)
	if err != nil {
		return fmt.Errorf("ERROR %v", err)
	}
	debugPrintoutDbontent(&dbcontent)

	// get a list of all files in directory
	fmt.Println("get data file list in dir... ", ctrl.DataDir)
	files, err := os.ReadDir(ctrl.DataDir)
	if err != nil {
		return fmt.Errorf("ERROR %v", err)
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

	// feed the files object content
	for i, file := range dataFiles {
		fmt.Print("   ", i+1, " ", file, "\n")
		f := File{ID: uint(i + 1), Name: file}
		model.Files = append(model.Files, f)
	}

	for _, fname := range dataFiles {
		// create full path and filename
		fullPath := filepath.Join(ctrl.DataDir, fname)
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
		line_count = 0
		for scanner.Scan() {
			res_count = res_count + 1
			line_count = line_count + 1
			res_value, _ := strconv.Atoi(scanner.Text())
			// assign result object and append it to the list
			r := Result{ID: uint(res_count),
				Number: line_count,
				Score:  res_value,
				FileID: file_count}
			model.Results = append(model.Results, r)
		}
	}
	return err
}
