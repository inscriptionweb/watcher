package tree_walker

import "os"
import "reflect"
import "strings"
import "testing"
import "time"

func createFolderHierarchy () {
	folders := []string{"level1", "level2", "level3", "level4", "level5", "level6"}

	os.RemoveAll("/tmp/level1")
	os.MkdirAll("/tmp/" + strings.Join(folders,"/"), 0777)

	for i := 1; i <= len(folders); i++ {
		os.Create("/tmp/" + strings.Join(folders[0:i],"/") + "/hello-world.txt")
	}
}

func TestWalkInFolder(t *testing.T) {
	createFolderHierarchy()
	duration := time.Duration(0)
	excludedFolderNames := make(map[string]bool)

	treeWalker := NewTreeWalker(duration, excludedFolderNames)

	path := "/tmp/level1"

	files, error := treeWalker.Process(&path)
	expectedResult := []string{
		"/tmp/level1/level2/level3/level4/level5/level6/hello-world.txt",
		"/tmp/level1/level2/level3/level4/level5/hello-world.txt",
		"/tmp/level1/level2/level3/level4/hello-world.txt",
		"/tmp/level1/level2/level3/hello-world.txt",
		"/tmp/level1/level2/hello-world.txt",
		"/tmp/level1/hello-world.txt",
	}

	if !reflect.DeepEqual(*files,expectedResult) || error.CodeInteger() != 0 {
		t.Error("List of files not as expected", files, error)
	}
}

func TestWalkUnexistingFolder(t *testing.T) {
	duration := time.Duration(0)
	excludedFolderNames := make(map[string]bool)

	treeWalker := NewTreeWalker(duration, excludedFolderNames)

	path := "/tmp/whatever"

	files, error := treeWalker.Process(&path)
	var expectedResult []string

	if error.CodeInteger() != PATH_NOT_FOUND || error.Path != "/tmp/whatever" || !reflect.DeepEqual(*files, expectedResult) {
		t.Error("Expect no file at all", files, error)
	}
}

func TestWalkAPathWhichPointOnFile(t *testing.T) {
	duration := time.Duration(0)
	excludedFolderNames := make(map[string]bool)

	treeWalker := NewTreeWalker(duration, excludedFolderNames)

	path := "/tmp/level1/hello-world.txt"

	files, error := treeWalker.Process(&path)
	var expectedResult []string

	if error.CodeInteger() != PATH_NOT_A_DIRECTORY || error.Path != "/tmp/level1/hello-world.txt" || !reflect.DeepEqual(*files, expectedResult) {
		t.Error("Expect error when walking a path which is not a directory", files, error)
	}
}

func TestWalkInFolderWithExcludedFolderNames(t *testing.T) {
	createFolderHierarchy()
	os.MkdirAll("/tmp/level1/level2/folder", 0777)
	os.Create("/tmp/level1/level2/folder/hello-world.txt")
	os.MkdirAll("/tmp/level1/level2/folder1", 0777)
	os.Create("/tmp/level1/level2/folder1/hello-world.txt")
	os.MkdirAll("/tmp/level1/level2/folder2", 0777)
	os.Create("/tmp/level1/level2/folder2/hello-world.txt")

	duration := time.Duration(0)
	excludedFolderNames := make(map[string]bool)
	excludedFolderNames["level4"] = true
	excludedFolderNames["folder"] = true
	excludedFolderNames["folder1"] = true
	excludedFolderNames["folder2"] = true
	excludedFolderNames["whatever"] = true

	treeWalker := NewTreeWalker(duration, excludedFolderNames)

	path := "/tmp/level1"

	files, error := treeWalker.Process(&path)
	expectedResult := []string{
		"/tmp/level1/level2/level3/hello-world.txt",
		"/tmp/level1/level2/hello-world.txt",
		"/tmp/level1/hello-world.txt",
	}

	if !reflect.DeepEqual(*files,expectedResult) || error.CodeInteger() != NO_ERROR {
		t.Error("Some files are not filtered properly", files, error)
	}
}
