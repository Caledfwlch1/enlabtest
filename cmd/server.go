package main

import (
	"flag"
	"log"

	"github.com/caledfwlch1/enlabtest/server"
)

var (
	srvIp   = flag.String("i", server.DefaultIpAddr, "listen on address")
	srvPort = flag.String("p", server.DefaultPort, "listen on port")
	dbHost  = flag.String("h", "", "database host")
	dbUser  = flag.String("u", "", "database user")
	dbPass  = flag.String("w", "", "database password")
	dbName  = flag.String("n", "", "database name")
	dbOpt   = flag.String("o", "", "database connect options")
	save    = flag.Bool("s", false, "save config")
)

func main() {
	flag.Parse()
	conf, err := server.NewConfig(*srvIp, *srvPort, *dbHost, *dbUser, *dbPass, *dbName, *dbOpt, *save)
	if err != nil {
		log.Println(err)
		return
	}

	err = server.Load(conf)
	if err != nil {
		log.Println(err)
	}
}
