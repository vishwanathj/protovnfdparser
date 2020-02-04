package main

import (
	"fmt"
	"io"
	"os"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/server"
	"github.com/vishwanathj/protovnfdparser/pkg/service"
	//"log"
)

const envLogFile = "LOGFILE"

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	/*log.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			repopath := fmt.Sprintf("%s/src/github.com/vishwanathj", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})*/

	//var filename  = "/tmp/go_web_server.log"
	var filename = os.Getenv(envLogFile)
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
		//log.SetOutput(f)
	}

	log.SetLevel(log.DebugLevel)

	log.SetReportCaller(true)
}

func main() {
	cfg := config.GetConfigInstance()
	v, err := service.GetVnfdServiceInstance(*cfg)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	s := server.NewServer(v, cfg.WebConfig)

	s.Start()
}
