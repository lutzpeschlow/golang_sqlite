package main

// (0) libraries
import (
	"fmt"
	"os"
	"runtime"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
	// "strings"
	// "io/ioutil"
	// "math/rand"
	// "time"
	// "bufio"
)

// (1) ojects
type Control_Object struct {
	Action  string
	DataDir string
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

func writeDb(db_name string, model *Model) error {
	fmt.Println("write db ...")
	fmt.Println(db_name)
	// connect SQLite DB
	// db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// err = db.AutoMigrate(&File{}, &Result{})
	// if err != nil {
	// 	panic("migration failed")
	// }
	// return err
}

// ======================================================================================
// (2) main
// ======================================================================================
func main() {
	// instance of objects
	ctrl_obj := Control_Object{}
	mod_obj := Model{}

	var db_name string
	db_name = "g.db"
	// system check
	fmt.Println(runtime.GOOS)

	// get control flag stored in control file and save in object
	err := ReadControlFile("control.txt", &ctrl_obj)
	if err != nil {
		fmt.Printf(" %v\n", err)
		os.Exit(1)
	}
	// content of control object
	fmt.Println("Settings:")
	fmt.Printf(" action: %s\n", ctrl_obj.Action)
	fmt.Printf(" datadir: %s\n", ctrl_obj.DataDir)

	// case handler
	// - FEED        - which calls to feed the database
	// - CONTENT     - get content of the database
	switch ctrl_obj.Action {
	case "FEED":
		err := getData(ctrl_obj.DataDir, &mod_obj)
		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
			err := writeDb(db_name, &mod_obj)
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
