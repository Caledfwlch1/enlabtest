package server

import "time"

const (
	configFileName        = "server.yml"
	DefaultIpAddr         = ""
	DefaultPort           = "8080"
	timeoutServerShutdown = time.Minute
)
