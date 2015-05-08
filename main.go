package main

import "github.com/antham/watcher/tree_walker"
import "github.com/antham/watcher/sender"
import "gopkg.in/alecthomas/kingpin.v1"
import "strings"
import "log"

var (
	maxChangeTime         = kingpin.Flag("max-change-time", "Maximal change time").Duration()
	excludedFoldersString = kingpin.Flag("excluded-paths", "Folder to exclude from lookup separated with comma").String()
	username              = kingpin.Flag("username", "Ssh username").String()
	ip                    = kingpin.Flag("ip", "Ssh ip").String()
	keyFile               = kingpin.Flag("key-file", "Ssh keyfile").String()

	localPath  = kingpin.Arg("local-path", "Local pathname to lookup").Required().String()
	remotePath = kingpin.Arg("remote-path", "Remote pathname to copy data in").Required().String()
)

func main() {
	kingpin.Parse()

	excludedFolders := make(map[string]bool)

	for _, value := range strings.Split(*excludedFoldersString, ",") {
		excludedFolders[value] = true
	}

	treeWalker := tree_walker.NewTreeWalker(*maxChangeTime, excludedFolders)

	if files, error := treeWalker.Process(localPath); error.CodeInteger() == tree_walker.NO_ERROR {
		fileSender := sender.NewSender(*username, *ip, *keyFile, *localPath, *remotePath)
		fileSender.Send(files)
	} else {
		log.Fatalf(error.Error())
	}
}
