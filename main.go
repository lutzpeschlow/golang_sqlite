package main

// (0) libraries
import (
	"fmt"
	"os"
	"runtime"

	"github.com/lutzpeschlow/golang_sqlite/ctrl"
	"github.com/lutzpeschlow/golang_sqlite/io"
)

// ======================================================================================
// (2) main
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := ctrl.Control_Object{}
	mod_obj := io.Model{}

	// system check
	osName := runtime.GOOS
	fmt.Println(osName)

	// get control flag stored in control file and save in object
	err := ctrl.ReadControlFile("control.txt", &ctrl_obj, osName)
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
		err := io.GetData(&ctrl_obj, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
			// err := writeDb(db_name, &mod_obj)
		}
		err = io.WriteDb(ctrl_obj.DbName, &mod_obj)
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
