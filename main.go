package main

import "github.com/antham/watcher/tree_walker"
import "github.com/antham/watcher/sender"
import "gopkg.in/alecthomas/kingpin.v1"
import "strings"
import "time"
import "github.com/Sirupsen/logrus"

var (
	verbose               = kingpin.Flag("verbose", "Report every operation occuring").Bool()
	maxChangeTime         = kingpin.Flag("max-change-time", "Maximal change time").Default("9s").Duration()
	intervalTime          = kingpin.Flag("interval-time", "Interval between two check").Default("10s").Duration()
	excludedFoldersString = kingpin.Flag("excluded-paths", "Folder to exclude from lookup separated with comma").String()
	username              = kingpin.Flag("username", "Ssh username").String()
	host                  = kingpin.Flag("host", "Ssh host").String()
	keyFile               = kingpin.Flag("key-file", "Ssh keyfile").String()

	localPath  = kingpin.Arg("local-path", "Local pathname to lookup").Required().String()
	remotePath = kingpin.Arg("remote-path", "Remote pathname to copy data in").Required().String()
)

func main() {
	kingpin.Parse()

	var loggingLevel logrus.Level

	if *verbose {
		loggingLevel = logrus.DebugLevel
	} else {
		loggingLevel = logrus.FatalLevel
	}

	logger := logrus.New()
	logger.Level = loggingLevel

	excludedFolders := make(map[string]bool)

	for _, value := range strings.Split(*excludedFoldersString, ",") {
		excludedFolders[value] = true
	}

	var fileSender *sender.Sender
	var senderError sender.SenderError

	treeWalker := tree_walker.NewTreeWalker(*maxChangeTime, excludedFolders, logger)

	if fileSender, senderError = sender.NewSender(*username, *host, *keyFile, *localPath, *remotePath, logger); senderError.CodeInteger() != sender.NO_ERROR {
		logger.WithFields(logrus.Fields{
			"message": senderError.Error(),
		}).Fatal("Something wrong happened")
	}

	go func() {

		for _ = range time.Tick(*intervalTime) {

			var files *[]string
			var treeWalkerError tree_walker.TreeWalkerError

			if files, treeWalkerError = treeWalker.Process(localPath); treeWalkerError.CodeInteger() != tree_walker.NO_ERROR {
				logger.WithFields(logrus.Fields{
					"message": treeWalkerError.Error(),
				}).Fatal("Something wrong happened")
			}

			if senderError = fileSender.Send(files); senderError.CodeInteger() != sender.NO_ERROR {
				logger.WithFields(logrus.Fields{
					"message": senderError.Error(),
				}).Fatal("Something wrong happened")
			}
		}
	}()

	select {}
}
