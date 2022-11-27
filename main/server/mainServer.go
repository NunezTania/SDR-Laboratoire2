package main

import (
	pm "SDR-Laboratoire1/main/server/processMutex"
	"SDR-Laboratoire1/main/server/server"
)

// Main function
func main() {
	conf := pm.Config
	server.LaunchNServ(conf.NServ)
}
