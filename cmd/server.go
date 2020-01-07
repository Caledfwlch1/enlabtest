package main

import (
	"flag"
	"log"

	"github.com/caledfwlch1/enlabtest/server"
)

var (
	srvIp   = flag.String("i", server.DefaultIpAddr, "listen on address")
	srvPort = flag.String("p", server.DefaultPort, "listen on port")
	connStr = flag.String("c", "", "database connection string")
)

func main() {
	flag.Parse()
	conf, err := server.NewConfig(*srvIp, *srvPort, *connStr)
	if err != nil {
		log.Fatalln(err)
	}

	err = server.ListenAndServe(conf)
	if err != nil {
		log.Fatalln(err)
	}
}
