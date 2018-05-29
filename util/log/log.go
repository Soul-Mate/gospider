package log

import (
	"github.com/Soul-Mate/gospider/util/conf"
	"log"
	"os"
	"strings"
	"fmt"
)

func InitLog() {
	c := conf.GlobalSharedConfig.Log
	if c.Enable {
		file := c.File
		if file == "" {
			file = "gospider.log"
		}
		fd, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(fd)
	}
}

func Printf(level string, format string, v ...interface{}) {
	c := conf.GlobalSharedConfig.Log
	if strings.ContainsAny(level, c.Level) {
		log.Printf(format, v)
	} else {
		fmt.Fprintf(os.Stdout, format, v)
	}
}
