package tree_walker

import "log"
import "time"
import "os"

// All logic to browse files
type TreeWalker struct {
	MaxChangeTime time.Duration
	ExcludedPaths map[string]bool
	fileNames     []string
}

// Constructor
func NewTreeWalker(maxChangeTime time.Duration, excludedPaths map[string]bool) *TreeWalker {
	return &TreeWalker{
		MaxChangeTime: maxChangeTime,
		ExcludedPaths: excludedPaths,
	}
}

// TreeWalker entrypoint
func (c *TreeWalker) Process(path *string) []string {
	c.walk(*path)

	return c.fileNames
}

// Recursive function to browse files
func (c *TreeWalker) walk(path string) {

	file, error := os.Open(path)
	defer file.Close()

	if error != nil {
		log.Fatal(log.Ldate, " - No folder found : "+path)
	}

	if fileStat, error := os.Stat(path); error != nil || !fileStat.IsDir() {
		log.Fatal(log.Ldate, " - No folder found : "+path)
	}

	if filesInfo, error := file.Readdir(0); error == nil {

		for i := 0; i < len(filesInfo); i++ {

			fileName := filesInfo[i].Name()
			builtPath := path + "/" + fileName

			if c.filterByPattern(fileName, c.ExcludedPaths) {
				if subFileStat, error := os.Stat(builtPath); error == nil {
					if !subFileStat.IsDir() {
						if c.filterFileByDate(builtPath, c.MaxChangeTime) {

							c.fileNames = append(c.fileNames, builtPath)

							log.Print("Find " + builtPath)
						}
					} else {
						c.walk(builtPath)
					}

				}
			}

		}
	}
}

// Filter unwanted prefix
func (c *TreeWalker) filterByPattern(fileName string, excludedPaths map[string]bool) bool {
	if excludedPaths[fileName] {
		return false
	}

	return true
}

// Filter according to duration
func (c *TreeWalker) filterFileByDate(path string, maxChangeTime time.Duration) bool {
	if maxChangeTime == 0 {
		return true
	}

	if stat, error := os.Stat(path); error == nil {
		return time.Since(stat.ModTime()) < maxChangeTime
	}

	return false
}
