package main

import "github.com/antham/watcher/tree_walker"
import "github.com/antham/watcher/sender"
import "gopkg.in/alecthomas/kingpin.v1"
import "strings"

var (
	maxChangeTime       = kingpin.Flag("max-change-time", "Maximal change time").Duration()
	excludedPathsString = kingpin.Flag("excluded-paths", "Path to exclude from lookup separated with comma").String()
	username            = kingpin.Flag("username", "Ssh username").String()
	ip                  = kingpin.Flag("ip", "Ssh ip").String()
	keyFile             = kingpin.Flag("key-file", "Ssh keyfile").String()

	localPath  = kingpin.Arg("local-path", "Local pathname to lookup").Required().String()
	remotePath = kingpin.Arg("remote-path", "Remote pathname to copy data in").Required().String()
)

func main() {
	kingpin.Parse()

	excludedPaths := make(map[string]bool)

	for _, value := range strings.Split(*excludedPathsString, ",") {
		excludedPaths[value] = true
	}

	treeWalker := tree_walker.NewTreeWalker(*maxChangeTime, excludedPaths)
	fileNames := treeWalker.Process(localPath)

	fileSender := sender.NewSender(*username, *ip, *keyFile, *localPath, *remotePath)

	fileSender.Send(fileNames)
}
