package main

import (
	"encoding/json"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"runtime"
)

// Web UI default Listen address
var G_AlexAddr = "0.0.0.0:8000"

// Job Storage Default Url
var G_MongoUrl = "localhost:27017"

// Global MongoDB object
var G_MongoSession *mgo.Session
var G_MongoDB *mgo.Database

// vegeta jobs current running
var G_RunningVegetaJobs = NewConcurrentSet()

// vegeta jobs will stopping
var G_StoppingVegetaJobs = NewConcurrentSet()

// boom jobs current running
var G_RunningBoomJobs = NewConcurrentSet()

// boom jobs will stopping
var G_StoppingBoomJobs = NewConcurrentSet()

// teams for grouping jobs
var G_AlexTeams = []string{"python"}

// Display Html page layout
var G_ShowLayout = true

// Configuration Object
type Config struct {
	BindAddr   string
	MongoUrl   string
	Teams      []string
	ShowLayout bool
}

// Load Config from external json file
func LoadConfig() {
	var cfile = flag.String("c", "", "json config file path")
	flag.Parse()
	if *cfile != "" {
		file, err := os.Open(*cfile)
		if err != nil {
			log.Panic("open config file fail")
		}
		decoder := json.NewDecoder(file)
		config := Config{}
		err = decoder.Decode(&config)
		if err != nil {
			log.Panic("config file not valid json")
		}
		G_AlexAddr = config.BindAddr
		G_AlexTeams = config.Teams
		G_MongoUrl = config.MongoUrl
		G_ShowLayout = config.ShowLayout
	}
}

func InitGlobals() {
	session, err := mgo.Dial(G_MongoUrl)
	if err != nil {
		log.Panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	G_MongoSession = session
	G_MongoDB = session.DB("alex")
	// set golang threads num
	runtime.GOMAXPROCS(runtime.NumCPU())
}
