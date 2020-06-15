package logging

import (
	"os"
	"strings"

	"github.com/alecthomas/log4go"
)

func InitConfigs() {
	log4goInit()
	log4go.LoadConfiguration("logging/log4go.xml")
	log4go.Info("log4go init ok.")
}

func log4goInit() {
	path, _ := os.Getwd()
	path = strings.Replace(path, "\\", "/", -1)
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	path = path[0:strings.LastIndex(path, "/")]
	path = path[0:strings.LastIndex(path, "/")+1] + "log"
	if !pathExists(path) {
		log4go.Warn("dir: logs/ not found.")
		err := os.MkdirAll(path, 0711)
		if err != nil {
			log4go.Error(err.Error())
		}
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
