package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"streamjury/confighandler"
	"streamjury/connections"
	"streamjury/protector"
	"time"
)

var reStripTrailingSlash = regexp.MustCompile(`/$`)
var loggingFileAbsPath string

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	binDir := path.Dir(ex)
	os.Chdir(binDir)
	fmt.Printf("Program directory is %s\n", binDir)

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s: config-file-path\n", os.Args[0])
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	// Log to file
	loggingFileAbsPath = path.Join(binDir, "streamjury.log")
	fmt.Printf("Logging to file %s\n", loggingFileAbsPath)
	log.SetFlags(log.Ldate | log.Ltime)
	f, err := os.OpenFile(
		loggingFileAbsPath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		fmt.Printf("Error opening file %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	filedata, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panicf("Couldn't read config file: %s\nERR: %s",
			os.Args[1],
			err)
	}
	config := confighandler.LoadConfig(filedata)
	protector.Protect(config.StreamjuryConfig.ResultsAbsPath)
	connections.SetConfigValues(config)

	// Start the Telegram connections
	connections.ConnectionHandler()

}
