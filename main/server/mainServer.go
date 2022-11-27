package main

import (
	pm "SDR-Laboratoire1/main/server/processMutex"
	"SDR-Laboratoire1/main/server/server"
)

func main() {
	conf := pm.Config
	server.LaunchNServ(conf.NServ)
}
