package io

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// WriteDb function to write data into sqlite database
//
// input:
//   - db_name: as name of the database
//   - model: pointer to model_object
//
// output:
//   - err
func WriteDb(db_name string, model *Model) error {
	var dbExist bool
	var files []File
	fmt.Println("write db ...")
	fmt.Println(db_name)

	// existence of database
	if _, err := os.Stat(db_name); err == nil {
		fmt.Println("WARNING: database already exists")
		dbExist = true
	} else if os.IsNotExist(err) {
		fmt.Println("database will be created from scratch")
		dbExist = false
	}

	// connect SQLite DB
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		return err
	}
	// if database exists, read already saved data files
	if dbExist {
		if err := db.Find(&files).Error; err != nil {
			return err
		}
		fmt.Printf(" number of filenames found: %d\n", len(files))
		fmt.Println("no further data written")

	} else {
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
	}
	return nil
}
