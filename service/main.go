// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/shell909090/influx-proxy/backend"
	"github.com/toolkits/file"
)

var (
	//ErrConfig   = errors.New("config parse error")
	ConfigFile  string
	//NodeName    string
	//RedisAddr   string
	LogFilePath string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.StringVar(&LogFilePath, "log", "var/app.log", "output file")
	flag.StringVar(&ConfigFile, "c", "", "config file")
	//flag.StringVar(&NodeName, "node", "l1", "node name")
	//flag.StringVar(&RedisAddr, "redis", "localhost:6379", "config file")
	flag.Parse()
}

//type Config struct {
//	redis.Options
//	Node string
//}

//func LoadJson(configfile string, cfg interface{}) (err error) {
//	file, err := os.Open(configfile)
//	if err != nil {
//		return
//	}
//	defer file.Close()
//
//	dec := json.NewDecoder(file)
//	err = dec.Decode(&cfg)
//	return
//}

func initLog() {
	if LogFilePath == "" {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LogFilePath,
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     7,
		})
	}
}

func main() {
	initLog()

	var err error
	//var cfg Config

	//if ConfigFile != "" {
	//	err = LoadJson(ConfigFile, &cfg)
	//	if err != nil {
	//		log.Print("load config failed: ", err)
	//		return
	//	}
	//	log.Printf("json loaded.")
	//}

	if ConfigFile == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(ConfigFile) {
		log.Fatalln("config file:", ConfigFile, "is not existent")
	}

	configContent, err := file.ToTrimString(ConfigFile)
	if err != nil {
		log.Fatalln("read config file:", ConfigFile, "fail:", err)
	}

	var cfg backend.GlobelConfig
	err = json.Unmarshal([]byte(configContent), &cfg)
	if err != nil {
		log.Fatalln("parse config file:", ConfigFile, "fail:", err)
	}

	//if NodeName != "" {
	//	cfg.Node = NodeName
	//}
	//
	//if RedisAddr != "" {
	//	cfg.Addr = RedisAddr
	//}

	//rcs := backend.NewRedisConfigSource(&cfg.Options, cfg.Node)
	//
	//nodecfg, err := rcs.LoadNode()
	//if err != nil {
	//	log.Printf("config source load failed.")
	//	return
	//}

	jcs := backend.NewJsonConfigSource(&cfg)
	ic := backend.NewInfluxCluster(jcs)
	ic.LoadConfig()

	mux := http.NewServeMux()
	NewHttpService(ic, cfg.Node.DB).Register(mux)

	log.Printf("http service start.")
	server := &http.Server{
		Addr:        cfg.Node.ListenAddr,
		Handler:     mux,
		IdleTimeout: time.Duration(cfg.Node.IdleTimeout) * time.Second,
	}
	if cfg.Node.IdleTimeout <= 0 {
		server.IdleTimeout = 10 * time.Second
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Print(err)
		return
	}
}
