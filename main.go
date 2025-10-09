package main

// (0) libraries
import (
	"fmt"
	"os"
	"runtime"
	// "strings"
	// "io/ioutil"
	// "math/rand"
	// "time"
	// "bufio"
)

// (1) ojects
type Control_Object struct {
	Ctrl string
}

type Model struct {
	Files   []File
	Results []Result
}

type Result struct {
	ID     uint `gorm:"primaryKey"`
	Score  int
	FileID int
}

type File struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

// ======================================================================================
// (2) main
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := Control_Object{}
	mod_obj := Model{}

	// set directory for data files, this will be set later on another way
	var dir string
	if runtime.GOOS == "windows" {
		dir = "c:\\tmp"
	} else {
		dir = "/home/lutz/test"
	}
	fmt.Println("directory setting: ", dir)

	// get control flag stored in control file and save in object
	ctrl, err := readControlFile("control.txt")
	if err != nil {
		fmt.Printf(" %v\n", err)
		os.Exit(1)
	}
	ctrl_obj.Ctrl = ctrl
	fmt.Printf("ctrl keyword: %s\n", ctrl_obj.Ctrl)

	// case handler
	// - FEED        - which calls to feed the database
	// - CONTENT     - get content of the database
	switch ctrl_obj.Ctrl {
	case "FEED":
		fmt.Println("FEED is active")

		err := getData(dir, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}

	case "CONTENT":
		fmt.Println("CONTENT is active")
	}

	// content of files
	fmt.Println("Files:")
	for _, file := range mod_obj.Files {
		fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	}
	// content of results
	fmt.Println("\nResults:")
	for _, result := range mod_obj.Results {
		fmt.Printf("ID: %d, Score: %d, FileID: %d\n", result.ID, result.Score, result.FileID)
	}
}

// TDL:
// - build object  data_set
// - include data reader
