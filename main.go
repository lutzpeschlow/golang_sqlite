package main

// (0) libraries
import (
	"fmt"
	"runtime"

	// "bufio"
	"os"
	// "strings"
	// "io/ioutil"
	// "math/rand"
	// "time"
)

// (1) ojects
type Control_Object struct {
	Ctrl string
}

type Data_Object struct {
	data_files []string
}

// ======================================================================================
// (2) main
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := Control_Object{}
	data_set := Data_Object{}

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

		files, err := getDataFiles(dir)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
		data_set.data_files = files

		for i, file := range files {
			fmt.Printf("%d: %s\n", i+1, file)
		}

		getDataFields(files)

	case "CONTENT":
		fmt.Println("CONTENT is active")
	}

}

// TDL:
// - build object  data_set
// - include data reader
