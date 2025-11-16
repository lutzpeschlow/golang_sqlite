package main

// libraries
import (
	"fmt"
	"os"
	"runtime"

	"github.com/lutzpeschlow/golang_sqlite/ctrl"
	"github.com/lutzpeschlow/golang_sqlite/io"
)

// ======================================================================================
// main function
// - read control file, according defined action
//   - read data files and store in sqlite
//   - get content of sqlite database
//
// input:
//
//	None
//
// output:
//
//	None
//
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := ctrl.Control_Object{}
	mod_obj := io.Model{}
	db_content := io.DbContent{}

	// system check
	osName := runtime.GOOS
	fmt.Print("OS: ", osName, "\n")

	// get control flag stored in control file and save in object
	err := ctrl.ReadControlFile("control.txt", &ctrl_obj, osName)
	if err != nil {
		fmt.Printf(" %v\n", err)
		os.Exit(1)
	}
	// content of control object
	fmt.Print("Settings: ", ctrl_obj.Action, " ", ctrl_obj.DataDir, " ", ctrl_obj.DbName, "\n")

	// case handler
	// - FEED        - which calls to feed the database
	// - CONTENT     - get content of the database
	switch ctrl_obj.Action {
	case "FEED":
		// read data from data files
		err := io.GetData(&ctrl_obj, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
		// write data into sqlite database
		err = io.WriteDb(ctrl_obj.DbName, &mod_obj)
		if err != nil {
			fmt.Printf("error writing db: %v\n", err)
		}
	case "CONTENT":
		fmt.Println("CONTENT is active")
		fmt.Print(" read db info ... \n")
		err = io.ReadDbInfo(ctrl_obj.DbName, &db_content)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}
		io.DebugPrintoutDbontent(&db_content)
	}

	// content of files as one attribute of model object
	// fmt.Println("Files:")
	// for _, file := range mod_obj.Files {
	// 	fmt.Printf("ID: %d, Name: %s\n", file.ID, file.Name)
	// }
	// // content of results as another attribute of model object
	// fmt.Println("\nResults:")
	// for _, result := range mod_obj.Results {
	// 	fmt.Printf("ID: %d, Score: %d, FileID: %d\n", result.ID, result.Score, result.FileID)
	// }

}

// ============================================================================
