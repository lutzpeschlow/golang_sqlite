package main

// (0) libraries
import (
	// "bufio"
	"fmt"
	"os"

	// "path/filepath"
	"runtime"
	// "strconv"
	// "strings"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
	// "strings"
	// "io/ioutil"
	// "math/rand"
	// "time"
	// "bufio"
	"github.com/lutzpeschlow/golang_sqlite/control"
	// "github.com/lutzpeschlow/golang_sqlite/readdata"
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
	err := control.ReadControlFile("control.txt", &ctrl_obj, osName)
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
		err := readdata.getData(ctrl_obj.DataDir, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
			// err := writeDb(db_name, &mod_obj)
		}
		err = writedata.writeDb(ctrl_obj.DbName, &mod_obj)
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
