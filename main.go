package main

// (0) libraries
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "strings"
	// "io/ioutil"
	// "math/rand"
	// "time"
	// "bufio"
)

// ========== object declaration ==============================================
//
// (1) ojects
type Control_Object struct {
	Action  string
	DataDir string
	DbName  string
}

type Model struct {
	Files   []File
	Results []Result
}

type Result struct {
	ID     uint `gorm:"primaryKey"`
	Number int
	Score  int
	FileID int
}

type File struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

// ===== planned package: pkg =================================================
//
// # ReadControlFile
//
// function to read a control file that contains values for the tool
// what action should be made
// where are the data
// what should be the name of the sqlite database
//
// input:
//
//	path - file name as string
//	control_object - pointer to control_object
//	osname - type of operating system
//
// output:
//
//	error - error value
func ReadControlFile(path string, obj *Control_Object, osName string) error {
	// Defaults
	obj.Action = "FEED"
	obj.DataDir = "."
	// pointer to file for later opening, err as interface value
	file, err := os.Open(path)
	// if we have an error go out with returning err value
	if err != nil {
		return err
	}
	// close file at the end of the function
	defer file.Close()
	// read content from file object and scan
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// trim and split the line
		// parts is a slice of strings (dynamical array)
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			switch parts[0] {
			case "ACTION":
				obj.Action = parts[1]
			case "DATALOC_WIN":
				if osName == "windows" {
					obj.DataDir = parts[1]
				}
			case "DATALOC_LINUX":
				if osName == "linux" {
					obj.DataDir = parts[1]
				}
			case "DBNAME":
				obj.DbName = parts[1]
			}
		}
	}
	// return value is the error interface value of the scanner
	return scanner.Err()
}

// ===== planned package: read data ===========================================
//
// get data
func getData(dir string, model *Model) error {
	// variables
	var file_count int = 0
	var res_count int = 0
	var line_count int = 0
	var dataFiles []string
	fmt.Println(" get data file list ... ", dir)

	// list of all files in directory dir

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ERROR %v", err)
	}

	// go through files in directory and check pre and suffix
	for _, file := range files {
		if !file.IsDir() {
			// get file name itself
			filename := file.Name()
			// fmt.Println(filename)
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
		line_count = 0
		fmt.Printf("... reading %s\n", fname)
		for scanner.Scan() {
			// fmt.Println(scanner.Text())
			res_count = res_count + 1
			line_count = line_count + 1
			res_value, _ := strconv.Atoi(scanner.Text())

			r := Result{ID: uint(res_count),
				Number: line_count,
				Score:  res_value,
				FileID: file_count}
			model.Results = append(model.Results, r)
		}

	}
	return err
}

// ===== write data from model object to sqlite database ======================
//
// using gorm so write to sqlite
// the objects are already adapted for direct export to sqlite
func writeDb(db_name string, model *Model) error {
	fmt.Println("write db ...")
	fmt.Println(db_name)
	// connect SQLite DB
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		return err
	}
	// automatic migration of tables for file and result attributes
	err = db.AutoMigrate(&File{}, &Result{})
	if err != nil {
		return err
	}
	// save files
	for _, file := range model.Files {
		if err := db.Create(&file).Error; err != nil {
			return err
		}
	}
	// save results
	for _, result := range model.Results {
		if err := db.Create(&result).Error; err != nil {
			return err
		}
	}
	fmt.Println("Database write completed.")
	return nil
}

// ======================================================================================
// (2) main
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := Control_Object{}
	mod_obj := Model{}

	// system check
	osName := runtime.GOOS
	fmt.Println(osName)

	// get control flag stored in control file and save in object
	err := ReadControlFile("control.txt", &ctrl_obj, osName)
	if err != nil {
		fmt.Printf(" %v\n", err)
		os.Exit(1)
	}
	// content of control object
	fmt.Println("Settings:")
	fmt.Printf(" action: %s\n", ctrl_obj.Action)
	fmt.Printf(" datadir: %s\n", ctrl_obj.DataDir)
	fmt.Printf(" dbname: %s\n", ctrl_obj.DbName)

	// case handler
	// - FEED        - which calls to feed the database
	// - CONTENT     - get content of the database
	switch ctrl_obj.Action {
	case "FEED":
		err := getData(ctrl_obj.DataDir, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
			// err := writeDb(db_name, &mod_obj)
		}
		err = writeDb(ctrl_obj.DbName, &mod_obj)
		if err != nil {
			fmt.Printf("error writing db: %v\n", err)
		}
	case "CONTENT":
		fmt.Println("CONTENT is active")
	}

	// content of files as one attribute of model object
	fmt.Println("Files:")
	for _, file := range mod_obj.Files {
		fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	}
	// content of results as another attribute of model object
	fmt.Println("\nResults:")
	for _, result := range mod_obj.Results {
		fmt.Printf("ID: %d, Score: %d, FileID: %d\n", result.ID, result.Score, result.FileID)
	}
}

// ============================================================================
