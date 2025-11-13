package ctrl

import (
	"bufio"
	"os"
	"strings"
)

// control object
type Control_Object struct {
	Action  string
	DataDir string
	DbName  string
}

// ReadControlFile function to read a control file from path and
// fills the values into object Control_Object. osName defines the operating
// system and can be used for platform dependent actions or settings
//
// input:
//   - path: file name as string
//   - obj: pointer to control_object
//   - osName: type of operating system
//
// output:
//   - error: if read or parse fails, put back error, else nil
func ReadControlFile(path string, obj *Control_Object, osName string) error {
	// defaults
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
