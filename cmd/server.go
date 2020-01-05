package main

import (
	"flag"
	"log"

	"github.com/caledfwlch1/enlabtest/server"
)

var (
	srvip   = flag.String("i", server.DefaultIpAddr, "listen on address")
	srvport = flag.String("p", server.DefaultPort, "listen on port")
)

func main() {
	flag.Parse()
	conf := server.NewConfig(*srvip, *srvport)

	err := server.Load(conf)
	if err != nil {
		log.Println(err)
	}
}
