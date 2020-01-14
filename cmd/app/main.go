package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/mongo"
	"github.com/vishwanathj/protovnfdparser/pkg/server"
	"github.com/vishwanathj/protovnfdparser/pkg/service"
	//"log"
)

const envLogFile = "LOGFILE"

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

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

const (
	mongoIP   = "localhost"
	mongoPort = 27017
	dbName    = "go_web_server"
	collName  = "vnfd"
)

var mgoIPPtr = flag.String("ip", mongoIP, "hostname or IP address")
var mgoPortPtr = flag.Int("port", mongoPort, "mongodb port number")
var mgoDbName = flag.String("dbname", dbName, "name of database")
var mgoCollName = flag.String("coll", collName, "collection name")

func main() {
	//mgoIPPtr := flag.String("ip", mongoIP, "hostname or IP address")
	//mgoPortPtr := flag.Int("port", mongoPort, "mongodb port number")
	//mgoDbName := flag.String("dbname", dbName, "name of database")
	//mgoCollName := flag.String("coll", collName, "collection name")
	flag.Parse()

	log.WithFields(log.Fields{
		"MongoDB_IP":    *mgoIPPtr,
		"MongoPort":     *mgoPortPtr,
		"MongoDB_Name":  *mgoDbName,
		"MongoCollName": *mgoCollName,
	}).Info("App startup parameters")
	//log.Info(*mgoIPPtr, *mgoPortPtr, *mgoDbName, *mgoCollName)
	//fmt.Println(mgoCollName, mgoDbName)

	dbURI := fmt.Sprintf("%s:%d", *mgoIPPtr, *mgoPortPtr)
	d, err := mongo.NewMongoDAL(dbURI, dbName, collName)
	if err != nil {
		log.Fatal("unable to connect to mongodb")
		log.Debug(err)
	}
	v := service.NewVnfdService(d)
	s := server.NewServer(v, nil)

	s.Start()
}
