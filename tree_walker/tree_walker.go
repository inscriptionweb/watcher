package tree_walker

import "fmt"
import "os"
import "time"
import "github.com/Sirupsen/logrus"

const (
	NO_ERROR = 0
	PATH_NOT_FOUND = 1
	PATH_NOT_A_DIRECTORY = 2
	PATH_STAT_FAILURE = 3
)

type TreeWalkerError struct {
	Path string
	code uint32
}

func (e TreeWalkerError) Error() string {
	return fmt.Sprintf("An error occured at path %v : %v", e.Path, e.CodeString())
}

func (e TreeWalkerError) CodeString() string {
	switch e.code {
	case PATH_NOT_FOUND:
		return  "PATH_NOT_FOUND"
	case PATH_NOT_A_DIRECTORY:
		return  "PATH_NOT_A_DIRECTORY"
	case PATH_STAT_FAILURE:
		return  "PATH_STAT_FAILURE"
	}

	return "NO_ERROR"
}

func (e TreeWalkerError) CodeInteger() uint32 {
	return e.code
}

// All logic to browse files
type TreeWalker struct {
	MaxChangeTime time.Duration
	ExcludedFolderNames map[string]bool
}

// Constructor
func NewTreeWalker(maxChangeTime time.Duration, excludedFolderNames map[string]bool) *TreeWalker {
	return &TreeWalker{
		MaxChangeTime: maxChangeTime,
		ExcludedFolderNames: excludedFolderNames,
	}
}

// TreeWalker entrypoint
func (c *TreeWalker) Process(path *string) (*[]string, TreeWalkerError) {
	var files []string

	error := c.walk(*path, &files)

	return &files, error
}

// Recursive function to browse files
func (c *TreeWalker) walk(path string, files *[]string) (TreeWalkerError) {

	file, error := os.Open(path)
	defer file.Close()

	if error != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Error("Path not found")

		return TreeWalkerError{
			path,
			PATH_NOT_FOUND,
		}
	}

	if fileStat, error := file.Stat(); error != nil || !fileStat.IsDir() {
		if (error != nil) {

			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Error("Can't stat path")

			return TreeWalkerError{
				path,
				PATH_STAT_FAILURE,
			}
		}

		if (!fileStat.IsDir()) {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Error("Path is not a directory")

			return TreeWalkerError{
				path,
				PATH_NOT_A_DIRECTORY,
			}
		}

	}

	if filesInfo, error := file.Readdir(0); error == nil {

		for i := 0; i < len(filesInfo); i++ {

			fileName := filesInfo[i].Name()
			builtPath := path + "/" + fileName

			if c.filterByFolderName(fileName) {
				if subFileStat, error := os.Stat(builtPath); error == nil {
					if !subFileStat.IsDir() {
						if c.filterFileByDate(builtPath) {

							logrus.WithFields(logrus.Fields{
								"file": builtPath,
							}).Info("Add file")

							*files = append(*files, builtPath)
						}
					} else {
						if treeWalkerError := c.walk(builtPath, files); treeWalkerError.CodeInteger() != 0 {
							return treeWalkerError
						}

						logrus.WithFields(logrus.Fields{
							"folder": builtPath,
						}).Info("Browse folder")
					}

				}
			}

		}
	}

	return TreeWalkerError{
		path,
		NO_ERROR,
	}
}

// Filter folder name
func (c *TreeWalker) filterByFolderName(fileName string) bool {
	if len(c.ExcludedFolderNames) == 0 {
		return true
	}

	return  !c.ExcludedFolderNames[fileName]
}

// Filter according to duration
func (c *TreeWalker) filterFileByDate(path string) bool {
	if c.MaxChangeTime == 0 {
		return true
	}

	if stat, error := os.Stat(path); error == nil {
		return time.Since(stat.ModTime()) < c.MaxChangeTime
	}

	return false
}
