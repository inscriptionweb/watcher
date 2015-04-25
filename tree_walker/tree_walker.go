package tree_walker

import "log"
import "time"
import "os"


// All logic to browse files
type TreeWalker struct {

}

// TreeWalker entrypoint
func (c *TreeWalker) Process(path *string) {

	file, fileError := os.Open(*path)
	fileStat, fileStatError := file.Stat()

	if fileError != nil || fileStatError != nil || !fileStat.IsDir() {
		log.SetPrefix("checker | ")
		log.Fatal(log.Ldate, " | No file found ")

		return
	}

	log.Print("Analyze folder " + *path)

	c.walk(file, path)
}

// Recursive function to browse files
func (c *TreeWalker) walk(file *os.File, path *string) *os.File {
	fileStat, fileStatError := file.Stat()

	if fileStat.IsDir() && fileStatError == nil {

		files, filesError := file.Readdir(0)

		if filesError == nil {

			for i := 0; i < len(files); i++ {
				builtPath := *path + "/" + files[i].Name()
				subFile, subFileError := os.Open(builtPath)

				if subFileError == nil && subFile != nil {
					c.walk(subFile, &builtPath)
				}
			}
		}
	}

	log.Print("Find " + file.Name())

	return file
}
